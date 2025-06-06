package rest

import (
	"context"
	"net/http"
	"strconv"
	"vrs-api/internal/dto"
	"vrs-api/internal/entity"

	"github.com/gin-gonic/gin"
)

type (
	VideoUsecase interface {
		CreateVideo(ctx context.Context, video *entity.Video) error
		GetVideos(ctx context.Context, params entity.GetVideosParams) (entity.GetVideosReturn, error)
	}

	VideoController struct {
		vuc VideoUsecase
	}
)

func NewVideoController(vuc VideoUsecase) *VideoController {
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
		RentPrice:  payload.RentPrice,
		// initially available stock = total stock
		AvailableStock:    payload.TotalStock,
		CoverPath:         payload.CoverPath,
		ProductionCompany: payload.ProductionCompany,
		GenreIDs:          payload.GenreIDs,
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
		RentPrice:         data.RentPrice,
		AvailableStock:    data.TotalStock,
		CoverPath:         data.CoverPath,
		ProductionCompany: data.ProductionCompany,
		GenreIDs:          data.GenreIDs,
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Data:    res,
	})
}

func (vc *VideoController) GetVideos(ctx *gin.Context) {
	var query dto.GetVideosQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.Error(err)
		return
	}

	limit, err := strconv.Atoi(query.Limit)
	if err != nil {
		defaultLimit := 0
		limit = defaultLimit
	}
	page, err := strconv.Atoi(query.Page)
	if err != nil || page <= 0 {
		defaultPage := 1
		page = defaultPage
	}

	params := entity.GetVideosParams{
		SortOrder: query.OrderSort,
		GenreIDs:  query.Genres,
		Title:     query.Title,
		OrderBy:   query.OrderBy,
		PaginationQuery: entity.PaginationQuery{
			Limit: limit,
			Page:  page,
		},
	}

	data, err := vc.vuc.GetVideos(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	videos := []dto.VideoRes{}
	for _, entry := range data.Entries {
		video := dto.VideoRes{
			ID:                entry.ID,
			Title:             entry.Title,
			Overview:          entry.Overview,
			Format:            entry.Format,
			TotalStock:        entry.TotalStock,
			RentPrice:         entry.RentPrice,
			AvailableStock:    entry.AvailableStock,
			CoverPath:         entry.CoverPath,
			ProductionCompany: entry.ProductionCompany,
			GenreIDs:          entry.GenreIDs,
			CreatedAt:         &entry.CreatedAt,
		}
		videos = append(videos, video)
	}

	var filters []dto.PageFilter
	for _, filter := range data.PageInfo.Filters {
		filters = append(filters, dto.PageFilter{
			Field: filter.Field,
			Value: filter.Value,
		})
	}

	getVideosReturn := dto.PaginatadResponse[dto.VideoRes]{
		Entries: videos,
		PageInfo: dto.PageInfo{
			Page:      data.PageInfo.Page,
			Limit:     data.PageInfo.Limit,
			OrderBy:   data.PageInfo.OrderBy,
			OrderSort: data.PageInfo.OrderSort,
			Filters:   filters,
			TotalRow:  data.PageInfo.TotalRow,
		},
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    getVideosReturn,
	})
}
