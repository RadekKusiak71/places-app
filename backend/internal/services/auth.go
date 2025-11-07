package services

import (
	"context"
	"errors"
	"time"

	errs "github.com/RadekKusiak71/places-app/internal/errors"
	jwt2 "github.com/RadekKusiak71/places-app/internal/jwt"
	"github.com/RadekKusiak71/places-app/internal/models"
	"github.com/RadekKusiak71/places-app/internal/password"
	"github.com/RadekKusiak71/places-app/internal/stores"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userStore *stores.UserStore
	rtStore   *stores.RefreshTokenStore
}

func NewAuthService(userStore *stores.UserStore, rtStore *stores.RefreshTokenStore) *AuthService {
	return &AuthService{
		userStore: userStore,
		rtStore:   rtStore,
	}
}

func (s *AuthService) RefreshTokensPair(ctx context.Context, refreshData *models.TokenPairRefreshRequest) (*models.Tokens, error) {
	if err := refreshData.Validate(); err != nil {
		return nil, err
	}

	refreshTokenClaims, err := jwt2.ParseAndVerifyRefreshToken(refreshData.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			if err := s.rtStore.Delete(ctx, refreshTokenClaims.JTI); err != nil {
				return nil, err
			}
		}
		return nil, errs.InvalidTokenErrorWithMessage(err)
	}

	previousRefreshToken, err := s.rtStore.Get(ctx, refreshTokenClaims.JTI)
	if err != nil {
		if errors.Is(err, stores.ErrRefreshTokenNotFound) {
			return nil, errs.InvalidTokenError()
		}
		return nil, err
	}

	if previousRefreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errs.InvalidTokenError()
	}

	newDBRefreshToken := &models.RefreshToken{
		UserID:    previousRefreshToken.UserID,
		ExpiresAt: jwt2.GetRefreshEXPTime(),
	}

	if err := s.rtStore.Rotate(ctx, previousRefreshToken.ID, newDBRefreshToken); err != nil {
		return nil, err
	}

	return s.generateNewTokenPair(previousRefreshToken.UserID, newDBRefreshToken.ID, newDBRefreshToken.ExpiresAt)
}

func (s *AuthService) ObtainTokensPair(ctx context.Context, loginData *models.LoginRequest) (*models.Tokens, error) {
	if err := loginData.Validate(); err != nil {
		return nil, err
	}

	user, err := s.userStore.GetByUsername(ctx, loginData.Username)
	if err != nil {
		if errors.Is(err, stores.ErrUserNotFound) {
			return nil, errs.ErrInvalidCredentials()
		}
		return nil, err
	}

	if err := password.Compare(loginData.Password, user.Password); err != nil {
		return nil, errs.ErrInvalidCredentials()
	}

	dbRefreshToken := &models.RefreshToken{
		UserID:    user.ID,
		ExpiresAt: jwt2.GetRefreshEXPTime(),
	}

	if err := s.rtStore.Create(ctx, dbRefreshToken); err != nil {
		return nil, err
	}

	return s.generateNewTokenPair(user.ID, dbRefreshToken.ID, dbRefreshToken.ExpiresAt)
}

func (s *AuthService) RegisterUser(ctx context.Context, registerData *models.RegisterRequest) (*models.User, error) {
	if err := registerData.Validate(); err != nil {
		return nil, err
	}

	_, err := s.userStore.GetByUsername(ctx, registerData.Username)

	if err == nil {
		return nil, errs.ErrUserAlreadyExists()
	}

	if !errors.Is(err, stores.ErrUserNotFound) {
		return nil, err
	}

	hashedPassword, err := password.Hash(registerData.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: registerData.Username,
		Password: hashedPassword,
	}

	if err := s.userStore.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) generateNewTokenPair(userID int, newRefreshTokenID string, rtExp time.Time) (*models.Tokens, error) {
	accessTokenExpTS := jwt2.GetAccessEXPTime()

	accessToken, err := jwt2.GenerateAccessToken(userID, accessTokenExpTS)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt2.GenerateRefreshToken(rtExp, newRefreshTokenID)
	if err != nil {
		return nil, err
	}

	return &models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
