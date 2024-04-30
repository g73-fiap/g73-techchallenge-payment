package entities

type PaymentRequest struct {
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

type PaymentQRCodeResponse struct {
	QrData       string `json:"qr_data"`
	StoreOrderId string `json:"in_store_order_id"`
}
