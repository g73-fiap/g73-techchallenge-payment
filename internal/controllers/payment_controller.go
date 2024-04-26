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
	paymentUsecase usecases.PaymentUsecase
}

func NewPaymentController() PaymentController {
	return PaymentController{}
}

func (c PaymentController) PaymentHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		handleBadRequestResponse(ctx, "[id] path parameter is required", errors.New("id is missing"))
		return
	}

	orderId, err := strconv.Atoi(id)
	if err != nil {
		handleBadRequestResponse(ctx, "[id] path parameter is invalid", err)
		return
	}

	var paymentNotification dto.PaymentNotificationDTO
	err = ctx.ShouldBindJSON(&paymentNotification)
	if err != nil {
		handleBadRequestResponse(ctx, "failed to bind payment notification payload", err)
		return
	}

	valid, err := paymentNotification.ValidatePaymentNotification()
	if !valid {
		handleBadRequestResponse(ctx, "invalid payment notification payload", err)
		return
	}

	err = c.paymentUsecase.CreateOrderPayment(orderId)
	if err != nil {
		handleInternalServerResponse(ctx, "failed to handle payment", err)
		return
	}

	ctx.Status(http.StatusOK)
}
