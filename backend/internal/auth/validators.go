package auth

import (
	"strings"

	"github.com/RadekKusiak71/places-app/internal/errors"
)

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

	if len(p.Username) == 0 {
		validationErrors.Add("username", []string{"username is required"})
	}

	if len(p.Password) == 0 {
		validationErrors.Add("password", []string{"password is required"})
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}
