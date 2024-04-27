package usecases

import (
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/payment"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/gateways"

	log "github.com/sirupsen/logrus"
)

type PaymentUsecase interface {
	GeneratePaymentQRCode(order entities.Order) (string, error)
	CreateOrderPayment(orderId int) error
}

type paymentUsecase struct {
	notificationUrl        string
	sponsorId              string
	paymentBroker          payment.PaymentBroker
	orderRepositoryGateway gateways.OrderRepositoryGateway
}

func NewPaymentUsecase(notificationUrl, sponsorId string, paymentBroker payment.PaymentBroker) PaymentUsecase {
	return paymentUsecase{
		notificationUrl: notificationUrl,
		sponsorId:       sponsorId,
		paymentBroker:   paymentBroker,
	}
}

func (u paymentUsecase) GeneratePaymentQRCode(order entities.Order) (string, error) {
	paymentRequest := u.createPaymentRequest(order)
	paymentResponse, err := u.paymentBroker.GeneratePaymentQRCode(paymentRequest)
	if err != nil {
		log.Errorf("failed to generate payment qrcode for the order [%d], error: %v", order.ID, err)
		return "", err
	}

	return paymentResponse.QrData, nil
}

func (u paymentUsecase) createPaymentRequest(order entities.Order) dto.PaymentQRCodeRequest {
	var items []dto.PaymentItemRequest
	for _, item := range order.Items {
		items = append(items, createPaymentItem(item))
	}

	return dto.PaymentQRCodeRequest{
		OrderId:     order.ID,
		Items:       items,
		TotalAmount: order.TotalAmount,
	}
}

func createPaymentItem(item entities.OrderItem) dto.PaymentItemRequest {
	paymentItem := dto.PaymentItemRequest{
		Quantity: item.Quantity,
		Product: dto.PaymentProductRequest{
			Name:        item.Product.Name,
			SkuId:       item.Product.SkuId,
			Description: item.Product.Description,
			Category:    item.Product.Category,
			Type:        item.Type,
			Price:       item.Product.Price,
		},
	}

	return paymentItem
}

func getUnitMeasure(itemType string) string {
	if itemType == string(dto.OrderItemTypeCustomCombo) {
		return "pack"
	}
	return "unit"
}

func (u paymentUsecase) CreateOrderPayment(orderId int) error {
	err := u.orderRepositoryGateway.UpdateOrderStatus(orderId, string(dto.OrderStatusPaid))
	if err != nil {
		log.Errorf("failed to update order status from order id [%d], error: %v", orderId, err)
		return err
	}

	return nil
}
