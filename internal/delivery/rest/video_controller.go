package rest

import (
	"context"
	"net/http"
	"vrs-api/internal/dto"
	"vrs-api/internal/entity"

	"github.com/gin-gonic/gin"
)

type VideoUsecase interface {
	CreateVideo(ctx context.Context, video *entity.Video) error
	GetVideos(ctx context.Context) (entity.Videos, error)
}

type VideoController struct {
	vuc VideoUsecase
}

func NewVideoController(router *gin.Engine, vuc VideoUsecase) *VideoController {
	return &VideoController{vuc}
}

func (vc *VideoController) CreateVideo(ctx *gin.Context) {
	var payload dto.CreateVideoReq
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	data := entity.Video{
		Title:      payload.Title,
		Overview:   payload.Overview,
		Format:     payload.Format,
		TotalStock: payload.TotalStock,
		// initially available stock = total stock
		AvailableStock:    payload.TotalStock,
		CoverPath:         payload.CoverPath,
		ProductionCompany: payload.ProductionCompany,
	}

	if err := vc.vuc.CreateVideo(ctx, &data); err != nil {
		ctx.Error(err)
		return
	}

	res := dto.CreateVideoRes{
		ID:                data.ID,
		Title:             data.Title,
		Overview:          data.Overview,
		Format:            data.Format,
		TotalStock:        data.TotalStock,
		AvailableStock:    data.TotalStock,
		CoverPath:         data.CoverPath,
		ProductionCompany: data.ProductionCompany,
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data:    res,
	})
}

func (vc *VideoController) GetVideos(ctx *gin.Context) {
	datas, err := vc.vuc.GetVideos(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	var videos dto.GetVideosRes
	for _, data := range datas {
		video := dto.VideosRes{
			ID:                data.ID,
			Title:             data.Title,
			Overview:          data.Overview,
			Format:            data.Format,
			TotalStock:        data.TotalStock,
			AvailableStock:    data.AvailableStock,
			CoverPath:         data.CoverPath,
			ProductionCompany: data.ProductionCompany,
		}
		videos = append(videos, video)
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    videos,
	})
}
