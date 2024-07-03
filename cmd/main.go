package main

import (
	"context"

	"github.com/IgorRamosBR/g73-techchallenge-payment/configs"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/api"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/controllers"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/dynamodb"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/http"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/payment"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/gateways"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsDynamoDb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config := configs.NewConfig()
	appConfig, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	// mercado pago payment broker
	paymentHttpClient := http.NewMockHttpClient()
	paymentBrokerConfig := payment.MercadoPagoBrokerConfig{
		HttpClient:      paymentHttpClient,
		BrokerUrl:       appConfig.PaymentBrokerURL,
		NotificationUrl: appConfig.NotificationURL,
		SponsorId:       appConfig.SponsorId,
	}
	paymentBroker := payment.NewMercadoPagoBroker(paymentBrokerConfig)

	// payment repository
	dynamodbClient, err := NewDynamoDBClient(appConfig.Environment, appConfig.PaymentTableEndpoint)
	if err != nil {
		panic(err)
	}
	paymentRepository := gateways.NewPaymentRepositoryGateway(dynamodbClient, appConfig.PaymentTable)

	// order api
	httpClient := http.NewHttpClient(appConfig.DefaultTimeout)
	orderClient := gateways.NewOrderClient(httpClient, appConfig.OrderApiUrl)

	// payment usecase
	paymentUseCaseConfig := usecases.PaymentUseCaseConfig{
		PaymentBroker:     paymentBroker,
		PaymentRepository: paymentRepository,
		OrderClient:       orderClient,
	}
	paymentUseCase := usecases.NewPaymentUseCase(paymentUseCaseConfig)

	// payment controller
	paymentController := controllers.NewPaymentController(paymentUseCase)

	api := api.NewApi(paymentController)
	api.Run(":" + appConfig.Port)
}

func NewDynamoDBClient(environment, endpoint string) (dynamodb.DynamoDBClient, error) {
	if environment == "local" {
		return createLocalDynamoDb(endpoint)
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	client := awsDynamoDb.NewFromConfig(cfg)
	return dynamodb.NewDynamoDBClient(client), nil
}

func createLocalDynamoDb(endpoint string) (dynamodb.DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "ANYKEYID", SecretAccessKey: "ANYSECRETKEY",
				Source: "Values not relevant in shardDb mode",
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	client := awsDynamoDb.NewFromConfig(cfg)
	output, err := client.CreateTable(context.Background(), &awsDynamoDb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("OrderId"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("OrderId"),
				KeyType:       types.KeyType(types.KeyTypeHash),
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String("Payment"),
	})
	if err != nil {
		logrus.Warnf("failed to create Payment table, error: %s", err.Error())
	}

	if output != nil {
		logrus.Info("table Payment created.")
	}

	return dynamodb.NewDynamoDBClient(client), nil
}
