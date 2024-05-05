package entities

type PaymentStatus string

var (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusPaid    PaymentStatus = "PAID"
)

type PaymentOrder struct {
	OrderId     int           `dynamodbav:"OrderId"`
	CustomerCPF string        `dynamodbav:"CustomerCPF"`
	TotalAmout  float64       `dynamodbav:"TotalAmount"`
	Status      PaymentStatus `dynamodbav:"Status"`
	QRCode      string        `dynamodbav:"QRCode"`
}
