package dto

import (
	"time"

	"github.com/asaskevich/govalidator"
)

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

type PaymentOrderStatusDTO struct {
	Status string `json:"status"`
}
