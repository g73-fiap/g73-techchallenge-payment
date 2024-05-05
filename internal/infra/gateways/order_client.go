package gateways

import (
	"encoding/json"
	"fmt"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/http"
)

type OrderClient interface {
	NotifyPaymentOrder(orderId int, status entities.PaymentStatus) error
}

type orderClient struct {
	httpClient  http.HttpClient
	orderApiUrl string
}

func NewOrderClient(httpClient http.HttpClient, orderApiUrl string) OrderClient {
	return orderClient{
		httpClient:  httpClient,
		orderApiUrl: orderApiUrl,
	}
}

func (o orderClient) NotifyPaymentOrder(orderId int, status entities.PaymentStatus) error {
	reqBody, err := json.Marshal(dto.PaymentOrderStatusDTO{
		Status: string(status),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payment qrcode request, error: %v", err)
	}

	response, err := o.httpClient.DoPut(fmt.Sprintf("%s/%d/status", o.orderApiUrl, orderId), reqBody)
	if err != nil {
		return fmt.Errorf("failed to call order api, error: %v", err)
	}

	if response.StatusCode > 299 || response.StatusCode < 200 {
		return fmt.Errorf("failed to call order api, status [%d] non-2xx", response.StatusCode)
	}

	return nil
}
