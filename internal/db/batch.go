package db

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type KeyValues struct {
	Pk string `json:"pk" dynamodbav:"pk"`
	Sk string `json:"sk" dynamodbav:"sk"`
}

func (kv *KeyValues) AsAttributeValues() (map[string]types.AttributeValue, error) {
	av, err := attributevalue.MarshalMap(kv)
	if err != nil {
		return nil, err
	}
	return av, nil
}
