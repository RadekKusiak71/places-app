package handlers

import (
	"net/http"

	"github.com/RadekKusiak71/places-app/internal/models"
	"github.com/RadekKusiak71/places-app/internal/services"
	"github.com/RadekKusiak71/places-app/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RefreshTokensPair(w http.ResponseWriter, r *http.Request) error {
	var refreshReq models.TokenPairRefreshRequest
	if err := utils.ReadJSON(r, &refreshReq); err != nil {
		return err
	}

	tokensPair, err := h.authService.RefreshTokensPair(r.Context(), &refreshReq)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, tokensPair)
}

func (h *AuthHandler) ObtainTokensPair(w http.ResponseWriter, r *http.Request) error {
	var loginReq models.LoginRequest
	if err := utils.ReadJSON(r, &loginReq); err != nil {
		return err
	}

	tokensPair, err := h.authService.ObtainTokensPair(r.Context(), &loginReq)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, tokensPair)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var registerReq models.RegisterRequest
	if err := utils.ReadJSON(r, &registerReq); err != nil {
		return err
	}

	createdUser, err := h.authService.RegisterUser(r.Context(), &registerReq)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, createdUser)
}
