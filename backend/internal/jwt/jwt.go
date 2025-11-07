package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/RadekKusiak71/places-app/config"
	"github.com/RadekKusiak71/places-app/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GetRefreshEXPTime() time.Time {
	return time.Now().Add(time.Duration(config.Config.JWT_REFRESH_EXP_SECONDS) * time.Second)
}
func GetAccessEXPTime() time.Time {
	return time.Now().Add(time.Duration(config.Config.JWT_ACCESS_EXP_SECONDS) * time.Second)
}

func signToken(token *jwt.Token) (string, error) {
	return token.SignedString([]byte(config.Config.JWT_SECRET))
}

func GenerateAccessToken(userID int, expTime time.Time) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, models.AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	})

	return signToken(claims)
}

func GenerateRefreshToken(expTime time.Time, jti string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, models.RefreshTokenClaims{
		JTI: jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	})
	return signToken(claims)
}

func ParseAndVerifyAccessToken(tokenStr string) (*models.AccessTokenClaims, error) {
	claims := &models.AccessTokenClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, tokenParseKeyFunc)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, err
		}
		return nil, err
	}

	return claims, nil
}

func ParseAndVerifyRefreshToken(tokenStr string) (*models.RefreshTokenClaims, error) {
	claims := &models.RefreshTokenClaims{}

	_, err := jwt.ParseWithClaims(tokenStr, claims, tokenParseKeyFunc)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, err
		}
		return nil, err
	}

	return claims, nil
}

func tokenParseKeyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(config.Config.JWT_SECRET), nil
}
