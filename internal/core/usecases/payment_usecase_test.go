package usecases

import (
	"errors"
	"testing"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	drivers "github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/payment"
	mock_payment "github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/payment/mocks"
	mock_gateways "github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/gateways/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPaymentUseCase_CreatePaymentOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	paymentBroker := mock_payment.NewMockPaymentBroker(ctrl)
	paymentRepository := mock_gateways.NewMockPaymentRepositoryGateway(ctrl)

	type args struct {
		paymentOrder dto.PaymentOrderDTO
	}
	type want struct {
		qrCode string
		err    error
	}
	type paymentBrokerCall struct {
		paymentOrder  dto.PaymentOrderDTO
		times         int
		paymentQRCode drivers.PaymentQRCodeResponse
		err           error
	}
	type paymentRepositoryCall struct {
		paymentOrder dto.PaymentOrderDTO
		qrCode       string
		times        int
		err          error
	}
	tests := []struct {
		name string
		args
		want
		paymentBrokerCall
		paymentRepositoryCall
	}{
		{
			name: "should fail to create payment order when payment broker returns error",
			args: args{
				paymentOrder: createPaymentOrderDTO(),
			},
			want: want{
				qrCode: "",
				err:    errors.New("internal server error"),
			},
			paymentBrokerCall: paymentBrokerCall{
				paymentOrder:  createPaymentOrderDTO(),
				times:         1,
				paymentQRCode: drivers.PaymentQRCodeResponse{},
				err:           errors.New("internal server error"),
			},
		},
		{
			name: "should fail to create payment order when payment repository returns error",
			args: args{
				paymentOrder: createPaymentOrderDTO(),
			},
			want: want{
				qrCode: "",
				err:    errors.New("internal server error"),
			},
			paymentBrokerCall: paymentBrokerCall{
				paymentOrder: createPaymentOrderDTO(),
				times:        1,
				paymentQRCode: drivers.PaymentQRCodeResponse{
					QrData:       "mercadopago123456",
					StoreOrderId: "98765",
				},
				err: nil,
			},
			paymentRepositoryCall: paymentRepositoryCall{
				paymentOrder: createPaymentOrderDTO(),
				qrCode:       "mercadopago123456",
				times:        1,
				err:          errors.New("internal server error"),
			},
		},
		{
			name: "should create payment order successfully",
			args: args{
				paymentOrder: createPaymentOrderDTO(),
			},
			want: want{
				qrCode: "mercadopago123456",
				err:    nil,
			},
			paymentBrokerCall: paymentBrokerCall{
				paymentOrder: createPaymentOrderDTO(),
				times:        1,
				paymentQRCode: drivers.PaymentQRCodeResponse{
					QrData:       "mercadopago123456",
					StoreOrderId: "98765",
				},
				err: nil,
			},
			paymentRepositoryCall: paymentRepositoryCall{
				paymentOrder: createPaymentOrderDTO(),
				qrCode:       "mercadopago123456",
				times:        1,
				err:          nil,
			},
		},
	}

	for _, tt := range tests {
		paymentBroker.EXPECT().
			GeneratePaymentQRCode(gomock.Eq(tt.paymentBrokerCall.paymentOrder)).
			Times(tt.paymentBrokerCall.times).
			Return(tt.paymentBrokerCall.paymentQRCode, tt.paymentBrokerCall.err)

		paymentRepository.EXPECT().
			SavePaymentOrder(gomock.Eq(tt.paymentRepositoryCall.paymentOrder), gomock.Eq(tt.paymentRepositoryCall.qrCode)).
			Times(tt.paymentRepositoryCall.times).
			Return(tt.paymentRepositoryCall.err)

		config := PaymentUseCaseConfig{
			PaymentBroker:     paymentBroker,
			PaymentRepository: paymentRepository,
			OrderClient:       nil,
		}
		paymentUseCase := NewPaymentUseCase(config)

		qrCode, err := paymentUseCase.CreatePaymentOrder(tt.args.paymentOrder)

		assert.Equal(t, tt.want.qrCode, qrCode)
		assert.Equal(t, tt.want.err, err)
	}
}

func TestPaymentUseCase_NotifyPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	paymentRepository := mock_gateways.NewMockPaymentRepositoryGateway(ctrl)
	orderClient := mock_gateways.NewMockOrderClient(ctrl)

	type args struct {
		orderId   int
		paymentId int
	}
	type want struct {
		err error
	}
	type orderClientCall struct {
		orderId int
		times   int
		err     error
	}
	type paymentRepositoryCall struct {
		orderId   int
		paymentId int
		times     int
		err       error
	}
	tests := []struct {
		name string
		args
		want
		orderClientCall
		paymentRepositoryCall
	}{
		{
			name: "should fail to notify payment when payment repository returns error",
			args: args{
				orderId:   123,
				paymentId: 111,
			},
			want: want{
				err: errors.New("internal server error"),
			},
			paymentRepositoryCall: paymentRepositoryCall{
				orderId:   123,
				paymentId: 111,
				times:     1,
				err:       errors.New("internal server error"),
			},
		},
		{
			name: "should fail to notify payment when order client returns error",
			args: args{
				orderId:   123,
				paymentId: 111,
			},
			want: want{
				err: errors.New("internal server error"),
			},
			paymentRepositoryCall: paymentRepositoryCall{
				orderId:   123,
				paymentId: 111,
				times:     1,
				err:       nil,
			},
			orderClientCall: orderClientCall{
				orderId: 123,
				times:   1,
				err:     errors.New("internal server error"),
			},
		},
		{
			name: "should notify payment successfully",
			args: args{
				orderId:   123,
				paymentId: 111,
			},
			want: want{
				err: nil,
			},
			paymentRepositoryCall: paymentRepositoryCall{
				orderId:   123,
				paymentId: 111,
				times:     1,
				err:       nil,
			},
			orderClientCall: orderClientCall{
				orderId: 123,
				times:   1,
				err:     nil,
			},
		},
	}

	for _, tt := range tests {
		paymentRepository.EXPECT().
			UpdatePaymentOrderStatus(gomock.Eq(tt.paymentRepositoryCall.orderId), gomock.Eq(tt.paymentRepositoryCall.paymentId), gomock.Eq(entities.PaymentStatusPaid)).
			Times(tt.paymentRepositoryCall.times).
			Return(tt.paymentRepositoryCall.err)

		orderClient.EXPECT().
			NotifyPaymentOrder(gomock.Eq(tt.orderClientCall.orderId), gomock.Eq(entities.PaymentStatusPaid)).
			Times(tt.orderClientCall.times).
			Return(tt.orderClientCall.err)

		config := PaymentUseCaseConfig{
			PaymentRepository: paymentRepository,
			OrderClient:       orderClient,
		}
		paymentUseCase := NewPaymentUseCase(config)

		err := paymentUseCase.NotifyPayment(tt.args.orderId, tt.args.paymentId)

		assert.Equal(t, tt.want.err, err)
	}
}

func createPaymentOrderDTO() dto.PaymentOrderDTO {
	return dto.PaymentOrderDTO{
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
	}
}
