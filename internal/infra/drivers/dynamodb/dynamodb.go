package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	client *dynamodb.Client
}

func NewDynamoDBClient(client *dynamodb.Client) *DynamoDBClient {
	return &DynamoDBClient{client: client}
}

func (d *DynamoDBClient) PutItem(tableName string, item map[string]types.AttributeValue) error {
	_, err := d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *DynamoDBClient) GetItem(tableName string, key map[string]types.AttributeValue) (map[string]types.AttributeValue, error) {
	result, err := d.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       key,
	})
	if err != nil {
		return nil, err
	}
	return result.Item, nil
}
