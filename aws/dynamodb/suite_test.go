package dynamodb_test

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const scanTableName = "scanTableTest"
const writeTableName = "writeTableTest"
const tableItems = 1000

type Suite struct {
	suite.Suite
	dynamoC testcontainers.Container
	db      *dynamodb.Client
}

func TestSuite(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping test in short mode")
	}
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	ctx := context.Background()

	var err error
	s.dynamoC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "amazon/dynamodb-local:latest",
			ExposedPorts: []string{"8000"},
			WaitingFor:   wait.ForListeningPort("8000"),
		},
		Started: true,
	})
	s.Require().NoError(err)

	host, err := s.dynamoC.Host(ctx)
	s.Require().NoError(err)

	mappedPort, err := s.dynamoC.MappedPort(ctx, "8000")
	s.Require().NoError(err)

	s.db = CreateDynamodbClient(s.T(), fmt.Sprintf("http://%s:%d", host, mappedPort.Int()))

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
	s.Require().NoError(testcontainers.TerminateContainer(s.dynamoC))
}

func CreateDynamodbClient(t *testing.T, dynamodbEndpoint string) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithBaseEndpoint(dynamodbEndpoint),
		config.WithDefaultRegion("eu-central-1"),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "foo",
				SecretAccessKey: "bar",
			},
		}))
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
