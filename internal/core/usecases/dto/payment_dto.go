package dto

import (
	"time"

	"github.com/asaskevich/govalidator"
)

type PaymentQRCodeRequest struct {
	OrderId     int                  `json:"orderId"`
	Items       []PaymentItemRequest `json:"items"`
	TotalAmount float64              `json:"totalAmount"`
}

type PaymentItemRequest struct {
	Quantity int                   `json:"quantity"`
	Product  PaymentProductRequest `json:"product"`
}

type PaymentProductRequest struct {
	Name        string  `json:"name"`
	SkuId       string  `json:"skuId"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Type        string  `json:"type"`
	Price       float64 `json:"price"`
}

type SponsorRequest struct {
	Id string `json:"id"`
}

type PaymentQRCodeResponse struct {
	QrData       string `json:"qr_data"`
	StoreOrderId string `json:"in_store_order_id"`
}

type PaymentQRCode struct {
	QRCode string `json:"qrcode"`
}

type PaymentNotificationDTO struct {
	Id          string      `json:"id"`
	LiveMode    bool        `json:"liveMode"`
	Type        string      `json:"type" valid:"in(payment),required~Type is invalid"`
	DateCreated time.Time   `json:"dateCreated"`
	UserId      int         `json:"userId"`
	ApiVersion  string      `json:"apiVersion"`
	Action      string      `json:"action"`
	Data        PaymentData `json:"data"`
}

type PaymentData struct {
	Id string `json:"id" valid:"required,numeric"`
}

func (p PaymentNotificationDTO) ValidatePaymentNotification() (bool, error) {
	if _, err := govalidator.ValidateStruct(p); err != nil {
		return false, err
	}

	return true, nil
}
