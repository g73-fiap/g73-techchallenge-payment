package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient interface {
	GetItem(tableName string, key map[string]types.AttributeValue) (map[string]types.AttributeValue, error)
	PutItem(tableName string, item map[string]types.AttributeValue) error
	UpdateItem(tableName string, key map[string]types.AttributeValue, expr expression.Expression) error
}

type dynamoDBClient struct {
	client *dynamodb.Client
}

func NewDynamoDBClient(client *dynamodb.Client) *dynamoDBClient {
	return &dynamoDBClient{client: client}
}

func (d *dynamoDBClient) PutItem(tableName string, item map[string]types.AttributeValue) error {
	_, err := d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *dynamoDBClient) GetItem(tableName string, key map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {
	result, err := d.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	return result.Item, nil
}

func (d *dynamoDBClient) UpdateItem(tableName string, key map[string]types.AttributeValue, expr expression.Expression) error {
	_, err := d.client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName:                 &tableName,
		Key:                       key,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})
	if err != nil {
		return err
	}
	return nil
}
