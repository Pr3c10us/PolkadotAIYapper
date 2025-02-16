package utils

import (
	"errors"
	"fmt"
	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
	Expires int64  `json:"expires"`
	jwt.RegisteredClaims
}

func CreateUserToken(user *authentication.User, secret string, duration time.Duration) (string, error) {
	now := time.Now()
	expires := now.Add(duration).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", appError.InternalServerError(errors.New("failed to authenticate user"))
	}
	return tokenStr, nil
}

func CreateUserRefreshToken(user *authentication.User, secret string, duration time.Duration) (string, error) {
	now := time.Now()
	expires := now.Add(duration).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"version": user.RefreshTokenVersion,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", appError.InternalServerError(errors.New("failed to authenticate user"))
	}
	return tokenStr, nil
}

func DecryptUserToken(tokenStr string, secret string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, appError.Unauthorized(fmt.Errorf("internal server error"))
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate the token and extract the claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if the token is expired
		if claims.Expires < time.Now().Unix() {
			return nil, appError.Unauthorized(errors.New("token is expired"))
		}
		return claims, nil
	}

	return nil, appError.Unauthorized(errors.New("invalid token"))
}
