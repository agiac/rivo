package dynamodb

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/agiac/rivo"

	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// TODO: consider using ForEachOutput function

type batchWriteOptions struct {
	PoolSize int
	ChanSize int
}

func newDefaultBatchWriteOptions() *batchWriteOptions {
	return &batchWriteOptions{
		PoolSize: 1,
		ChanSize: 0,
	}
}

type BatchWriteOption func(*batchWriteOptions) error

func applyBatchWriteOptions(o *batchWriteOptions, opts []BatchWriteOption) (*batchWriteOptions, error) {
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

func BatchWritePoolSize(poolSize int) BatchWriteOption {
	return func(o *batchWriteOptions) error {
		if poolSize < 1 {
			return fmt.Errorf("poolSize must be greater than 0")
		}
		o.PoolSize = poolSize
		return nil
	}
}

func BatchWriteChanSize(chanSize int) BatchWriteOption {
	return func(o *batchWriteOptions) error {
		if chanSize < 0 {
			return fmt.Errorf("chanSize must be greater than or equal to 0")
		}
		o.ChanSize = chanSize
		return nil
	}
}

// BatchWrite returns a pipeline which writes the input stream to the provided DynamoDB using the BatchWriteItem API.
func BatchWrite(client *awsdynamodb.Client, opt ...BatchWriteOption) rivo.Pipeline[*awsdynamodb.BatchWriteItemInput, *awsdynamodb.BatchWriteItemOutput] {
	o, err := applyBatchWriteOptions(newDefaultBatchWriteOptions(), opt)
	if err != nil {
		panic(fmt.Sprintf("invalid BatchWriteOption: %v", err))
	}

	return func(ctx context.Context, in rivo.Stream[*awsdynamodb.BatchWriteItemInput], errs chan<- error) rivo.Stream[*awsdynamodb.BatchWriteItemOutput] {
		out := make(chan *awsdynamodb.BatchWriteItemOutput)

		go func() {
			defer close(out)

			wg := sync.WaitGroup{}
			wg.Add(o.PoolSize)

			for range o.PoolSize {
				go func() {
					defer wg.Done()
					for i := range rivo.OrDone(ctx, in) {
						res, err := batchWriteItem(ctx, client, i, 0)
						if err != nil {
							select {
							case <-ctx.Done():
								return
							case errs <- err:
								continue
							}
						}

						select {
						case <-ctx.Done():
							return
						case out <- res:
						}
					}
				}()
			}

			wg.Wait()
		}()

		return out
	}
}

// BatchPutItems returns a pipeline which writes the input stream to the provided DynamoDB using the BatchWriteItem API, but only for PutItem operations;
func BatchPutItems(client *awsdynamodb.Client, tableName string, opt ...BatchWriteOption) rivo.Pipeline[types.PutRequest, *awsdynamodb.BatchWriteItemOutput] {
	batchedItems := rivo.Batch[types.PutRequest](25)

	batchWriteRequests := rivo.Map[[]types.PutRequest, *awsdynamodb.BatchWriteItemInput](func(ctx context.Context, r []types.PutRequest) (*awsdynamodb.BatchWriteItemInput, error) {
		writeRequests := make([]types.WriteRequest, 0, len(r))
		for _, putRequest := range r {
			writeRequests = append(writeRequests, types.WriteRequest{PutRequest: &putRequest})
		}

		return &awsdynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: writeRequests,
			},
		}, nil
	})

	return rivo.Pipe3(batchedItems, batchWriteRequests, BatchWrite(client, opt...))
}

func batchWriteItem(ctx context.Context, client *awsdynamodb.Client, item *awsdynamodb.BatchWriteItemInput, retries int) (*awsdynamodb.BatchWriteItemOutput, error) {
	const maxRetries = 5

	res, err := client.BatchWriteItem(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("failed to batch write items %+v: %w", item.RequestItems, err)
	}

	if len(res.UnprocessedItems) != 0 {
		if retries < maxRetries {
			item.RequestItems = res.UnprocessedItems
			time.Sleep(time.Duration(2^retries) * time.Second)
			return batchWriteItem(ctx, client, item, retries+1)
		}

		ui, _ := json.Marshal(res.UnprocessedItems)
		return nil, fmt.Errorf("failed to batch write items: max retries exceeded: unprocessed items: %s", ui)
	}

	return res, nil
}
