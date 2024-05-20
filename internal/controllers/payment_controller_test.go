package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	mock_usecases "github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var paymentRequestValid, _ = os.ReadFile("./testdata/payment_request_valid.json")
var paymentRequestInvalid, _ = os.ReadFile("./testdata/payment_request_invalid.json")

func TestPaymentController_CreatePaymentOrderHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	paymentUseCase := mock_usecases.NewMockPaymentUseCase(ctrl)
	paymentController := NewPaymentController(paymentUseCase)

	type args struct {
		reqBody string
	}
	type want struct {
		statusCode int
		respBody   string
	}
	type paymentUseCaseCall struct {
		paymentOrder dto.PaymentOrderDTO
		times        int
		qrCode       string
		err          error
	}
	tests := []struct {
		name string
		args
		want
		paymentUseCaseCall
	}{
		{
			name: "should return bad request when req body is not a json",
			args: args{
				reqBody: "<invalidJson>",
			},
			want: want{
				statusCode: 404,
				respBody:   `{"message":"failed to bind payment order payload","error":"invalid character '\u003c' looking for beginning of value"}`,
			},
		},
		{
			name: "should return bad request when req body is invalid",
			args: args{
				reqBody: string(paymentRequestInvalid),
			},
			want: want{
				statusCode: 400,
				respBody:   `{"message":"invalid payment order payload","error":"Customer CPF is required"}`,
			},
		},
		{
			name: "should return internal server error when payment use case fails to create payment order",
			args: args{
				reqBody: string(paymentRequestValid),
			},
			want: want{
				statusCode: 500,
				respBody:   `{"message":"failed to create payment order","error":"internal server error"}`,
			},
			paymentUseCaseCall: paymentUseCaseCall{
				paymentOrder: createPaymentOrder(),
				times:        1,
				qrCode:       "",
				err:          errors.New("internal server error"),
			},
		},
		{
			name: "should return ok when creates payment order successfully",
			args: args{
				reqBody: string(paymentRequestValid),
			},
			want: want{
				statusCode: 200,
				respBody:   `{"qrcode":"mercadopago123456"}`,
			},
			paymentUseCaseCall: paymentUseCaseCall{
				paymentOrder: createPaymentOrder(),
				times:        1,
				qrCode:       "mercadopago123456",
				err:          nil,
			},
		},
	}

	for _, tt := range tests {
		paymentUseCase.EXPECT().
			CreatePaymentOrder(gomock.Eq(tt.paymentUseCaseCall.paymentOrder)).
			Times(tt.paymentUseCaseCall.times).
			Return(tt.paymentUseCaseCall.qrCode, tt.paymentUseCaseCall.err)

		router := createRouter(paymentController)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/v1/paymentOrder", strings.NewReader(tt.args.reqBody))
		router.ServeHTTP(w, req)

		assert.Equal(t, tt.want.statusCode, w.Code)
		assert.Equal(t, tt.want.respBody, w.Body.String())
	}
}

func TestPaymentController_NotifyPaymentHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	paymentUseCase := mock_usecases.NewMockPaymentUseCase(ctrl)
	paymentController := NewPaymentController(paymentUseCase)

	type args struct {
		id      string
		reqBody string
	}
	type want struct {
		statusCode int
		respBody   string
	}
	type paymentUseCaseCall struct {
		orderId   int
		paymentId int
		times     int
		err       error
	}
	tests := []struct {
		name string
		args
		want
		paymentUseCaseCall
	}{
		{
			name: "should return bad request when orderId is empty",
			args: args{
				id:      "",
				reqBody: "",
			},
			want: want{
				statusCode: 400,
				respBody:   `{"message":"[id] path parameter is required","error":"id is missing"}`,
			},
		},
		{
			name: "should return bad request when orderId is not a number",
			args: args{
				id:      "abc",
				reqBody: "",
			},
			want: want{
				statusCode: 400,
				respBody:   `{"message":"[id] path parameter is invalid","error":"strconv.Atoi: parsing \"abc\": invalid syntax"}`,
			},
		},
		{
			name: "should return bad request when req body is not a json",
			args: args{
				id:      "123",
				reqBody: "<invalidJson>",
			},
			want: want{
				statusCode: 400,
				respBody:   `{"message":"failed to bind payment notification payload","error":"invalid character '\u003c' looking for beginning of value"}`,
			},
		},

		{
			name: "should return internal server error when payment use case fails to notify payment",
			args: args{
				id: "123",
				reqBody: `{
					"description": "Payment received for order 123",
					"merchant_order": 123456,
					"payment_id": 7890
				}`,
			},
			want: want{
				statusCode: 500,
				respBody:   `{"message":"failed to notify payment","error":"internal server error"}`,
			},
			paymentUseCaseCall: paymentUseCaseCall{
				orderId:   123,
				paymentId: 7890,
				times:     1,
				err:       errors.New("internal server error"),
			},
		},
		{
			name: "should return ok when creates payment order successfully",
			args: args{
				id: "123",
				reqBody: `{
					"description": "Payment received for order 123",
					"merchant_order": 123456,
					"payment_id": 7890
				}`,
			},
			want: want{
				statusCode: 200,
			},
			paymentUseCaseCall: paymentUseCaseCall{
				orderId:   123,
				paymentId: 7890,
				times:     1,
				err:       nil,
			},
		},
	}

	for _, tt := range tests {
		paymentUseCase.EXPECT().
			NotifyPayment(gomock.Eq(tt.paymentUseCaseCall.orderId), gomock.Eq(tt.paymentUseCaseCall.paymentId)).
			Times(tt.paymentUseCaseCall.times).
			Return(tt.paymentUseCaseCall.err)

		router := createRouter(paymentController)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/payment/%s/notify", tt.args.id), strings.NewReader(tt.args.reqBody))
		router.ServeHTTP(w, req)

		assert.Equal(t, tt.want.statusCode, w.Code)
		assert.Equal(t, tt.want.respBody, w.Body.String())
	}
}

func createRouter(paymenteControler PaymentController) *gin.Engine {

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/payment/:id/notify", paymenteControler.NotifyPaymentHandler)
		v1.POST("/paymentOrder", paymenteControler.CreatePaymentOrderHandler)
	}
	return router
}

func createPaymentOrder() dto.PaymentOrderDTO {
	return dto.PaymentOrderDTO{
		OrderId:     123456,
		CustomerCPF: "123.456.789-00",
		Items: []dto.PaymentOrderItem{
			{
				Quantity: 2,
				Product: dto.OrderItemProduct{
					Name:        "Product A",
					SkuId:       "SKU123",
					Description: "Description of Product A",
					Category:    "Category1",
					Type:        string(dto.OrderItemTypeUnit),
					Price:       19.99,
				},
			},
			{
				Quantity: 1,
				Product: dto.OrderItemProduct{
					Name:        "Product B",
					SkuId:       "SKU124",
					Description: "Description of Product B",
					Category:    "Category2",
					Type:        string(dto.OrderItemTypeCombo),
					Price:       49.99,
				},
			},
		},
		TotalAmount: 89.97,
	}
}
