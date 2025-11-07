package auth

import (
	"strings"

	"github.com/RadekKusiak71/places-app/internal/errors"
)

func (p *RefreshTokenPayload) Validate() *errors.ValidationError {
	validationErrors := errors.NewValidationError()

	p.RefreshToken = strings.TrimSpace(p.RefreshToken)

	if p.RefreshToken == "" {
		validationErrors.Add("refresh_token", []string{"refresh_token is missing"})
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

func (p *RegisterPayload) Validate() *errors.ValidationError {
	validationErrors := errors.NewValidationError()

	p.Username = strings.TrimSpace(p.Username)
	p.Password = strings.TrimSpace(p.Password)

	usernameLength := len(p.Username)
	passwordLength := len(p.Password)

	if usernameLength < 8 || usernameLength > 60 {
		validationErrors.Add("username", []string{"username has to be between 8 and 60 length"})
	}

	if passwordLength < 8 || passwordLength > 120 {
		validationErrors.Add("password", []string{"password has to be between 8 and 120 length"})
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}

func (p *LoginPayload) Validate() *errors.ValidationError {
	validationErrors := errors.NewValidationError()

	p.Username = strings.TrimSpace(p.Username)
	p.Password = strings.TrimSpace(p.Password)

	if p.Username == "" {
		validationErrors.Add("username", []string{"username is missing"})
	}

	if p.Password == "" {
		validationErrors.Add("password", []string{"password is missing"})
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}
