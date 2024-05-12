package gateways

import (
	"strconv"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/dynamodb"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type PaymentRepositoryGateway interface {
	SavePaymentOrder(paymentOrderDTO dto.PaymentOrderDTO, qrCode string) error
	UpdatePaymentOrderStatus(orderId, paymentId int, status entities.PaymentStatus) error
}

type paymentRepositoryGateway struct {
	paymentTable   string
	dynamodbClient dynamodb.DynamoDBClient
}

func NewPaymentRepositoryGateway(dynamodbClient dynamodb.DynamoDBClient, paymentTable string) PaymentRepositoryGateway {
	return paymentRepositoryGateway{
		dynamodbClient: dynamodbClient,
		paymentTable:   paymentTable,
	}
}

func (p paymentRepositoryGateway) SavePaymentOrder(paymentOrderDTO dto.PaymentOrderDTO, qrCode string) error {
	paymentOrder := paymentOrderDTO.ToPaymentOrder(qrCode)

	av, err := attributevalue.MarshalMap(paymentOrder)
	if err != nil {
		return err
	}

	err = p.dynamodbClient.PutItem(p.paymentTable, av)
	if err != nil {
		return err
	}

	return nil
}

func (p paymentRepositoryGateway) UpdatePaymentOrderStatus(orderId, paymentId int, status entities.PaymentStatus) error {
	key := map[string]types.AttributeValue{
		"OrderId": &types.AttributeValueMemberN{Value: strconv.Itoa(orderId)},
	}
	update := expression.Set(expression.Name("Status"), expression.Value(status))
	update.Set(expression.Name("PaymentId"), expression.Value(paymentId))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	err = p.dynamodbClient.UpdateItem(p.paymentTable, key, expr)
	if err != nil {
		return err
	}

	return nil
}
