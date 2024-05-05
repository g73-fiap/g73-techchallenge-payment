package usecases

import (
	"fmt"
	"strconv"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/entities"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/payment"
	drivers "github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/drivers/payment"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/infra/gateways"

	log "github.com/sirupsen/logrus"
)

type PaymentUsecase interface {
	CreatePaymentOrder(paymentOrder dto.PaymentOrderDTO) (string, error)
	NotifyPayment(orderId int) error
}

type paymentUsecase struct {
	notificationUrl   string
	sponsorId         string
	paymentBroker     drivers.PaymentBroker
	paymentRepository gateways.PaymentRepositoryGateway
	orderClient       gateways.OrderClient
}

func NewPaymentUsecase(notificationUrl, sponsorId string, paymentBroker payment.PaymentBroker) paymentUsecase {
	return paymentUsecase{
		notificationUrl: notificationUrl,
		sponsorId:       sponsorId,
		paymentBroker:   paymentBroker,
	}
}

func (u paymentUsecase) CreatePaymentOrder(paymentOrder dto.PaymentOrderDTO) (string, error) {
	paymentRequest := u.createPaymentRequest(paymentOrder)
	paymentQRCode, err := u.paymentBroker.GeneratePaymentQRCode(paymentRequest)
	if err != nil {
		log.Errorf("failed to generate payment qrcode for the order [%d], error: %v", paymentOrder.OrderId, err)
		return "", err
	}

	err = u.paymentRepository.SavePaymentOrder(paymentOrder, paymentQRCode.QrData)
	if err != nil {
		log.Errorf("failed to save payment order [%d], error: %v", paymentOrder.OrderId, err)
		return "", err
	}

	return paymentQRCode.QrData, err
}

func (u paymentUsecase) NotifyPayment(orderId int) error {
	err := u.paymentRepository.UpdatePaymentOrderStatus(orderId, entities.PaymentStatusPaid)
	if err != nil {
		log.Errorf("failed to payment payment status for the order [%d], error: %v", orderId, err)
		return err
	}

	err = u.orderClient.NotifyPaymentOrder(orderId, entities.PaymentStatusPaid)
	if err != nil {
		log.Errorf("failed to notify payment order for the order [%d], error: %v", orderId, err)
		return err
	}

	return nil
}

func (u paymentUsecase) createPaymentRequest(paymentOrder dto.PaymentOrderDTO) payment.PaymentRequest {
	var items []payment.PaymentItemRequest
	for _, item := range paymentOrder.Items {
		items = append(items, createPaymentItem(item))
	}

	return payment.PaymentRequest{
		ExternalReference: strconv.FormatUint(uint64(paymentOrder.OrderId), 10),
		Title:             fmt.Sprintf("Order %d for the Customer[%s]", paymentOrder.OrderId, paymentOrder.CustomerCPF),
		NotificationURL:   fmt.Sprintf("%s/payment/%d/notify", u.notificationUrl, paymentOrder.OrderId),
		TotalAmount:       paymentOrder.TotalAmount,
		Items:             items,
		Sponsor:           u.sponsorId,
	}
}

func createPaymentItem(item dto.PaymentOrderItem) payment.PaymentItemRequest {
	paymentItem := payment.PaymentItemRequest{
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
