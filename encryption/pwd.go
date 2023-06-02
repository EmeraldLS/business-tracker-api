package encryption

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	passByte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(passByte), nil
}

func ValidatePassword(providedPwd string, hashedPwd string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(providedPwd)); err != nil {
		return errors.New("incorrect password")
	}
	return nil
}
