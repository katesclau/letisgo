package store

// TODO: Implement search
// TODO: Implement soft delete

import (
	"context"
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"
)

type Interface[T any] interface {
	Get(ctx context.Context, keys []string) (*T, error)
	Update(ctx context.Context, keys []string, item T) (*T, error)
	Create(ctx context.Context, item T) error
	Delete(ctx context.Context, keys []string) (*T, error)
	Upsert(ctx context.Context, keys []string, item T) error
	Search(ctx context.Context,
		filterExpression string,
		expressionAttributeValues map[string]types.AttributeValue,
		limit int32,
		startAt map[string]types.AttributeValue,
	) (*SearchResult[T], error)
}

type DDBStoreConfig struct {
	Client    *dynamodb.Client
	TableName string
	Keys      []string
}

type DDBStore[T any] struct {
	config DDBStoreConfig
}

type SearchResult[T any] struct {
	Items            []T
	LastEvaluatedKey map[string]*types.AttributeValue
}

func New[T any](cfg DDBStoreConfig) Interface[T] {
	return &DDBStore[T]{
		config: cfg,
	}
}

func (s *DDBStore[T]) Get(ctx context.Context, keys []string) (*T, error) {
	var item T
	input := &dynamodb.GetItemInput{
		TableName: aws.String(s.config.TableName),
		Key:       s.getKeysMap(keys),
	}

	res, err := s.config.Client.GetItem(ctx, input)
	if err != nil {
		return &item, fmt.Errorf("[DDBStore] failed to get item: %w", err)
	}

	err = attributevalue.UnmarshalMap(res.Item, &item)
	if err != nil {
		return &item, fmt.Errorf("[DDBStore] failed to unmarshal item: %w", err)
	}

	return &item, nil
}

func (s *DDBStore[T]) Update(ctx context.Context, keys []string, item T) (*T, error) {
	var updatedItem T

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return &updatedItem, fmt.Errorf("[DDBStore] failed to marshal item: %w", err)
	}

	updateBuilder := expression.UpdateBuilder{}
	for key, value := range av {
		if slices.Contains(s.config.Keys, key) {
			continue
		}
		var v any
		err := attributevalue.Unmarshal(value, &v)
		if err != nil {
			return &updatedItem, fmt.Errorf("[DDBStore] failed to unmarshal value for key %s: %w", key, err)
		}
		logrus.Debugf("Setting key: %s, value: %v", key, v)
		updateBuilder = updateBuilder.Set(expression.Name(key), expression.Value(v))
	}

	expr, err := expression.NewBuilder().WithUpdate(updateBuilder).Build()
	if err != nil {
		return &updatedItem, fmt.Errorf("[DDBStore] failed to build update expression: %w", err)
	}

	res, err := s.config.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(s.config.TableName),
		Key:                       s.getKeysMap(keys),
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              types.ReturnValueAllNew,
	})
	if err != nil {
		return &updatedItem, fmt.Errorf("failed to update item: %w", err)
	}

	err = attributevalue.UnmarshalMap(res.Attributes, &updatedItem)
	if err != nil {
		return &updatedItem, fmt.Errorf("[DDBStore] failed to unmarshal item: %w", err)
	}

	return &updatedItem, nil
}

func (s *DDBStore[T]) Create(ctx context.Context, item T) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("[DDBStore] failed to marshal item: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"av": av,
	}).Trace("[DDBStore] Create.Input")
	input := &dynamodb.PutItemInput{
		TableName:    aws.String(s.config.TableName),
		Item:         av,
		ReturnValues: types.ReturnValueAllOld,
	}

	_, err = s.config.Client.PutItem(ctx, input)
	if err != nil {
		return fmt.Errorf("[DDBStore] failed to create item: %w", err)
	}

	return nil
}

func (s *DDBStore[T]) SoftDelete(ctx context.Context, keys []string) (*T, error) {
	var deletedItem T
	input := &dynamodb.DeleteItemInput{
		TableName:    aws.String(s.config.TableName),
		Key:          s.getKeysMap(keys),
		ReturnValues: types.ReturnValueAllOld,
	}

	res, err := s.config.Client.DeleteItem(ctx, input)
	if err != nil {
		return &deletedItem, err
	}

	err = attributevalue.UnmarshalMap(res.Attributes, &deletedItem)
	if err != nil {
		return &deletedItem, err
	}

	return &deletedItem, err
}

func (s *DDBStore[T]) Delete(ctx context.Context, keys []string) (*T, error) {
	var deletedItem T
	input := &dynamodb.DeleteItemInput{
		TableName:    aws.String(s.config.TableName),
		Key:          s.getKeysMap(keys),
		ReturnValues: types.ReturnValueAllOld,
	}

	res, err := s.config.Client.DeleteItem(ctx, input)
	if err != nil {
		return &deletedItem, err
	}

	err = attributevalue.UnmarshalMap(res.Attributes, &deletedItem)
	if err != nil {
		return &deletedItem, err
	}

	return &deletedItem, err
}

func (s *DDBStore[T]) Upsert(ctx context.Context, keys []string, item T) error {
	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %w", err)
	}

	_, err = s.config.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(s.config.TableName),
		Item:                av,
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

func (s *DDBStore[T]) Search(
	ctx context.Context,
	filterExpression string,
	expressionAttributeValues map[string]types.AttributeValue,
	limit int32,
	startAt map[string]types.AttributeValue,
) (*SearchResult[T], error) {
	var result SearchResult[T]

	input := &dynamodb.ScanInput{
		TableName:                 aws.String(s.config.TableName),
		FilterExpression:          aws.String(filterExpression),
		ExpressionAttributeValues: expressionAttributeValues,
		Limit:                     aws.Int32(limit),
	}
	if startAt != nil {
		input.ExclusiveStartKey = startAt
	}

	res, err := s.config.Client.Scan(ctx, input)
	if err != nil {
		return &result, err
	}

	var items []T
	err = attributevalue.UnmarshalListOfMaps(res.Items, &items)
	if err != nil {
		return &result, err
	}

	result.Items = items
	return &result, nil
}

func (s *DDBStore[T]) getKeysMap(keys []string) map[string]types.AttributeValue {
	var keysMap = make(map[string]types.AttributeValue)
	for i, key := range keys {
		keysMap[s.config.Keys[i]] = &types.AttributeValueMemberS{
			Value: key,
		}
	}
	return keysMap
}
