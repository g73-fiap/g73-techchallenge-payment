package payment

import (
	"encoding/json"
	"fmt"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/http"
)

type PaymentBroker interface {
	GeneratePaymentQRCode(dto.PaymentQRCodeRequest) (dto.PaymentQRCodeResponse, error)
}

type mercadoPagoBroker struct {
	httpClient http.HttpClient
	brokerPath string
}

func NewMercadoPagoBroker(httpClient http.HttpClient, brokerPath string) PaymentBroker {
	return mercadoPagoBroker{
		httpClient: httpClient,
		brokerPath: brokerPath,
	}
}

func (b mercadoPagoBroker) GeneratePaymentQRCode(request dto.PaymentQRCodeRequest) (dto.PaymentQRCodeResponse, error) {
	reqBody, err := json.Marshal(&request)
	if err != nil {
		return dto.PaymentQRCodeResponse{}, fmt.Errorf("failed to marshal payment qrcode request, error: %v", err)
	}

	response, err := b.httpClient.DoPost(b.brokerPath, reqBody)
	if err != nil {
		return dto.PaymentQRCodeResponse{}, fmt.Errorf("failed to call mercado pago broker, error: %v", err)
	}
	defer response.Body.Close()

	var paymentQRCodeResponse dto.PaymentQRCodeResponse
	err = json.NewDecoder(response.Body).Decode(&paymentQRCodeResponse)
	if err != nil {
		return dto.PaymentQRCodeResponse{}, fmt.Errorf("failed to decode mercado pago response, error: %v", err)
	}

	return paymentQRCodeResponse, nil
}
