package gateways

import (
	"errors"
	"net/http"
	"testing"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	mock_http "github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/http/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOrderClient_NotifyPaymentOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_http.NewMockHttpClient(ctrl)

	type args struct {
		orderId int
		status  entities.PaymentStatus
	}
	type want struct {
		err error
	}
	type clientCall struct {
		orderApiUrl string
		times       int
		response    *http.Response
		err         error
	}
	tests := []struct {
		name string
		args
		want
		clientCall
	}{
		{
			name: "should fail to notify payment when http client returns error",
			args: args{
				orderId: 123,
				status:  entities.PaymentStatusPending,
			},
			want: want{
				errors.New("failed to call order api, error: internal error"),
			},
			clientCall: clientCall{
				orderApiUrl: "/order/123456/status",
				times:       1,
				response:    &http.Response{},
				err:         errors.New("internal error"),
			},
		},
		{
			name: "should fail to notify payment when response is non-2xx",
			args: args{
				orderId: 123,
				status:  entities.PaymentStatusPending,
			},
			want: want{
				errors.New("failed to call order api, status [500] non-2xx"),
			},
			clientCall: clientCall{
				orderApiUrl: "/order/1234/status",
				times:       1,
				response: &http.Response{
					StatusCode: 500,
				},
				err: nil,
			},
		},
		{
			name: "should notify payment when response is 2xx",
			args: args{
				orderId: 123,
				status:  entities.PaymentStatusPending,
			},
			want: want{
				nil,
			},
			clientCall: clientCall{
				orderApiUrl: "/order/12345/status",
				times:       1,
				response: &http.Response{
					StatusCode: 200,
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		httpClient.EXPECT().DoPut(gomock.Eq(tt.clientCall.orderApiUrl), gomock.Any()).
			Times(tt.clientCall.times).
			Return(tt.clientCall.response, tt.clientCall.err)

		orderClient := NewOrderClient(httpClient, "/order")
		err := orderClient.NotifyPaymentOrder(tt.args.orderId, tt.args.status)

		assert.Equal(t, tt.want.err, err)
	}
}
