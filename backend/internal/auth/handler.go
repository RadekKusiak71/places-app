package auth

import (
	"net/http"

	"github.com/RadekKusiak71/places-app/internal/utils"
)

type AuthHandler interface {
	Register(w http.ResponseWriter, r *http.Request) error
	ObtainJWT(w http.ResponseWriter, r *http.Request) error
}

type Handler struct {
	authService AuthService
}

func NewAuthHandler(authService AuthService) AuthHandler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) ObtainJWT(w http.ResponseWriter, r *http.Request) error {
	var loginPayload LoginPayload
	if err := utils.ReadJSON(r, &loginPayload); err != nil {
		return err
	}

	tokenResponse, err := h.authService.ObtainTokens(&loginPayload)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, tokenResponse)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) error {
	var registerPayload RegisterPayload

	if err := utils.ReadJSON(r, &registerPayload); err != nil {
		return err
	}

	user, err := h.authService.RegisterUser(&registerPayload)

	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, user)
}
