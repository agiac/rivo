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

func (s *Suite) TestScanTable() {
	ctx := context.Background()

	errs := make(chan error, 1)

	out := rivodynamodb.Scan(s.db, &dynamodb.ScanInput{
		TableName: aws.String(scanTableName),
	}, rivodynamodb.ScanPoolSize(runtime.NumCPU()))(ctx, nil, errs)

	var got []map[string]types.AttributeValue
	for o := range out {
		got = append(got, o.Items...)
	}

	close(errs)
	for err := range errs {
		s.NoError(err)
	}

	s.Len(got, tableItems)
	for i := range tableItems {
		s.Contains(got, map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PK-%d", i)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SK-%d", i)},
		})
	}
}

func (s *Suite) TestScanTableItems() {
	ctx := context.Background()

	errs := make(chan error, 1)

	out := rivodynamodb.ScanItems(s.db, &dynamodb.ScanInput{
		TableName: aws.String(scanTableName),
	}, rivodynamodb.ScanPoolSize(runtime.NumCPU()))(ctx, nil, errs)

	var got []map[string]types.AttributeValue
	for o := range out {
		got = append(got, o)
	}

	close(errs)
	for err := range errs {
		s.NoError(err)
	}

	s.Len(got, tableItems)
	for i := range tableItems {
		s.Contains(got, map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PK-%d", i)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SK-%d", i)},
		})
	}
}
