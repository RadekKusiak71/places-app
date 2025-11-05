package auth

import "github.com/RadekKusiak71/places-app/internal/errors"

func (p *RegisterPayload) Validate() *errors.ValidationError {
	validationErrors := errors.NewValidationError()
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
