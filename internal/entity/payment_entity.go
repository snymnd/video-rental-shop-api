package entity

import (
	"time"
	"vrs-api/internal/constant"
)

type (
	Payment struct {
		ID          int
		UserID      string
		TotalPrice  float64
		Method      *constant.PaymentMethod
		ExpiredTime time.Time
		Status      constant.PaymentStatus
		CreatedAt   time.Time
	}

	UpdatePaymentParams struct {
		ID          int
		TotalPrice  *float64
		Method      *constant.PaymentMethod
		ExpiredTime *time.Time
		Status      *constant.PaymentStatus
	}
)
