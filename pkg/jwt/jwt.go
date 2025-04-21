package jwt

import (
	"time"

	"github.com/fatihsen-dev/kanban-backend/config"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string, name string, email string, isAdmin bool) (string, error) {
	appConfig := config.Read()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		ID:      userID,
		Name:    name,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	})

	return token.SignedString([]byte(appConfig.JWTSecret))
}

func VerifyToken(tokenString string) (*UserClaims, error) {
	appConfig := config.Read()

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(appConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(*UserClaims), nil
}
