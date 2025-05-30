package dynamodb

import (
	"context"
	"fmt"
	"sync"

	"github.com/agiac/rivo"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type scanOptions struct {
	PoolSize int
}

func newDefaultScanOptions() *scanOptions {
	return &scanOptions{
		PoolSize: 1,
	}
}

type ScanOption func(*scanOptions) error

func applyScanOptions(o *scanOptions, opts []ScanOption) (*scanOptions, error) {
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	return o, nil
}

func ScanPoolSize(poolSize int) ScanOption {
	return func(o *scanOptions) error {
		if poolSize < 1 {
			return fmt.Errorf("poolSize must be greater than 0")
		}
		o.PoolSize = poolSize
		return nil
	}
}

// Scan returns a pipeline which scans the provided DynamoDB table and emits the scan output responses;
func Scan(client *dynamodb.Client, input *dynamodb.ScanInput, opt ...ScanOption) rivo.Pipeline[rivo.None, *dynamodb.ScanOutput] {
	o, err := applyScanOptions(newDefaultScanOptions(), opt)
	if err != nil {
		panic(fmt.Sprintf("invalid ScanOption: %v", err))
	}

	return func(ctx context.Context, _ rivo.Stream[rivo.None]) rivo.Stream[*dynamodb.ScanOutput] {
		out := make(chan rivo.Item[*dynamodb.ScanOutput])

		go func() {
			defer close(out)

			nSegments := o.PoolSize

			wg := sync.WaitGroup{}
			wg.Add(nSegments)

			for i := 0; i < nSegments; i++ {
				go func(segment int) {
					defer wg.Done()

					input := *input
					input.TotalSegments = aws.Int32(int32(nSegments))
					input.Segment = aws.Int32(int32(segment))

					paginator := dynamodb.NewScanPaginator(client, &input)

					for paginator.HasMorePages() {
						output, err := paginator.NextPage(ctx)
						if err != nil {
							select {
							case <-ctx.Done():
								return
							case out <- rivo.Item[*dynamodb.ScanOutput]{Err: fmt.Errorf("failed to scan: %w", err)}:
								continue
							}
						}

						if output == nil {
							continue
						}

						select {
						case <-ctx.Done():
							return
						case out <- rivo.Item[*dynamodb.ScanOutput]{Val: output}:
						}
					}
				}(i)
			}

			wg.Wait()
		}()

		return out
	}
}

// ScanItems returns a pipeline which scans the provided DynamoDB table and emits the items from the scan output responses.
func ScanItems(client *dynamodb.Client, input *dynamodb.ScanInput, opt ...ScanOption) rivo.Pipeline[rivo.None, map[string]types.AttributeValue] {
	tableScan := Scan(client, input, opt...)

	itemsMap := rivo.Map[*dynamodb.ScanOutput, []map[string]types.AttributeValue](func(ctx context.Context, i rivo.Item[*dynamodb.ScanOutput]) ([]map[string]types.AttributeValue, error) {
		if i.Err != nil {
			return nil, i.Err
		}

		return i.Val.Items, nil
	})

	return rivo.Pipe3(tableScan, itemsMap, rivo.Flatten[map[string]types.AttributeValue]())
}
