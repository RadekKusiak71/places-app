package auth

import (
	"github.com/RadekKusiak71/places-app/internal/users"
	"github.com/RadekKusiak71/places-app/internal/utils"
)

type AuthService interface {
	RegisterUser(registerData *RegisterPayload) (*users.User, error)
}

type Service struct {
	userStore users.UserStore
}

func NewService(userStore users.UserStore) AuthService {
	return &Service{userStore: userStore}
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
