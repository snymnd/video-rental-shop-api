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
	keyCache := rbaccr.getCacheRepositoryKey(fmt.Sprintf("%d:%d:%d", role, permission, resource))
	rbac := rbaccr.redisClient.Get(ctx, keyCache)
	if rbac == nil {
		err := fmt.Errorf("cache key %s is not found", keyCache)
		return nil, err
	}
	hasAccess, err := rbac.Bool()
	if err != nil {
		fmt.Println("error =>", err)
		return nil, err
	}

	return &hasAccess, nil
}

func (rbaccr *RBACCacheRepository) SetCheckRoleAccess(ctx context.Context, role, permission, resource, duration int, value bool) error {
	keyCache := rbaccr.getCacheRepositoryKey(fmt.Sprintf("%d:%d:%d", role, permission, resource))
	if err := rbaccr.redisClient.Set(ctx, keyCache, value, time.Second*time.Duration(duration)).Err(); err != nil {
		return err
	}

	return nil
}

const basePrefix = "rbac"

func (rbaccr *RBACCacheRepository) getCacheRepositoryKey(key string) string {
	return fmt.Sprintf("%s-%s", basePrefix, key)
}
