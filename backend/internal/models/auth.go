package models

import (
	"strings"

	"github.com/RadekKusiak71/places-app/internal/errors"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenPairRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (r *TokenPairRefreshRequest) Validate() error {
	validationErrors := errors.NewValidationError()

	r.RefreshToken = strings.TrimSpace(r.RefreshToken)

	if r.RefreshToken == "" {
		validationErrors.Add("refresh_token", []string{"refresh_token is required"})
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	validationErrors := errors.NewValidationError()

	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)

	if r.Username == "" {
		validationErrors.Add("username", []string{"username is required"})
	}

	if r.Password == "" {
		validationErrors.Add("password", []string{"password is required"})
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Validate() error {
	validationErrors := errors.NewValidationError()

	r.Username = strings.TrimSpace(r.Username)
	r.Password = strings.TrimSpace(r.Password)

	if r.Username == "" {
		validationErrors.Add("username", []string{"username is required"})
	} else {
		if len(r.Username) < 6 || len(r.Username) > 60 {
			validationErrors.Add("username", []string{"username must be between 6 and 60 characters"})
		}
	}

	if r.Password == "" {
		validationErrors.Add("password", []string{"password is required"})
	} else {
		if len(r.Password) < 8 || len(r.Password) > 120 {
			validationErrors.Add("password", []string{"password must be between 8 and 120 characters"})
		}
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}
