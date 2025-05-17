package rest

import (
	"context"
	"net/http"
	"vrs-api/internal/constant"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/dto"
	"vrs-api/internal/entity"
	util "vrs-api/internal/util/jwt"

	"github.com/gin-gonic/gin"
)

type RentalUsecase interface {
	RentVideos(ctx context.Context, rentVideos entity.RentVideoParam) (entity.RentVideoReturn, error)
}

type RentalController struct {
	ruc RentalUsecase
}

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
	authPayload := ctx.Value(constant.CTX_AUTH_PAYLOAD_KEY).(*util.JWTCustomClaims)
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
			// available stock is havent
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
