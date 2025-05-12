package middleware

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"vrs-api/internal/constant"
	"vrs-api/internal/customerrors"
	util "vrs-api/internal/util/jwt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

type Token interface {
	Parse(tokenString string) (*util.JWTCustomClaims, error)
}

func AuthenticateMiddleware(token Token) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			log.Println(err)

			ctx.Error(customerrors.NewError(
				"authorization header not found",
				err,
				customerrors.Unauthenticate,
			))
			ctx.Abort()
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := errors.New("invalid token format")
			log.Println(err)

			ctx.Error(customerrors.NewError(
				"invalid authorization header format",
				err,
				customerrors.Unauthenticate,
			))
			ctx.Abort()
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			log.Println(err)

			ctx.Error(customerrors.NewError(
				"unsuppoerted authorization type",
				err,
				customerrors.Unauthenticate,
			))
			ctx.Abort()
			return
		}

		accessToken := fields[1]
		payload, err := token.Parse(accessToken)
		if err != nil {
			log.Println(err)
			errorMessage := "cannot parse token"
			if errors.Is(err, jwt.ErrTokenExpired) {
				errorMessage = "token has expired"
			}
			ctx.Error(customerrors.NewError(
				errorMessage,
				err,
				customerrors.Unauthenticate,
			))
			ctx.Abort()
			return
		}

		ctx.Set(constant.CTX_AUTH_PAYLOAD_KEY, payload)
		ctx.Next()
	}
}
