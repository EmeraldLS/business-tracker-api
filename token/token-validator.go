package token

import (
	"errors"
	"time"

	"github.com/EmeraldLS/phsps-api/config"
	"github.com/dgrijalva/jwt-go"
)

func ValidateToken(signedToken string) error {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret_key), nil
	})

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

	return nil
}
