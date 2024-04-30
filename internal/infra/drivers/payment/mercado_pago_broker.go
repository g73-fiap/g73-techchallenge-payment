package payment

import (
	"encoding/json"
	"fmt"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/http"
)

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

func (b mercadoPagoBroker) GeneratePaymentQRCode(paymentRequest PaymentRequest) (PaymentQRCodeResponse, error) {
	reqBody, err := json.Marshal(&paymentRequest)
	if err != nil {
		return PaymentQRCodeResponse{}, fmt.Errorf("failed to marshal payment qrcode request, error: %v", err)
	}

	response, err := b.httpClient.DoPost(b.brokerPath, reqBody)
	if err != nil {
		return PaymentQRCodeResponse{}, fmt.Errorf("failed to call mercado pago broker, error: %v", err)
	}
	defer response.Body.Close()

	var paymentQRCodeResponse PaymentQRCodeResponse
	err = json.NewDecoder(response.Body).Decode(&paymentQRCodeResponse)
	if err != nil {
		return PaymentQRCodeResponse{}, fmt.Errorf("failed to decode mercado pago response, error: %v", err)
	}

	return paymentQRCodeResponse, nil
}
