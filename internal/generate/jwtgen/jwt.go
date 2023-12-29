package jwtgen

import (
	"github.com/MorZLE/auth/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewJWT(user models.User, app models.App, timeS time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["uid"] = user.ID
	claims["login"] = user.Login
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(timeS).Unix()

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
