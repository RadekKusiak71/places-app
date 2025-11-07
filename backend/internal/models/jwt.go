package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type RefreshTokenClaims struct {
	JTI string `json:"jti"`
	jwt.RegisteredClaims
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
}
