package gateways

import (
	"errors"
	"testing"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestPaymentRepository_SavePaymentOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	dynamodbClient := mock_dynamodb.NewMockDynamoDBClient(ctrl)

	type args struct {
		paymentDto dto.PaymentOrderDTO
		qrCode     string
	}
	type want struct {
		err error
	}
	type dynamodbCall struct {
		table string
		times int
		err   error
	}
	tests := []struct {
		name string
		args
		want
		dynamodbCall
	}{
		{
			name: "should fail to save payment order when dynamodb client returns error",
			args: args{
				paymentDto: dto.PaymentOrderDTO{
					OrderId:     123,
					CustomerCPF: "111222333444",
					Items: []dto.PaymentOrderItem{
						{
							Quantity: 1,
							Product: dto.OrderItemProduct{
								Name:        "Batata frita",
								SkuId:       "333",
								Description: "Batata canoa",
								Category:    "Acompanhamento",
								Type:        "UNIT",
								Price:       9.99,
							},
						},
					},
					TotalAmount: 9.99,
				},
				qrCode: "mercadopago1234566778",
			},
			want: want{
				errors.New("internal error"),
			},
			dynamodbCall: dynamodbCall{
				table: "Payment",
				times: 1,
				err:   errors.New("internal error"),
			},
		},
		{
			name: "should save payment order when dynamodb client does not return error",
			args: args{
				paymentDto: dto.PaymentOrderDTO{
					OrderId:     123,
					CustomerCPF: "111222333444",
					Items: []dto.PaymentOrderItem{
						{
							Quantity: 1,
							Product: dto.OrderItemProduct{
								Name:        "Batata frita",
								SkuId:       "333",
								Description: "Batata canoa",
								Category:    "Acompanhamento",
								Type:        "UNIT",
								Price:       9.99,
							},
						},
					},
					TotalAmount: 9.99,
				},
				qrCode: "mercadopago1234566778",
			},
			want: want{
				nil,
			},
			dynamodbCall: dynamodbCall{
				table: "Payment",
				times: 1,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		dynamodbClient.EXPECT().PutItem(gomock.Eq(tt.dynamodbCall.table), gomock.Any()).
			Times(tt.dynamodbCall.times).
			Return(tt.dynamodbCall.err)

		paymentRepository := NewPaymentRepositoryGateway(dynamodbClient, "Payment")
		err := paymentRepository.SavePaymentOrder(tt.args.paymentDto, tt.args.qrCode)

		assert.Equal(t, tt.want.err, err)
	}
}

func TestPaymentRepository_UpdatePaymentOrderStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	dynamodbClient := mock_dynamodb.NewMockDynamoDBClient(ctrl)

	type args struct {
		orderId   int
		paymentId int
		status    entities.PaymentStatus
	}
	type want struct {
		err error
	}
	type dynamodbCall struct {
		table string
		times int
		err   error
	}
	tests := []struct {
		name string
		args
		want
		dynamodbCall
	}{
		{
			name: "should fail to update payment order status when dynamodb client returns error",
			args: args{
				orderId:   123,
				paymentId: 999,
				status:    entities.PaymentStatusPaid,
			},
			want: want{
				errors.New("internal error"),
			},
			dynamodbCall: dynamodbCall{
				table: "Payment",
				times: 1,
				err:   errors.New("internal error"),
			},
		},
		{
			name: "should update payment order when dynamodb client does not return error",
			args: args{
				orderId:   123,
				paymentId: 999,
				status:    entities.PaymentStatusPaid,
			},
			want: want{
				nil,
			},
			dynamodbCall: dynamodbCall{
				table: "Payment",
				times: 1,
				err:   nil,
			},
		},
	}

	for _, tt := range tests {
		dynamodbClient.EXPECT().UpdateItem(gomock.Eq(tt.dynamodbCall.table), gomock.Any(), gomock.Any()).
			Times(tt.dynamodbCall.times).
			Return(tt.dynamodbCall.err)

		paymentRepository := NewPaymentRepositoryGateway(dynamodbClient, "Payment")
		err := paymentRepository.UpdatePaymentOrderStatus(tt.args.orderId, tt.args.paymentId, tt.args.status)

		assert.Equal(t, tt.want.err, err)
	}
}
