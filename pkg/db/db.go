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
	"mnesis.com/pkg/log"
)

type DynamoDBHandler interface {
	Insert(ctx context.Context, item any) (*dynamodb.PutItemOutput, error)
}

type dynamoDBHandler struct {
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

func NewDynamoDBHandler(ddb *dynamodb.Client, tableName string) DynamoDBHandler {
	return &dynamoDBHandler{
		client:    ddb,
		tableName: tableName,
	}
}

func (h *dynamoDBHandler) Insert(ctx context.Context, item any) (*dynamodb.PutItemOutput, error) {
	recordProvider, ok := item.(RecordInterface)
	if !ok {
		return nil, errors.New("item does not implement Record() method")
	}
	record := recordProvider.Record()

	logrus.WithFields(logrus.Fields{
		"record": log.GetJsonString(record),
	}).Info("Record in JSON format")
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

	logrus.WithFields(logrus.Fields{
		"item":  item,
		"table": h.tableName,
		"av":    av,
	}).Info("Values...")
	input := dynamodb.PutItemInput{
		Item:                     av,
		TableName:                &h.tableName,
		ConditionExpression:      expr.Condition(),
		ExpressionAttributeNames: expr.Names(),
	}
	logrus.WithFields(logrus.Fields{
		"input": log.GetJsonString(input),
	}).Info("Input")

	resp, err := h.client.PutItem(ctx, &input)
	var ccf *types.ConditionalCheckFailedException
	if errors.As(err, &ccf) {
		return nil, ErrDuplicatedID
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}
