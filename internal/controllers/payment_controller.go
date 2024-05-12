package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases"
	"github.com/IgorRamosBR/g73-techchallenge-payment/internal/core/usecases/dto"
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentUsecase usecases.PaymentUseCase
}

func NewPaymentController(paymentUsecase usecases.PaymentUseCase) PaymentController {
	return PaymentController{
		paymentUsecase: paymentUsecase,
	}
}

func (p PaymentController) CreatePaymentOrderHandler(c *gin.Context) {
	var paymentOrder dto.PaymentOrderDTO
	err := c.ShouldBindJSON(&paymentOrder)
	if err != nil {
		handleBadRequestResponse(c, "failed to bind payment order payload", err)
		return
	}

	valid, err := paymentOrder.ValidatePaymentOrder()
	if !valid {
		handleBadRequestResponse(c, "invalid payment order payload", err)
		return
	}

	paymentQRCode, err := p.paymentUsecase.CreatePaymentOrder(paymentOrder)
	if err != nil {
		handleInternalServerResponse(c, "failed to create payment order", err)
		return
	}

	c.JSON(http.StatusOK, dto.PaymentQRCode{QRCode: paymentQRCode})
}

func (p PaymentController) NotifyPaymentHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		handleBadRequestResponse(c, "[id] path parameter is required", errors.New("id is missing"))
		return
	}

	orderId, err := strconv.Atoi(id)
	if err != nil {
		handleBadRequestResponse(c, "[id] path parameter is invalid", err)
		return
	}

	var paymentNotification dto.PaymentNotificationDTO
	err = c.ShouldBindJSON(&paymentNotification)
	if err != nil {
		handleBadRequestResponse(c, "failed to bind payment notification payload", err)
		return
	}

	err = p.paymentUsecase.NotifyPayment(orderId, paymentNotification.PaymentId)
	if err != nil {
		handleInternalServerResponse(c, "failed to notify payment", err)
		return
	}

	c.Status(http.StatusOK)
}
