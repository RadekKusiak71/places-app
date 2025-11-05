package auth

import (
	"errors"
	"time"

	"github.com/RadekKusiak71/places-app/internal/users"
	"github.com/RadekKusiak71/places-app/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	RegisterUser(registerData *RegisterPayload) (*users.User, error)
	ObtainTokens(loginData *LoginPayload) (*TokenResponse, error)
	RefreshTokens(refreshData *RefreshPayload) (*TokenResponse, error)
}

type Service struct {
	userStore users.UserStore
	authStore AuthStore
}

func NewService(userStore users.UserStore, authStore AuthStore) AuthService {
	return &Service{userStore: userStore, authStore: authStore}
}

func (s *Service) RefreshTokens(refreshData *RefreshPayload) (*TokenResponse, error) {
	if err := refreshData.Validate(); err != nil {
		return nil, err
	}
	token, err := ValidateToken(refreshData.RefreshToken)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, InvalidRefreshToken()
	}

	claimJTI, ok := claims["jti"].(string)
	if !ok {
		return nil, InvalidRefreshToken()
	}

	tokenDB, err := s.authStore.GetRefreshTokenByID(claimJTI)

	if err != nil {
		if errors.Is(err, ErrRefreshTokenNotFound) {
			return nil, InvalidRefreshToken()
		}
		return nil, err
	}
	if time.Now().Compare(tokenDB.ExpiresAt) > 0 {
		return nil, InvalidRefreshToken()
	}
	newRefreshToken := &RefreshToken{
		UserID:    tokenDB.UserID,
		ExpiresAt: CalculateRefreshTokenExp(),
	}

	if err := s.authStore.RotateRefreshToken(tokenDB.ID, newRefreshToken); err != nil {
		return nil, err
	}
	refreshJWT, err := GenerateRefreshToken(newRefreshToken)
	if err != nil {
		return nil, err
	}
	accessJWT, err := GenerateAccessToken(tokenDB.UserID)
	if err != nil {
		return nil, err
	}
	return &TokenResponse{RefreshToken: refreshJWT, AccessToken: accessJWT}, nil
}

func (s *Service) ObtainTokens(loginData *LoginPayload) (*TokenResponse, error) {
	if err := loginData.Validate(); err != nil {
		return nil, err
	}

	user, err := s.userStore.GetUser(loginData.Username)
	if err != nil {
		if errors.Is(err, users.ErrUserNotFound) {
			return nil, InvalidCredentials()
		}
	}

	if !utils.CheckPasswordHash(loginData.Password, user.Password) {
		return nil, InvalidCredentials()
	}

	return s.issueTokenPair(user.ID)
}

func (s *Service) issueTokenPair(userID int) (*TokenResponse, error) {
	refreshToken := &RefreshToken{
		UserID:    userID,
		ExpiresAt: CalculateRefreshTokenExp(),
	}

	if err := s.authStore.CreateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	refreshJWT, err := GenerateRefreshToken(refreshToken)
	if err != nil {
		// Best-effort cleanup: If JWT generation fails, delete the
		// orphaned token we just created.
		_ = s.authStore.DeleteRefreshToken(refreshToken.ID)
		return nil, err
	}

	accessJWT, err := GenerateAccessToken(userID)
	if err != nil {
		_ = s.authStore.DeleteRefreshToken(refreshToken.ID)
		return nil, err
	}

	return &TokenResponse{RefreshToken: refreshJWT, AccessToken: accessJWT}, nil
}

func (s *Service) RegisterUser(registerData *RegisterPayload) (*users.User, error) {
	if err := registerData.Validate(); err != nil {
		return nil, err
	}

	_, err := s.userStore.GetUser(registerData.Username)
	if err == nil {
		return nil, ErrUsernameIsTaken()
	}

	if !errors.Is(err, users.ErrUserNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(registerData.Password)
	if err != nil {
		return nil, err
	}

	user := &users.User{
		Username: registerData.Username,
		Password: hashedPassword,
	}

	if err := s.userStore.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}
