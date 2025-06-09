package rest

import (
	"context"
	"errors"
	"net/http"
	"vrs-api/internal/constant"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/dto"
	"vrs-api/internal/entity"
	"vrs-api/internal/util/token"

	"github.com/gin-gonic/gin"
)

type (
	RentalUsecase interface {
		RentVideos(ctx context.Context, rentVideos entity.RentVideoParam) (entity.RentVideoReturn, error)
		ReturnVideos(ctx context.Context, rentVideos entity.ReturnVideoParam) (entity.ReturnVideoReturn, error)
	}

	RentalController struct {
		ruc RentalUsecase
	}
)

func NewRentalController(ruc RentalUsecase) *RentalController {
	return &RentalController{ruc}
}

func (rc *RentalController) RentVideos(ctx *gin.Context) {
	var payload dto.RentVideosReq
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	// get token payload and user id from token subject claims
	authPayload := ctx.Value(constant.CTX_AUTH_PAYLOAD_KEY).(*token.JWTCustomClaims)
	if authPayload == nil {
		ctx.Error(customerrors.ErrFailedGetAuthPayload)
		return
	}
	userID := authPayload.Subject
	if userID == "" {
		ctx.Error(customerrors.ErrUserIDNotFoundInAuthPayload)
		return
	}

	data := entity.RentVideoParam{
		VideosID: payload.VideosID,
		UserID:   userID,
	}

	rentVideos, err := rc.ruc.RentVideos(ctx, data)
	if err != nil {
		ctx.Error(err)
		return
	}

	var videosRes []dto.VideoRes
	for _, video := range rentVideos.Videos {
		video := dto.VideoRes{
			ID:         video.ID,
			Title:      video.Title,
			Overview:   video.Overview,
			Format:     video.Format,
			RentPrice:  video.RentPrice,
			TotalStock: video.TotalStock,
			// available stock is haven't being reduced by 1
			AvailableStock:    video.AvailableStock - 1,
			CoverPath:         video.CoverPath,
			ProductionCompany: video.ProductionCompany,
			GenreIDs:          video.GenreIDs,
		}
		videosRes = append(videosRes, video)
	}

	res := dto.RentVideosRes{
		PaymentID:   rentVideos.PaymentID,
		ExpiredTime: rentVideos.ExpiredTime,
		TotalPrice:  rentVideos.TotalPrice,
		Videos:      videosRes,
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data:    res,
	})

}

func (rc *RentalController) ReturnVideos(ctx *gin.Context) {
	var payload dto.ReturnVideosReq
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	if len(payload.VideosID) <= 0 {
		ctx.Error(customerrors.NewError(
			"Must return at least 1 video",
			errors.New("VideoIDs length must be longer than 1"),
			customerrors.InvalidAction,
		))
		return
	}

	data := entity.ReturnVideoParam{
		UserID:   payload.UserID,
		VideoIDs: payload.VideosID,
	}

	returnVideos, err := rc.ruc.ReturnVideos(ctx, data)
	if err != nil {
		ctx.Error(err)
		return
	}

	lateRentals := make([]dto.LateRental, 0)
	for _, lateRentalData := range returnVideos.LateRentals {

		video := dto.VideoRes{
			ID:                lateRentalData.Video.ID,
			Title:             lateRentalData.Video.Title,
			Overview:          lateRentalData.Video.Overview,
			Format:            lateRentalData.Video.Format,
			TotalStock:        lateRentalData.Video.TotalStock,
			AvailableStock:    lateRentalData.Video.AvailableStock,
			RentPrice:         lateRentalData.Video.RentPrice,
			CoverPath:         lateRentalData.Video.CoverPath,
			ProductionCompany: lateRentalData.Video.ProductionCompany,
			GenreIDs:          lateRentalData.Video.GenreIDs,
		}

		lateRental := dto.LateRental{
			RentalID: lateRentalData.RentalID,
			DaysLate: lateRentalData.DaysLate,
			Videos:   video,
			LateFee:  lateRentalData.LateFee,
			DueDate:  lateRentalData.DueDate,
		}
		lateRentals = append(lateRentals, lateRental)
	}

	res := dto.ReturnVideosRes{
		ExpiredTime:      returnVideos.ExpiredTime,
		LateFeePaymentID: returnVideos.LateFeePaymentID,
		TotalPrice:       returnVideos.TotalPrice,
		LateRentals:      lateRentals,
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    res,
	})

}
