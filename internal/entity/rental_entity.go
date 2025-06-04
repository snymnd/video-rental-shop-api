package entity

import (
	"time"
	"vrs-api/internal/constant"
)

type (
	RentVideoParam struct {
		VideosID []int
		UserID   string
	}

	RentVideoReturn struct {
		PaymentID   int
		ExpiredTime time.Time
		TotalPrice  float64
		Videos      Videos
	}

	ReturnVideoParam struct {
		VideoIDs []int
		UserID   string
	}

	ReturnVideoReturn struct {
		LateFeePaymentID *int
		ExpiredTime      *time.Time
		TotalPrice       *float64
		LateRentals      []LateRental
	}

	Rental struct {
		ID               int
		VideoID          int
		UserID           string
		RentalPaymentID  int
		LateFeePaymentID *int
		Status           constant.RentalStatus
		DueDate          time.Time
		ReturnDate       *time.Time
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}
	Rentals []Rental

	LateRental struct {
		RentalID int
		DaysLate int
		DueDate  time.Time
		Video    Video
		LateFee  float64
	}

	MultipleRentParams struct {
		VideosID  []int
		UserID    string
		DueDate   time.Time
		PaymentID int
	}
)
