package dynamodb_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// TODO: use Test Containers https://golang.testcontainers.org/features/docker_compose/

const scanTableName = "scanTableTest"
const writeTableName = "writeTableTest"
const tableItems = 1000

type Suite struct {
	suite.Suite
	db *dynamodb.Client
}

func TestSuite(t *testing.T) {
	if os.Getenv("INTEGRATION_TESTS") != "true" {
		t.Skip("skipping integration tests")
	}

	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	s.db = CreateDynamodbClient(s.T())

	for _, tableName := range []string{scanTableName, writeTableName} {
		CreateDynamodbTable(s.T(), s.db, tableName)
	}

	for i := range tableItems {
		PutDynamodbItem(s.T(), s.db, scanTableName, map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("PK-%d", i)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("SK-%d", i)},
		})
	}
}

func (s *Suite) TearDownSuite() {
	for _, tableName := range []string{scanTableName, writeTableName} {
		DeleteDynamodbTable(s.T(), s.db, tableName)
	}
}

func CreateDynamodbClient(t *testing.T) *dynamodb.Client {
	dynamodbEndpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if dynamodbEndpoint == "" {
		dynamodbEndpoint = "http://localhost:8000"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithBaseEndpoint(dynamodbEndpoint))
	require.NoError(t, err)
	return dynamodb.NewFromConfig(cfg)
}

func CreateDynamodbTable(t *testing.T, client *dynamodb.Client, tableName string) {
	_, err := client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	require.NoError(t, err)
}

func DeleteDynamodbTable(t *testing.T, client *dynamodb.Client, tableName string) {
	_, err := client.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	require.NoError(t, err)
}

func PutDynamodbItem(t *testing.T, client *dynamodb.Client, tableName string, item map[string]types.AttributeValue) {
	_, err := client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	require.NoError(t, err)
}

func GetDynamodbItem(t *testing.T, client *dynamodb.Client, tableName string, key map[string]types.AttributeValue) map[string]types.AttributeValue {
	resp, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	})
	require.NoError(t, err)
	return resp.Item
}

func CountDynamodbTableItems(t *testing.T, client *dynamodb.Client, tableName string) int {
	count := 0

	var lastEvaluatedKey map[string]types.AttributeValue
	for {
		resp, err := client.Scan(context.Background(), &dynamodb.ScanInput{
			TableName:         aws.String(tableName),
			ExclusiveStartKey: lastEvaluatedKey,
		})
		require.NoError(t, err)

		count += len(resp.Items)

		lastEvaluatedKey = resp.LastEvaluatedKey
		if lastEvaluatedKey == nil {
			break
		}

	}

	return count
}
