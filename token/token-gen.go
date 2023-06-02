package token

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTClaim struct {
	FirstName string
	LastName  string
	Email     string
	jwt.StandardClaims
}

var secret_key = os.Getenv("jwt_key")

func GenerateToken(firstname string, lastname string, email string) (string, string, int64, error, JWTClaim) {

	expires_at := time.Now().Add(1 * time.Hour).Unix()
	claim := JWTClaim{
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires_at,
		},
	}
	refreshClaim := JWTClaim{
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(secret_key))
	if err != nil {
		return "", "", 0, err, JWTClaim{}
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim).SignedString([]byte(secret_key))
	if err != nil {
		return "", "", 0, err, JWTClaim{}
	}

	return token, refreshToken, expires_at, nil, claim

}
