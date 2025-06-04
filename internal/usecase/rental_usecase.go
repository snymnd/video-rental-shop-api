package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"vrs-api/internal/constant"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/dto"
	"vrs-api/internal/entity"
)

type (
	RentalRepository interface {
		Creates(ctx context.Context, rentals entity.MultipleRentParams) error
		FetchMultipleRentals(ctx context.Context, videosID []int, userID string, status constant.RentalStatus) (entity.Rentals, error)
		UpdatesAddLatefee(ctx context.Context, rentalIDs []int, lateFeePaymentId int) error
		UpdatesRentalStatus(ctx context.Context, rentalIDs []int, status constant.RentalStatus) error
	}
	RentalVideoRepository interface {
		FetchMultipleVideos(ctx context.Context, videosID []int) (entity.Videos, error)
		RentMultipleVideos(ctx context.Context, videosID []int) error
		ReturnMultipleVideos(ctx context.Context, videosID []int) error
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
					Title:   strconv.Itoa(video.ID),
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
			UserID:    rentVideosParams.UserID,
			VideosID:  rentVideosParams.VideosID,
			DueDate:   time.Now().Local().Add(time.Hour * 24 * constant.DEFAULT_RENTAL_DUE),
			PaymentID: payment.ID,
		}
		if createRentalsErr := ru.rr.Creates(txCtx, rentals); createRentalsErr != nil {
			return createRentalsErr
		}

		// update (decrease) videos stock
		if err := ru.vr.RentMultipleVideos(txCtx, rentVideosParams.VideosID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return rentReturn, err
	}

	return rentReturn, nil
}

func (ru *RentalUsecase) ReturnVideos(ctx context.Context, renturnVideosParams entity.ReturnVideoParam) (returnVideosReturn entity.ReturnVideoReturn, err error) {
	if err = ru.txr.WithTx(ctx, func(txCtx context.Context) error {
		// fetch rentals based on video id and user id
		rentals, fetchVideoErr := ru.rr.FetchMultipleRentals(txCtx, renturnVideosParams.VideoIDs, renturnVideosParams.UserID, constant.RENTAL_RENTED)
		if fetchVideoErr != nil {
			return fetchVideoErr
		}
		if len(rentals) <= 0 {
			return customerrors.NewError(
				"No rentals found for the specified user and videos",
				errors.New("no rental records found for the provided userID and videoIDs"),
				customerrors.ItemNotExist,
			)
		}

		videos, fetchVideoErr := ru.vr.FetchMultipleVideos(txCtx, renturnVideosParams.VideoIDs)
		if fetchVideoErr != nil {
			return fetchVideoErr
		}
		videoIDVideoMap := make(map[int]entity.Video)
		for _, video := range videos {
			videoIDVideoMap[video.ID] = video
		}

		var totalLateFee float64
		lateRentalIDs := make([]int, 0)
		rentalIDs := make([]int, 0)
		lateRentals := make([]entity.LateRental, 0)
		// check rentals due date and calculate any late fee
		for _, rental := range rentals {
			// HACK: use time.Now but only change the time zone (default: time.Local) to time.UTC+0
			// this allow us to compare with postgres time return with UTC+0 timezone
			now := time.Now()
			today := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.UTC)

			rentalIDs = append(rentalIDs, rental.ID)

			fmt.Println(today, rental.DueDate, today.Sub(rental.DueDate).Hours(), "check")
			if today.After(rental.DueDate) {
				daysLate := int(today.Sub(rental.DueDate).Hours() / 24)
				if daysLate >= 1 {
					lateFee := float64(daysLate) * constant.LATE_FEE

					if lateFee > constant.MAX_CAP_FEE {
						lateFee = constant.MAX_CAP_FEE
					}

					video := videoIDVideoMap[rental.VideoID]
					lateRentals = append(lateRentals, entity.LateRental{
						Video:    video,
						RentalID: rental.ID,
						DaysLate: daysLate,
						DueDate:  rental.DueDate,
						LateFee:  lateFee,
					})
					lateRentalIDs = append(lateRentalIDs, rental.ID)

					totalLateFee += lateFee
				}
			}
		}
		returnVideosReturn.LateRentals = lateRentals

		// update rental status
		if err := ru.rr.UpdatesRentalStatus(txCtx, rentalIDs, constant.RENTAL_RETURNED); err != nil {
			return err
		}
		// update (increase) videos stock
		if err := ru.vr.ReturnMultipleVideos(txCtx, renturnVideosParams.VideoIDs); err != nil {
			return err
		}
		// return if there are no late rentals
		if len(lateRentalIDs) <= 0 {
			return nil
		}

		// if there are late rentals
		returnVideosReturn.TotalPrice = &totalLateFee
		// create payment
		paymentExpiredTime := time.Now().Add(time.Hour * constant.DEFAULT_PAYMENT_EXPIRED_DUE)
		payment := entity.Payment{
			UserID:      renturnVideosParams.UserID,
			TotalPrice:  totalLateFee,
			ExpiredTime: paymentExpiredTime,
		}
		if createPaymentErr := ru.pr.Create(txCtx, &payment); createPaymentErr != nil {
			return createPaymentErr
		}
		returnVideosReturn.LateFeePaymentID = &payment.ID
		returnVideosReturn.ExpiredTime = &paymentExpiredTime

		// update rental with adding late fee payment id
		if err := ru.rr.UpdatesAddLatefee(txCtx, lateRentalIDs, payment.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return returnVideosReturn, err
	}

	return returnVideosReturn, nil
}
