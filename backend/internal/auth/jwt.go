package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/RadekKusiak71/places-app/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager interface {
	GenerateAccessToken(userID int) (string, error)
	GenerateRefreshToken(refreshToken *RefreshToken) (string, error)
	ValidateAccessToken(tokenString string) (*AccessTokenClaim, error)
	ValidateRefreshToken(tokenString string) (*RefreshTokenClaim, error)
	GetRefreshTokenExpiry() time.Time
}

type RefreshTokenClaim struct {
	jwt.RegisteredClaims
}

type AccessTokenClaim struct {
	jwt.RegisteredClaims
}

type jwtService struct {
	secret         []byte
	accessExpSecs  int
	refreshExpSecs int
}

func NewJWTService() JWTManager {
	return &jwtService{
		secret:         []byte(config.Config.JWT_SECRET),
		accessExpSecs:  config.Config.JWT_ACCESS_EXP_SECONDS,
		refreshExpSecs: config.Config.JWT_REFRESH_EXP_SECONDS,
	}
}

func (s *jwtService) GetRefreshTokenExpiry() time.Time {
	return time.Now().Add(time.Second * time.Duration(s.refreshExpSecs))
}

func (s *jwtService) GenerateRefreshToken(refreshToken *RefreshToken) (string, error) {
	claims := RefreshTokenClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshToken.ExpiresAt),
			ID:        refreshToken.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *jwtService) GenerateAccessToken(userID int) (string, error) {
	accessExp := time.Second * time.Duration(s.accessExpSecs)
	claims := AccessTokenClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExp)),
			Subject:   strconv.Itoa(userID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *jwtService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaim, error) {
	claims := &RefreshTokenClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *jwtService) ValidateAccessToken(tokenString string) (*AccessTokenClaim, error) {
	claims := &AccessTokenClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
