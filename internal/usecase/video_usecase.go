package usecase

import (
	"context"
	"vrs-api/internal/entity"
)

type VideoRepository interface {
	Create(ctx context.Context, video *entity.Video) error
	FetchAll(ctx context.Context) (entity.Videos, error)
}

type VideoUsecase struct {
	vr VideoRepository
}

func NewVideoUsecase(vr VideoRepository) *VideoUsecase {
	return &VideoUsecase{vr}
}

func (vu *VideoUsecase) CreateVideo(ctx context.Context, video *entity.Video) error {
	if err := vu.vr.Create(ctx, video); err != nil {
		return err
	}

	return nil
}

func (vu *VideoUsecase) GetVideos(ctx context.Context) (videos entity.Videos, err error) {
	videos, err = vu.vr.FetchAll(ctx)
	if err != nil {
		return videos, err
	}

	return videos, nil
}
