package handlers

import (
	"github.com/RadekKusiak71/places-app/internal/errors"
	"github.com/RadekKusiak71/places-app/internal/middlewares"
	"github.com/RadekKusiak71/places-app/internal/models"
	"github.com/RadekKusiak71/places-app/internal/services"
	"github.com/RadekKusiak71/places-app/internal/utils"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

type PlacesHandler struct {
	placesService *services.PlacesService
}

func NewPlacesHandler(placesService *services.PlacesService) *PlacesHandler {
	return &PlacesHandler{
		placesService: placesService,
	}
}

func (h *PlacesHandler) retrieveUserFromContext(r *http.Request) (*models.User, error) {
	user, ok := r.Context().Value(middlewares.UserContextKey).(*models.User)
	log.Println(user)
	if !ok {
		return &models.User{}, errors.InternalServerError()
	}
	return user, nil
}

func (h *PlacesHandler) retrievePlaceIDFromRequest(r *http.Request) (int, error) {
	placeIDStr := chi.URLParam(r, "placeID")
	placeID, err := strconv.Atoi(placeIDStr)
	if err != nil {
		return 0, err
	}
	return placeID, nil
}

func (h *PlacesHandler) ListPlaces(w http.ResponseWriter, r *http.Request) error {
	user, err := h.retrieveUserFromContext(r)
	if err != nil {
		return err
	}

	places, err := h.placesService.ListPlacesForUser(r.Context(), user.ID)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, places)
}

func (h *PlacesHandler) CreatePlace(w http.ResponseWriter, r *http.Request) error {
	user, err := h.retrieveUserFromContext(r)
	if err != nil {
		return err
	}

	var createPlaceRequest models.CreatePlaceRequest
	if err := utils.ReadJSON(r, &createPlaceRequest); err != nil {
		return errors.InvalidRequestError()
	}

	place, err := h.placesService.CreatePlaceForUser(r.Context(), user.ID, &createPlaceRequest)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusCreated, place)
}

func (h *PlacesHandler) RetrievePlace(w http.ResponseWriter, r *http.Request) error {
	user, err := h.retrieveUserFromContext(r)
	if err != nil {
		return err
	}

	placeID, err := h.retrievePlaceIDFromRequest(r)
	if err != nil {
		return errors.InvalidRequestError()
	}

	place, err := h.placesService.GetPlaceByIDForUser(r.Context(), placeID, user.ID)
	if err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusOK, place)
}

func (h *PlacesHandler) DeletePlace(w http.ResponseWriter, r *http.Request) error {
	user, err := h.retrieveUserFromContext(r)
	if err != nil {
		return err
	}

	placeID, err := h.retrievePlaceIDFromRequest(r)
	if err != nil {
		return errors.InvalidRequestError()
	}

	if err := h.placesService.DeletePlaceByIDForUser(r.Context(), placeID, user.ID); err != nil {
		return err
	}

	return utils.WriteJSON(w, http.StatusNoContent, nil)
}
