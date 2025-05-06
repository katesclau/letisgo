package tests

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/sirupsen/logrus"
	ddbContainer "github.com/testcontainers/testcontainers-go/modules/dynamodb"
	"mnesis.com/pkg/db"
)

var (
	ddbClient *dynamodb.Client
	ddb       db.DynamoDBHandler[TestRecord]
	mtx       = &sync.Mutex{}
)

const (
	tableName = "test-table"
)

type testStage struct {
	t   *testing.T
	ddb db.DynamoDBHandler[TestRecord]
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	mustBeNilErr(err)

	// ==========================================================================
	// Local DynamoDB
	// ==========================================================================
	dynamodbContainer, err := ddbContainer.Run(context.Background(), "amazon/dynamodb-local:2.2.1", ddbContainer.WithDisableTelemetry())
	mustBeNilErr(err)

	defer dynamodbContainer.Terminate(ctx) //nolint:errcheck

	dynamoEndpoint, err := dynamodbContainer.Endpoint(ctx, "http")
	mustBeNilErr(err)

	ddbClient = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider("key", "secret", "session")
		o.BaseEndpoint = &dynamoEndpoint
		o.Region = "us-west-2"
	})
	mustBeNilErr(err)

	ddb = db.NewDynamoDBHandler[TestRecord](
		ddbClient,
		tableName,
	)

	m.Run()
}

func newTestStage(t *testing.T) *testStage {
	mtx.Lock()
	defer mtx.Unlock()

	stage := testStage{
		t:   t,
		ddb: ddb,
	}

	os.Setenv("ENV", "test")
	InitDynamoTables()
	logrus.Info("DDB Table created")

	return &stage
}

func InitDynamoTables() {
	ctx := context.Background()

	// Table: listings
	_, err := ddbClient.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []types.KeySchemaElement{
			{KeyType: types.KeyTypeHash, AttributeName: aws.String("pk")},
			{KeyType: types.KeyTypeRange, AttributeName: aws.String("sk")},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("sk"), AttributeType: types.ScalarAttributeTypeS},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(100),
			WriteCapacityUnits: aws.Int64(100),
		},
	})
	mustBeNilErr(err)
}

func mustBeNilErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
