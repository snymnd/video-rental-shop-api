package usecase

import (
	"context"
	"errors"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/entity"
	"vrs-api/internal/util/logger"
)

type (
	VideoRepository interface {
		Create(ctx context.Context, video *entity.Video) error
		FetchAll(ctx context.Context, params entity.GetVideosParams) (entity.GetVideosReturn, error)
	}
	VideoCacheRepository interface {
		FetchAll(ctx context.Context, params entity.GetVideosParams) (entity.GetVideosReturn, error)
		SetFetchAll(ctx context.Context, params entity.GetVideosParams, value entity.GetVideosReturn) error
		GetFetchAllKey(params entity.GetVideosParams) string
	}

	videoUsecase struct {
		vr  VideoRepository
		vcr VideoCacheRepository
	}
)

func NewVideoUsecase(vr VideoRepository, vcr VideoCacheRepository) *videoUsecase {
	return &videoUsecase{vr, vcr}
}

func (vu *videoUsecase) CreateVideo(ctx context.Context, video *entity.Video) error {
	if err := vu.vr.Create(ctx, video); err != nil {
		return err
	}

	return nil
}

func (vu *videoUsecase) GetVideos(ctx context.Context, params entity.GetVideosParams) (videos entity.GetVideosReturn, err error) {
	log := logger.GetLogger()
	videos, err = vu.vcr.FetchAll(ctx, params)
	if err != nil && !errors.Is(err, customerrors.ErrCacheKeyNotFound) {
		log.Errorf("failed to get videos from cache, error: %s\n", err.Error())
	}

	// If cache hit, return cached data
	if err == nil {
		cacheKey := vu.vcr.GetFetchAllKey(params)
		log.Infof("use videos data from redis cache key %s", cacheKey)
		return videos, nil
	}

	videos, err = vu.vr.FetchAll(ctx, params)
	if err != nil {
		return videos, err
	}
	if err := vu.vcr.SetFetchAll(ctx, params, videos); err != nil {
		log.Errorf("failed to set videos to cache, error: %s\n", err.Error())
	}

	return videos, nil
}
