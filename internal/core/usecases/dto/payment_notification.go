package dto

type PaymentNotificationDTO struct {
	Description   string `json:"description"`
	MerchantOrder int    `json:"merchant_order"`
	PaymentId     int    `json:"payment_id"`
}

type PaymentOrderStatusDTO struct {
	Status string `json:"status"`
}
