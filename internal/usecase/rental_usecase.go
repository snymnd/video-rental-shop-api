package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"vrs-api/internal/constant"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/dto"
	"vrs-api/internal/entity"
)

type (
	RentalRepository interface {
		Create(ctx context.Context, rentals entity.MultipleRentParams) error
	}
	RentalVideoRepository interface {
		FetchMultipleVideos(ctx context.Context, videosID []int) (entity.Videos, error)
		RentMultipleVideos(ctx context.Context, videosID []int) error
	}
	PaymentRepository interface {
		Create(ctx context.Context, payment *entity.Payment) error
	}
	RentTxRepository interface {
		WithTx(ctx context.Context, tFunc func(ctx context.Context) error) error
	}
)

type RentalUsecase struct {
	rr  RentalRepository
	vr  RentalVideoRepository
	pr  PaymentRepository
	txr RentTxRepository
}

func NewRentalUsecase(rr RentalRepository, vr RentalVideoRepository, pr PaymentRepository, tx RentTxRepository) *RentalUsecase {
	return &RentalUsecase{rr, vr, pr, tx}
}

func (ru *RentalUsecase) RentVideos(ctx context.Context, rentVideosParams entity.RentVideoParam) (rentReturn entity.RentVideoReturn, err error) {
	if err = ru.txr.WithTx(ctx, func(txCtx context.Context) error {
		// fetch needed videos
		videos, fetchVideoErr := ru.vr.FetchMultipleVideos(txCtx, rentVideosParams.VideosID)
		if fetchVideoErr != nil {
			return fetchVideoErr
		}

		// check available stock and count total price
		var totalPrice float64
		var unavailableVideosErr []dto.DetailsError
		for _, video := range videos {
			if video.AvailableStock < 1 {
				unavailableVideosErr = append(unavailableVideosErr, dto.DetailsError{
					Title:   video.ID,
					Message: fmt.Sprintf("%s is not available", video.Title),
				})
			}
			totalPrice += video.RentPrice
		}
		if len(unavailableVideosErr) > 0 {
			return customerrors.NewError(
				"Some videos are not available",
				errors.New("some videos are not available"),
				customerrors.InvalidAction,
				unavailableVideosErr,
			)
		}
		rentReturn.TotalPrice = totalPrice
		rentReturn.Videos = videos

		// create payment
		paymentExpiredTime := time.Now().Add(time.Hour * constant.DEFAULT_PAYMENT_EXPIRED_DUE)
		payment := entity.Payment{
			UserID:      rentVideosParams.UserID,
			TotalPrice:  totalPrice,
			ExpiredTime: paymentExpiredTime,
		}
		if createPaymentErr := ru.pr.Create(txCtx, &payment); createPaymentErr != nil {
			return createPaymentErr
		}
		rentReturn.PaymentID = payment.ID
		rentReturn.ExpiredTime = paymentExpiredTime

		// crate videos rental record
		rentals := entity.MultipleRentParams{
			VideosID:  rentVideosParams.VideosID,
			DueDate:   time.Now().Add(time.Hour * 24 * constant.DEFAULT_RENTAL_DUE),
			PaymentID: payment.ID,
		}
		if createRentalsErr := ru.rr.Create(txCtx, rentals); createRentalsErr != nil {
			return createRentalsErr
		}

		// update (decrease) videos stock
		if err := ru.vr.RentMultipleVideos(txCtx, rentVideosParams.VideosID); err != nil {
			return err
		}
		// rentReturn.Videos = videos

		return nil
	}); err != nil {
		return rentReturn, err
	}

	return rentReturn, nil
}
