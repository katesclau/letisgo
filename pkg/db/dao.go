package db

import (
    "context"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoDBHandler struct {
    svc *dynamodb.DynamoDB
    tableName string
}

func NewDynamoDBHandler(sess *session.Session, tableName string) *DynamoDBHandler {
    return &DynamoDBHandler{
        svc: dynamodb.New(sess),
        tableName: tableName,
    }
}

func (h *DynamoDBHandler) Query(ctx context.Context, input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
    return h.svc.QueryWithContext(ctx, input)
}

func (h *DynamoDBHandler) Put(ctx context.Context, item interface{}) (*dynamodb.PutItemOutput, error) {
    av, err := dynamodbattribute.MarshalMap(item)
    if err != nil {
        return nil, err
    }
    input := &dynamodb.PutItemInput{
        TableName: aws.String(h.tableName),
        Item: av,
    }
    return h.svc.PutItemWithContext(ctx, input)
}

func (h *DynamoDBHandler) Scan(ctx context.Context, input *dynamodb.ScanInput) (*dynamodb.ScanOutput, error) {
    return h.svc.ScanWithContext(ctx, input)
}

func (h *DynamoDBHandler) Get(ctx context.Context, key map[string]*dynamodb.AttributeValue) (*dynamodb.GetItemOutput, error) {
    input := &dynamodb.GetItemInput{
        TableName: aws.String(h.tableName),
        Key: key,
    }
    return h.svc.GetItemWithContext(ctx, input)
}

func (h *DynamoDBHandler) Update(ctx context.Context, key map[string]*dynamodb.AttributeValue, updateExpression string, expressionAttributeValues map[string]*dynamodb.AttributeValue) (*dynamodb.UpdateItemOutput, error) {
    input := &dynamodb.UpdateItemInput{
        TableName: aws.String(h.tableName),
        Key: key,
        UpdateExpression: aws.String(updateExpression),
        ExpressionAttributeValues: expressionAttributeValues,
    }
    return h.svc.UpdateItemWithContext(ctx, input)
}

func (h *DynamoDBHandler) Insert(ctx context.Context, item interface{}) (*dynamodb.PutItemOutput, error) {
    return h.Put(ctx, item)
}
