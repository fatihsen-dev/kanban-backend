package helpers

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHash(str string) (string, error) {

	if str == "" {
		return "", errors.New("password cannot be empty")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func ValidateHash(hashedStr string, str string) error {
	if hashedStr == "" || str == "" {
		return errors.New("password cannot be empty")
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(str))
}
