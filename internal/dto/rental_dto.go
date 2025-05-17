package dto

import "time"

type (
	RentVideosReq struct {
		VideosID []int `json:"videos_id"`
	}
	RentVideosRes struct {
		PaymentID   int        `json:"payment_id"`
		ExpiredTime time.Time  `json:"expired_time"`
		TotalPrice  float64    `json:"total_price"`
		Videos      []VideoRes `json:"videos_id"`
	}
)
