package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RBACCacheRepository struct {
	redisClient *redis.Client
}

func NewRBACCacheRepository(redisClient *redis.Client) *RBACCacheRepository {
	return &RBACCacheRepository{redisClient}
}

func (rbaccr *RBACCacheRepository) CheckRoleAccess(ctx context.Context, role, permission, resource int) (*bool, error) {
	cacheKey := rbaccr.getRBACCacheRepositoryKey(fmt.Sprintf("%d:%d:%d", role, permission, resource))
	rbac := rbaccr.redisClient.Get(ctx, cacheKey)
	if rbac == nil {
		err := fmt.Errorf("cache key %s is not found", cacheKey)
		return nil, err
	}
	hasAccess, err := rbac.Bool()
	if err != nil {
		return nil, err
	}

	return &hasAccess, nil
}

func (rbaccr *RBACCacheRepository) SetCheckRoleAccess(ctx context.Context, role, permission, resource, duration int, value bool) error {
	cacheKey := rbaccr.getRBACCacheRepositoryKey(fmt.Sprintf("%d:%d:%d", role, permission, resource))
	if err := rbaccr.redisClient.Set(ctx, cacheKey, value, time.Second*time.Duration(duration)).Err(); err != nil {
		return err
	}

	return nil
}

const rbacCacheKeyPrefix = "rbac"

func (rbaccr *RBACCacheRepository) getRBACCacheRepositoryKey(key string) string {
	return fmt.Sprintf("%s-%s", rbacCacheKeyPrefix, key)
}
