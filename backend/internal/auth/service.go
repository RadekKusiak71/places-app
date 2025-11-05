package auth

import (
	"errors"

	"github.com/RadekKusiak71/places-app/internal/users"
	"github.com/RadekKusiak71/places-app/internal/utils"
)

type AuthService interface {
	RegisterUser(registerData *RegisterPayload) (*users.User, error)
	ObtainTokens(loginData *LoginPayload) (*TokenResponse, error)
}

type Service struct {
	userStore users.UserStore
	authStore AuthStore
}

func NewService(userStore users.UserStore, authStore AuthStore) AuthService {
	return &Service{userStore: userStore, authStore: authStore}
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

	refreshToken := &RefreshToken{
		UserID:    user.ID,
		ExpiresAt: CalculateRefreshTokenExp(),
	}

	// If something after execution of this function fails it will lead to leaving orphaned token in database,
	// I don't think It's critical, but could be solved with executing whole logic inside Transaction
	// We cannot move right before return because we generate UUID on db side
	if err := s.authStore.CreateRefreshToken(refreshToken); err != nil {
		return nil, err
	}

	refreshJWT, err := GenerateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessJWT, err := GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{RefreshToken: refreshJWT, AccessToken: accessJWT}, nil
}

func (s *Service) RegisterUser(registerData *RegisterPayload) (*users.User, error) {
	if err := registerData.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.userStore.GetUser(registerData.Username); err == nil {
		return nil, ErrUsernameIsTaken()
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
