package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"mnesis.com/pkg/store"
)

var dynamoClient *dynamodb.Client //nolint:gochecknoglobals,gofumpt

type Data struct {
	PK        string    `json:"pk"`
	SK        string    `json:"sk"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

type testStage struct {
	t     *testing.T
	store store.Interface[Data]
}

// TestMain test initialization that runs only once before running all tests
func TestMain(m *testing.M) {
	ctx := context.Background()

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	os.Setenv("ENV", "dev")
	mustBeNilErr(err)

	// ==========================================================================
	// Local DynamoDB
	// ==========================================================================
	req := testcontainers.ContainerRequest{
		Image:        "amazon/dynamodb-local",
		ExposedPorts: []string{"8000/tcp"},
		WaitingFor:   wait.ForExposedPort(),
	}
	dynamo, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	mustBeNilErr(err)

	dynamoEndpoint, err := dynamo.Endpoint(ctx, "http")
	mustBeNilErr(err)

	dynamoClient = dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.Credentials = credentials.NewStaticCredentialsProvider("key", "secret", "session")
		o.BaseEndpoint = &dynamoEndpoint
		o.Region = "us-west-2"
	})

	m.Run()
}

// newTestStage initializes and returns a listings testing stage
func newTestStage(t *testing.T) testStage {
	storeCfg := store.DDBStoreConfig{
		Client:    dynamoClient,
		TableName: "mediapoints-metadata",
		Keys:      []string{"PK", "SK"},
	}

	initDynamoTables(dynamoClient, storeCfg)

	store := store.New[Data](storeCfg)

	return testStage{
		t:     t,
		store: store,
	}
}

func initDynamoTables(dynamoClient *dynamodb.Client, cfg store.DDBStoreConfig) {
	ctx := context.Background()

	// Table: mediapoints-metadata
	_, err := dynamoClient.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: &cfg.TableName,
		KeySchema: []types.KeySchemaElement{
			{KeyType: types.KeyTypeHash, AttributeName: aws.String("PK")},
			{KeyType: types.KeyTypeRange, AttributeName: aws.String("SK")},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("PK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("SK"), AttributeType: types.ScalarAttributeTypeS},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(100),
			WriteCapacityUnits: aws.Int64(100),
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("PK-SK-Index"),
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
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(100),
					WriteCapacityUnits: aws.Int64(100),
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeInclude,
					NonKeyAttributes: []string{
						"Updated",
					},
				},
			},
		},
	})
	mustBeNilErr(err)
}

func mustBeNilErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
