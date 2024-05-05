package payment

import "github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"

type PaymentBroker interface {
	GeneratePaymentQRCode(paymentOrder dto.PaymentOrderDTO) (PaymentQRCodeResponse, error)
}
type PaymentRequest struct {
	ExternalReference string               `json:"external_reference"`
	Title             string               `json:"title"`
	NotificationURL   string               `json:"notification_url"`
	TotalAmount       float64              `json:"total_amount"`
	Items             []PaymentItemRequest `json:"items"`
	Sponsor           string               `json:"sponsor"`
}

type PaymentItemRequest struct {
	SkuNumber   string  `json:"sku_number"`
	Category    string  `json:"category"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
	UnitMeasure string  `json:"unit_measure"`
	TotalAmount float64 `json:"total_amount"`
}

type PaymentQRCodeResponse struct {
	QrData       string `json:"qr_data"`
	StoreOrderId string `json:"in_store_order_id"`
}
