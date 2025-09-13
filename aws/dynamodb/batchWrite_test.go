package dynamodb_test

import (
	"context"
	"fmt"
	"runtime"

	rivodynamodb "github.com/agiac/rivo/aws/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (s *Suite) TestBatchPutItems() {
	ctx := context.Background()

	errs := make(chan error, 1)

	in := make(chan types.PutRequest)
	go func() {
		defer close(in)
		for i := range tableItems {
			in <- types.PutRequest{
				Item: map[string]types.AttributeValue{
					"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PK-%d", i)},
					"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SK-%d", i)},
				},
			}
		}
	}()

	out := rivodynamodb.BatchPutItems(s.db, writeTableName, rivodynamodb.BatchWritePoolSize(runtime.NumCPU()), rivodynamodb.BatchWriteChanSize(1))(ctx, in, errs)
	for o := range out {
		s.NotNil(o)
	}

	close(errs)
	for err := range errs {
		s.NoError(err)
	}

	for i := range tableItems {
		item, err := s.db.GetItem(context.Background(), &dynamodb.GetItemInput{
			TableName: aws.String(writeTableName),
			Key: map[string]types.AttributeValue{
				"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PK-%d", i)},
				"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SK-%d", i)},
			},
		})
		s.NoError(err)
		s.NotNil(item.Item)
	}
}
