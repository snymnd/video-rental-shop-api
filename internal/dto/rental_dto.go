package dto

import "time"

type (
	RentVideosReq struct {
		VideosID []int `json:"videos_id" binding:"required,dive,gt=0"`
	}
	RentVideosRes struct {
		PaymentID   int        `json:"payment_id"`
		ExpiredTime time.Time  `json:"expired_time"`
		TotalPrice  float64    `json:"total_price"`
		Videos      []VideoRes `json:"videos_id"`
	}
	ReturnVideosRes struct {
		LateFeePaymentID *int         `json:"late_fee_payment_id,omitempty"`
		ExpiredTime      *time.Time   `json:"expired_time,omitempty"`
		TotalPrice       *float64     `json:"total_price,omitempty"`
		LateRentals      []LateRental `json:"late_rentals"`
	}
	LateRental struct {
		RentalID int       `json:"rental_id"`
		DueDate  time.Time `json:"due_date"`
		DaysLate int       `json:"days_late"`
		Videos   VideoRes  `json:"videos_id"`
		LateFee  float64   `json:"late_fee"`
	}
	ReturnVideosReq struct {
		UserID   string `json:"user_id" binding:"required"`
		VideosID []int  `json:"videos_id" binding:"required,dive,gt=0"`
	}
)
