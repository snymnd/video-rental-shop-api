package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"vrs-api/internal/constant"
	"vrs-api/internal/customerrors"
	util "vrs-api/internal/util/jwt"

	"github.com/gin-gonic/gin"
)

type (
	RBACRepository interface {
		CheckRoleAccess(ctx context.Context, role, permission, resource int) (bool, error)
	}
	RBACCacheRepository interface {
		CheckRoleAccess(ctx context.Context, role, permission, resource int) (*bool, error)
		SetCheckRoleAccess(ctx context.Context, role, permission, resource, duration int, value bool) error
	}
)

func AuthorizationMiddleware(permission, resource int, rbacr RBACRepository, rbacCache RBACCacheRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token payload and user id from token subject claims
		authPayloadCtx := ctx.Value(constant.CTX_AUTH_PAYLOAD_KEY)
		if authPayloadCtx == nil {
			ctx.Error(
				customerrors.NewError(
					"cannot identified user",
					errors.New("cannot get auth payload from auth payload context"),
					customerrors.Unauthenticate,
				))
			ctx.Abort()
			return
		}
		authPayload := authPayloadCtx.(*util.JWTCustomClaims)
		role := authPayload.Role

		cacheHasAccess, err := rbacCache.CheckRoleAccess(ctx, role, permission, resource)
		if err != nil {
			log.Printf("failed to get data from cache, %s", err.Error())
		}
		unauthorizedErr := customerrors.NewError(
			"unauthorize",
			fmt.Errorf("rbac for role: %d, permission: %d, resource: %d is not found in rbac records", role, permission, resource),
			customerrors.Unauthorized,
		)
		if cacheHasAccess != nil {
			log.Printf("use rbac data from redis cache, rbac-%d:%d:%d => %v", role, permission, resource, *cacheHasAccess)
			if !*cacheHasAccess {
				ctx.Error(unauthorizedErr)
				ctx.Abort()
				return
			}

			ctx.Next()
			return
		}

		hasAccess, checkRoleErr := rbacr.CheckRoleAccess(ctx, role, permission, resource)
		if checkRoleErr != nil {
			ctx.Error(customerrors.NewError(
				"cannot identified access",
				errors.New("cannot check if role has access"),
				customerrors.Unauthorized,
			))
			ctx.Abort()
			return
		}

		const cacheDuration = 604800 // 1 week
		if !hasAccess {
			ctx.Error(unauthorizedErr)
			if err := rbacCache.SetCheckRoleAccess(ctx, role, permission, resource, cacheDuration, false); err != nil {
				log.Printf("failed to set rbac key rbac-%d:%d:%d, error: %s\n", role, permission, resource, err.Error())
			}
			ctx.Abort()
			return
		}

		if err := rbacCache.SetCheckRoleAccess(ctx, role, permission, resource, cacheDuration, true); err != nil {
			log.Printf("failed to set rbac key rbac-%d:%d:%d, error: %s\n", role, permission, resource, err.Error())
		}

		ctx.Next()
	}
}
