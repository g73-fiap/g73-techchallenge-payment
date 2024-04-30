package controllers

import (
	"net/http"

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
	var paymentOrder dto.PaymentOrder
	err := ctx.ShouldBindJSON(&paymentOrder)
	if err != nil {
		handleBadRequestResponse(ctx, "failed to bind payment order payload", err)
		return
	}

	valid, err := paymentOrder.ValidatePaymentOrder()
	if !valid {
		handleBadRequestResponse(ctx, "invalid payment order payload", err)
		return
	}

	err = c.paymentUsecase.CreatePaymentOrder(paymentOrder)
	if err != nil {
		handleInternalServerResponse(ctx, "failed to handle payment", err)
		return
	}

	ctx.Status(http.StatusOK)
}
