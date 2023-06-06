package token

import (
	"errors"
	"time"

	"github.com/EmeraldLS/phsps-api/config"
	"github.com/EmeraldLS/phsps-api/model"
	"github.com/dgrijalva/jwt-go"
)

func ValidateToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
	})

	if err != nil {
		err = errors.New("invalid token")
		return err
	}

	_, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("invalid token")
		return err
	}

	user, err := config.FindUserByToken(signedToken)
	if user.ExpiresAt < time.Now().Unix() {
		err = errors.New("token has expired")
		return err
	}

	if err != nil {
		err = errors.New("invalid token")
		return err
	}

	return nil
}

func ValidateRefreshToken(refresh_token string) (model.User, error) {
	signedRefresh_token, err := jwt.ParseWithClaims(refresh_token, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
	})

	if err != nil {
		err = errors.New("invalid refresh token")
		return model.User{}, err
	}

	_, ok := signedRefresh_token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("invalid refresh token")
		return model.User{}, err
	}

	user, err := config.FindUserByRefreshToken(refresh_token)
	if user.ExpiresAt < time.Now().Unix() {
		err = errors.New("refresh token has expired")
		return model.User{}, err
	}

	if err != nil {
		err = errors.New("invalid refresh token")
		return model.User{}, err
	}

	return user, nil
}
