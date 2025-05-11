package config

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

type TokenManager struct {
	Config *viper.Viper
}

type JWTCustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewTokenManager(config *viper.Viper) *TokenManager {
	return &TokenManager{
		Config: config,
	}
}

func (tm *TokenManager) Generate(userID string, role string) (string, error) {
	now := time.Now()

	registeredClaims := JWTCustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:  tm.Config.GetString("APP_NAME"),
			Subject: userID,
			IssuedAt: &jwt.NumericDate{
				Time: now,
			},
			ExpiresAt: &jwt.NumericDate{
				Time: now.Add(24 * time.Hour),
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)

	secret := tm.Config.GetString("JWT_SECRET")
	// sign the token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (tm *TokenManager) Parse(tokenString string) (*JWTCustomClaims, error) {
	var customClaims JWTCustomClaims
	token, err := jwt.ParseWithClaims(
		tokenString,
		&customClaims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			secret := tm.Config.GetString("JWT_SECRET")
			return []byte(secret), nil
		},
		// parser options
		jwt.WithIssuer(tm.Config.GetString("APP_NAME")),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) { // error handling
			return nil, err
		}
		return nil, err
	}

	claims, ok := token.Claims.(*JWTCustomClaims)
	if !ok {
		return nil, errors.New("unknown claims type")
	}

	userID := claims.Subject
	if userID == "" {
		return nil, errors.New("user id is not found on claims")
	}

	userRole := claims.Role
	if userRole == "" {
		return nil, errors.New("user role is not found on claims")
	}

	return claims, nil
}
