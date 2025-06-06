package usecase

import (
	"context"
	"vrs-api/internal/entity"
)

type (
	VideoRepository interface {
		Create(ctx context.Context, video *entity.Video) error
		FetchAll(ctx context.Context, params entity.GetVideosParams) (entity.GetVideosReturn, error)
	}

	videoUsecase struct {
		vr VideoRepository
	}
)

func NewVideoUsecase(vr VideoRepository) *videoUsecase {
	return &videoUsecase{vr}
}

func (vu *videoUsecase) CreateVideo(ctx context.Context, video *entity.Video) error {
	if err := vu.vr.Create(ctx, video); err != nil {
		return err
	}

	return nil
}

func (vu *videoUsecase) GetVideos(ctx context.Context, params entity.GetVideosParams) (videos entity.GetVideosReturn, err error) {

	videos, err = vu.vr.FetchAll(ctx, params)
	if err != nil {
		return videos, err
	}

	return videos, nil
}
