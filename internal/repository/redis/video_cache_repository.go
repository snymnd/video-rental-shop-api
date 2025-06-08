package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/entity"

	"github.com/go-redis/redis/v8"
)

type videoCacheRepository struct {
	redisClient *redis.Client
}

func NewVideoCacheRepository(redisClient *redis.Client) *videoCacheRepository {
	return &videoCacheRepository{redisClient}
}

func (vcr *videoCacheRepository) FetchAll(ctx context.Context, params entity.GetVideosParams) (entity.GetVideosReturn, error) {
	cacheKey := vcr.GetFetchAllKey(params)
	videosBytes, err := vcr.redisClient.Get(ctx, cacheKey).Bytes()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return entity.GetVideosReturn{}, customerrors.ErrCacheKeyNotFound
		}
		return entity.GetVideosReturn{}, err
	}
	videosReturn := entity.GetVideosReturn{}
	if err = json.Unmarshal(videosBytes, &videosReturn); err != nil {
		return entity.GetVideosReturn{}, err
	}

	return videosReturn, nil
}

const cacheDuration = 86400 // 1 day
func (vcr *videoCacheRepository) SetFetchAll(ctx context.Context, params entity.GetVideosParams, value entity.GetVideosReturn) error {
	cacheKey := vcr.GetFetchAllKey(params)
	fmt.Println(cacheKey)

	valueBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := vcr.redisClient.Set(ctx, cacheKey, valueBytes, time.Second*time.Duration(cacheDuration)).Err(); err != nil {
		return err
	}

	return nil
}

const videoCacheKeyPrefx = "video"

func (vcr *videoCacheRepository) GetFetchAllKey(params entity.GetVideosParams) string {
	if params.SortOrder == "" {
		params.SortOrder = "-"
	}
	if params.Title == "" {
		params.Title = "-"
	}

	key := fmt.Sprintf("%s-%s:%d:%d:%s", videoCacheKeyPrefx, params.Title, params.Limit, params.Page, params.SortOrder)
	for _, genreID := range params.GenreIDs {
		key += fmt.Sprintf(":%d", genreID)
	}
	for _, order := range params.OrderBy {
		key += fmt.Sprintf(":%s", order)
	}
	return key
}
