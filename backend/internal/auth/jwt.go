package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/RadekKusiak71/places-app/config"
	"github.com/RadekKusiak71/places-app/internal/users"
	"github.com/golang-jwt/jwt/v5"
)

func CalculateRefreshTokenExp() time.Time {
	return time.Now().Add(time.Second * time.Duration(config.Config.JWT_REFRESH_EXP_SECONDS))
}

func GenerateRefreshToken(refreshToken *RefreshToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshToken.ExpiresAt),
		Subject:   strconv.Itoa(refreshToken.UserID),
		ID:        refreshToken.ID,
	})
	tokenString, err := token.SignedString([]byte(config.Config.JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateAccessToken(user *users.User) (string, error) {
	accessExp := time.Second * time.Duration(config.Config.JWT_ACCESS_EXP_SECONDS)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExp)),
		Subject:   strconv.Itoa(user.ID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString([]byte(config.Config.JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string) (jwt.RegisteredClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(config.Config.JWT_SECRET), nil
	})
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}

	claims, ok := token.Claims.(jwt.RegisteredClaims)
	if !ok {
		return jwt.RegisteredClaims{}, errors.New("invalid token")
	}
	return claims, nil
}
