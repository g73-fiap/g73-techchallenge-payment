package payment

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/http"
)

type mercadoPagoBroker struct {
	httpClient      http.HttpClient
	brokerPath      string
	notificationUrl string
	sponsorId       string
}

type MercadoPagoBrokerConfig struct {
	HttpClient      http.HttpClient
	BrokerUrl       string
	NotificationUrl string
	SponsorId       string
}

func NewMercadoPagoBroker(config MercadoPagoBrokerConfig) PaymentBroker {
	return mercadoPagoBroker{
		httpClient:      config.HttpClient,
		brokerPath:      config.BrokerUrl,
		notificationUrl: config.NotificationUrl,
		sponsorId:       config.SponsorId,
	}
}

func (b mercadoPagoBroker) GeneratePaymentQRCode(paymentOrder dto.PaymentOrderDTO) (PaymentQRCodeResponse, error) {
	paymentRequest := b.createPaymentRequest(paymentOrder)

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

func (b mercadoPagoBroker) createPaymentRequest(paymentOrder dto.PaymentOrderDTO) PaymentRequest {
	var items []PaymentItemRequest
	for _, item := range paymentOrder.Items {
		items = append(items, createPaymentItem(item))
	}

	return PaymentRequest{
		ExternalReference: strconv.FormatUint(uint64(paymentOrder.OrderId), 10),
		Title:             fmt.Sprintf("Order %d for the Customer[%s]", paymentOrder.OrderId, paymentOrder.CustomerCPF),
		NotificationURL:   fmt.Sprintf("%s/payment/%d/notify", b.notificationUrl, paymentOrder.OrderId),
		TotalAmount:       paymentOrder.TotalAmount,
		Items:             items,
		Sponsor:           b.sponsorId,
	}
}

func createPaymentItem(item dto.PaymentOrderItem) PaymentItemRequest {
	paymentItem := PaymentItemRequest{
		SkuNumber:   item.Product.SkuId,
		Category:    item.Product.Category,
		Title:       item.Product.Name,
		Description: item.Product.Description,
		UnitPrice:   item.Product.Price,
		Quantity:    item.Quantity,
		UnitMeasure: getUnitMeasure(item.Product.Type),
		TotalAmount: item.Product.Price * float64(item.Quantity),
	}

	return paymentItem
}

func getUnitMeasure(itemType string) string {
	if itemType == string(dto.OrderItemTypeCustomCombo) {
		return "pack"
	}
	return "unit"
}
