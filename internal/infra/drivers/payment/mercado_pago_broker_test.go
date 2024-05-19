package payment

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	mock_http "github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/http/mocks"
	"github.com/go-playground/assert/v2"
	"go.uber.org/mock/gomock"
)

func TestMercadoPagoBroker_GeneratePaymentQRCode(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_http.NewMockHttpClient(ctrl)

	type args struct {
		paymentOrder dto.PaymentOrderDTO
	}
	type want struct {
		qrCodeResponse PaymentQRCodeResponse
		err            error
	}
	type clientCall struct {
		brokerPath string
		times      int
		response   *http.Response
		err        error
	}
	tests := []struct {
		name string
		args
		want
		clientCall
	}{
		{
			name: "should fail to generate payment qrcode when http client returns error",
			args: args{
				paymentOrder: dto.PaymentOrderDTO{
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
			},
			want: want{
				qrCodeResponse: PaymentQRCodeResponse{},
				err:            errors.New("failed to call mercado pago broker, error: internal error"),
			},
			clientCall: clientCall{
				brokerPath: "/mercadopago",
				times:      1,
				response:   &http.Response{},
				err:        errors.New("internal error"),
			},
		},
		{
			name: "should fail to generate payment qrcode when response is invalid",
			args: args{
				paymentOrder: dto.PaymentOrderDTO{
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
			},
			want: want{
				qrCodeResponse: PaymentQRCodeResponse{},
				err:            errors.New("failed to decode mercado pago response, error: invalid character '<' looking for beginning of value"),
			},
			clientCall: clientCall{
				brokerPath: "/mercadopago",
				times:      1,
				response: &http.Response{
					Body: io.NopCloser(strings.NewReader("<invalid json>")),
				},
				err: nil,
			},
		},
		{
			name: "should generate payment qrcode",
			args: args{
				paymentOrder: dto.PaymentOrderDTO{
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
			},
			want: want{
				qrCodeResponse: PaymentQRCodeResponse{
					QrData:       "mercadopago123456789",
					StoreOrderId: 9876,
				},
				err: nil,
			},
			clientCall: clientCall{
				brokerPath: "/mercadopago",
				times:      1,
				response: &http.Response{
					StatusCode: 200,
					Body:       io.NopCloser(strings.NewReader(`{"qr_data":"mercadopago123456789","in_store_order_id":"9876"}`)),
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		httpClient.EXPECT().DoPost(gomock.Eq(tt.clientCall.brokerPath), gomock.Any()).
			Times(tt.clientCall.times).
			Return(tt.clientCall.response, tt.clientCall.err)

		config := MercadoPagoBrokerConfig{
			HttpClient:      httpClient,
			BrokerUrl:       "/mercadopago",
			NotificationUrl: "/notification",
			SponsorId:       "3333",
		}
		mercadoPagoBroker := NewMercadoPagoBroker(config)
		mercadoPagoResponse, err := mercadoPagoBroker.GeneratePaymentQRCode(tt.args.paymentOrder)

		assert.Equal(t, tt.want.qrCodeResponse, mercadoPagoResponse)
		assert.Equal(t, tt.want.err, err)
	}

}
