// https://adrianhesketh.com/2020/04/17/single-table-pattern-dynamodb-with-go-part-1/
package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
)

type DynamoDBHandler[T any] interface {
	Insert(ctx context.Context, item any) (*dynamodb.PutItemOutput, error)
	Get(ctx context.Context, pk string, sk string) (T, error)
	Delete(ctx context.Context, pk string, sk string) (T, error)
	BatchGet(ctx context.Context, keys []KeyValues) ([]T, error)
	// Update(ctx context.Context, pk string, sk string, updateExpression string, expressionValues map[string]types.AttributeValue) (T, error)
	// Scan(ctx context.Context, filterExpression string, expressionValues map[string]types.AttributeValue) ([]T, error)
}

type dynamoDBHandler[T any] struct {
	client    *dynamodb.Client
	tableName string
}

var (
	ErrDuplicatedID = errors.New("duplicated id")
)

const (
	pk = "pk"
	sk = "sk"
)

func NewDynamoDBHandler[T any](ddb *dynamodb.Client, tableName string) DynamoDBHandler[T] {
	return &dynamoDBHandler[T]{
		client:    ddb,
		tableName: tableName,
	}
}

func (h *dynamoDBHandler[T]) Insert(ctx context.Context, item any) (*dynamodb.PutItemOutput, error) {
	recordProvider, ok := item.(RecordInterface)
	if !ok {
		return nil, errors.New("item does not implement Record() method")
	}
	record := recordProvider.Record()

	av, err := attributevalue.MarshalMap(record)
	if err != nil {
		return nil, err
	}

	// This conditional expression is dependant on the relation type of the model being inserted
	cond := expression.AttributeNotExists(expression.Name(sk))
	expr, err := expression.NewBuilder().WithCondition(cond).Build()
	if err != nil {
		return nil, err
	}

	input := dynamodb.PutItemInput{
		Item:                     av,
		TableName:                &h.tableName,
		ConditionExpression:      expr.Condition(),
		ExpressionAttributeNames: expr.Names(),
	}

	resp, err := h.client.PutItem(ctx, &input)
	var ccf *types.ConditionalCheckFailedException
	if errors.As(err, &ccf) {
		logrus.Debug("faield on conditional check", err)
		return nil, ErrDuplicatedID
	}
	if err != nil {
		logrus.Debug("failed on put item action", err)
		return nil, err
	}

	return resp, nil
}

func (h *dynamoDBHandler[T]) Get(ctx context.Context, part string, rang string) (T, error) {
	var ret T
	keyCond := expression.Key(pk).Equal(expression.Value(part)).And(expression.Key(sk).Equal(expression.Value(rang)))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return ret, err
	}

	input := &dynamodb.QueryInput{
		TableName:                 &h.tableName,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	output, err := h.client.Query(ctx, input)
	if err != nil {
		logrus.Debug("failed query", err)
		return ret, err
	}

	if len(output.Items) == 0 {
		logrus.Debug("not found", err)
		return ret, errors.New("record not found")
	}

	err = attributevalue.UnmarshalMap(output.Items[0], &ret)
	if err != nil {
		logrus.Debug("failed to unmarshal", err)
		return ret, errors.New("record not found")

	}

	return ret, nil
}

func (h *dynamoDBHandler[T]) Delete(ctx context.Context, part string, rang string) (T, error) {
	var ret T
	key := map[string]types.AttributeValue{
		pk: &types.AttributeValueMemberS{Value: part},
		sk: &types.AttributeValueMemberS{Value: rang},
	}

	input := &dynamodb.DeleteItemInput{
		TableName:    &h.tableName,
		Key:          key,
		ReturnValues: types.ReturnValueAllOld,
	}

	output, err := h.client.DeleteItem(ctx, input)
	if err != nil {
		logrus.Debug("failed delete", err)
		return ret, err
	}

	if output.Attributes == nil {
		logrus.Debug("not found", err)
		return ret, errors.New("record not found")
	}

	err = attributevalue.UnmarshalMap(output.Attributes, &ret)
	if err != nil {
		logrus.Debug("failed to unmarshal", err)
		return ret, errors.New("failed to unmarshal deleted record")
	}

	return ret, nil
}

func (h *dynamoDBHandler[T]) BatchGet(ctx context.Context, keys []KeyValues) ([]T, error) {
	var results []T

	if len(keys) == 0 {
		return results, nil
	}

	avKeys := make([]map[string]types.AttributeValue, 0, len(keys))
	for _, key := range keys {
		keyMap, err := attributevalue.MarshalMap(key)
		if err != nil {
			logrus.Debug("failed to marshal key", err)
			// return nil, err
		}
		avKeys = append(avKeys, keyMap)
	}

	batchInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			h.tableName: {
				Keys: avKeys,
			},
		},
	}

	output, err := h.client.BatchGetItem(ctx, batchInput)
	if err != nil {
		logrus.Debug("failed batch get", err)
		return nil, err
	}

	items, ok := output.Responses[h.tableName]
	if !ok {
		logrus.Debug("no items found in batch get response")
		return results, nil
	}

	results = make([]T, 0, len(keys))
	for _, item := range items {
		var record T
		err := attributevalue.UnmarshalMap(item, &record)
		if err != nil {
			logrus.Warn("failed to unmarshal item", err)
			// return nil, err
		}
		results = append(results, record)
	}

	return results, nil
}
