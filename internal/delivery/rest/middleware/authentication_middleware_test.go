package middleware

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"
	"vrs-api/internal/customerrors"
	"vrs-api/internal/util/token"
	"vrs-api/internal/util/viper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func generateJWTToken(userID string, role int, issuetAt time.Time, expiredAt time.Time) (string, error) {
	config := viper.NewViper()
	claims := token.JWTCustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  config.GetString("APP_NAME"),
			Subject: userID,
			IssuedAt: &jwt.NumericDate{
				Time: issuetAt,
			},
			ExpiresAt: &jwt.NumericDate{
				Time: expiredAt,
			},
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := config.GetString("JWT_SECRET")
	// sign the token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func addAuthorization(
	t *testing.T,
	ctx *gin.Context,
	authorizationType string,
	userId string,
	role int,
	expiredAt time.Duration,
) {
	token, err := generateJWTToken(userId, role, time.Now(), time.Now().Add(expiredAt))
	assert.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	ctx.Request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestAuthMiddleware(t *testing.T) {
	user := struct {
		userID string
		role   int
	}{
		userID: "someUserID",
		role:   2, // user
	}

	testCases := []struct {
		name           string
		setupAuth      func(t *testing.T, ctx *gin.Context)
		wantErrMessage string
	}{
		{
			name: "should return 200 when valid token provided",
			setupAuth: func(t *testing.T, ctx *gin.Context) {
				addAuthorization(t, ctx, authorizationTypeBearer, user.userID, user.role, time.Minute)
			},
		},
		{
			name: "should return 401 when no token provided",
			setupAuth: func(t *testing.T, ctx *gin.Context) {
			},
			wantErrMessage: "authorization header not found",
		},
		{
			name: "should return 401 when unsupported token type provided",
			setupAuth: func(t *testing.T, ctx *gin.Context) {
				addAuthorization(t, ctx, "unsupported", user.userID, user.role, time.Minute)
			},
			wantErrMessage: "unsupported authorization type",
		},
		{
			name: "should return 401 when token format invalid",
			setupAuth: func(t *testing.T, ctx *gin.Context) {
				addAuthorization(t, ctx, "wrong token bearer format", user.userID, user.role, time.Minute)
			},
			wantErrMessage: "invalid authorization header format",
		},
		{
			name: "should return 401 when token expired",
			setupAuth: func(t *testing.T, ctx *gin.Context) {
				addAuthorization(t, ctx, authorizationTypeBearer, user.userID, user.role, -time.Minute)
			},
			wantErrMessage: "token has expired",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)

			ctx.Request = httptest.NewRequest("GET", "/", nil)
			tt.setupAuth(t, ctx)
			config := viper.NewViper()

			AuthenticateMiddleware(&token.TokenManager{Config: config})(ctx)

			var got string
			if len(ctx.Errors) > 0 {
				err := ctx.Errors[0]
				var ce *customerrors.CustomError
				if errors.As(err, &ce) {
					got = ce.ErrorMessage
				}

			}
			assert.Equal(t, tt.wantErrMessage, got)
		})
	}
}
