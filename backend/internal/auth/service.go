package auth

import (
	"context"
	"errors"
	"time"

	"github.com/RadekKusiak71/places-app/internal/users"
)

type AuthService interface {
	RegisterUser(ctx context.Context, registerData *RegisterPayload) (*users.User, error)
	ObtainJWT(ctx context.Context, loginData *LoginPayload) (*TokenResponse, error)
	RotateRefreshToken(ctx context.Context, refreshTokenPayload *RefreshTokenPayload) (*TokenResponse, error)
}

type Service struct {
	userStore  users.UserStore
	jwtService JWTManager
	authStore  AuthStore
}

func NewService(userStore users.UserStore, jwtService JWTManager, authStore AuthStore) AuthService {
	return &Service{userStore: userStore, jwtService: jwtService, authStore: authStore}
}

func (s *Service) RotateRefreshToken(ctx context.Context, refreshTokenPayload *RefreshTokenPayload) (*TokenResponse, error) {
	if err := refreshTokenPayload.Validate(); err != nil {
		return nil, err
	}

	refreshClaims, err := s.jwtService.ValidateRefreshToken(refreshTokenPayload.RefreshToken)
	if err != nil {
		return nil, InvalidRefreshToken()
	}

	storedToken, err := s.authStore.GetRefreshToken(ctx, refreshClaims.ID)
	if err != nil {
		if errors.Is(err, ErrRefreshTokenNotFound) {
			return nil, InvalidRefreshToken()
		}
		return nil, err
	}

	if storedToken.ExpiresAt.Before(time.Now()) {
		return nil, InvalidRefreshToken()
	}

	newRefreshToken := &RefreshToken{
		UserID:    storedToken.UserID,
		ExpiresAt: s.jwtService.GetRefreshTokenExpiry(),
	}

	if err := s.authStore.RotateRefreshToken(ctx, storedToken.ID, newRefreshToken); err != nil {
		return nil, err
	}

	return s.generateTokenPair(storedToken.UserID, newRefreshToken)
}

func (s *Service) ObtainJWT(ctx context.Context, loginData *LoginPayload) (*TokenResponse, error) {
	if err := loginData.Validate(); err != nil {
		return nil, err
	}

	user, err := s.userStore.GetUser(ctx, loginData.Username)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return nil, InvalidCredentials()
		}
		return nil, err
	}

	if !CheckPassword(loginData.Password, user.Password) {
		return nil, InvalidCredentials()
	}

	refreshToken := &RefreshToken{
		UserID:    user.ID,
		ExpiresAt: s.jwtService.GetRefreshTokenExpiry(),
	}

	if err := s.authStore.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	return s.generateTokenPair(user.ID, refreshToken)
}

func (s *Service) RegisterUser(ctx context.Context, registerData *RegisterPayload) (*users.User, error) {
	if err := registerData.Validate(); err != nil {
		return nil, err
	}

	_, err := s.userStore.GetUser(ctx, registerData.Username)
	if err == nil {
		return nil, ErrUsernameIsTaken()
	}

	if !errors.Is(err, users.ErrUserNotFound) {
		return nil, err
	}

	hashedPassword, err := HashPassword(registerData.Password)
	if err != nil {
		return nil, err
	}

	user := &users.User{
		Username: registerData.Username,
		Password: hashedPassword,
	}

	if err := s.userStore.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) generateTokenPair(userID int, refreshToken *RefreshToken) (*TokenResponse, error) {
	refreshTokenString, err := s.jwtService.GenerateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessTokenString, err := s.jwtService.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
