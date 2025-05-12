package middleware

import (
	"context"
	"errors"
	"vrs-api/internal/constant"
	"vrs-api/internal/customerrors"
	util "vrs-api/internal/util/jwt"

	"github.com/gin-gonic/gin"
)

type RBACRepository interface {
	HasAccess(ctx context.Context, role int, permission int, resource int) (bool, error)
}

func AuthorizationMiddleware(permission int, resource int, rbacr RBACRepository) gin.HandlerFunc {
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

		ok, hasAccessErr := rbacr.HasAccess(ctx, role, permission, resource)
		if hasAccessErr != nil {
			ctx.Error(
				customerrors.NewError(
					"cannot identified access",
					errors.New("cannot check if role has access"),
					customerrors.Unauthorized,
				))
			ctx.Abort()
			return
		}
		if !ok {
			ctx.Error(
				customerrors.NewError(
					"unauthorize",
					errors.New("user's role dont have access to this permission/resource"),
					customerrors.Unauthorized,
				))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
