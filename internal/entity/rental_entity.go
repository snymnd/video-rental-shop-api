package entity

import "time"

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

	Rental struct {
		ID         int
		VideoID    int
		UserID     string
		PaymentID  int
		DueDate    time.Time
		ReturnDate *time.Time
		CreatedAt  time.Time
		UpdatedAt  time.Time
	}

	MultipleRentParams struct {
		VideosID  []int
		DueDate   time.Time
		PaymentID int
	}
)
