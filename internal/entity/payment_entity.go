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
		Method      *string
		ExpiredTime time.Time
		Status      constant.PaymentStatus
		CreatedAt   time.Time
	}
)
