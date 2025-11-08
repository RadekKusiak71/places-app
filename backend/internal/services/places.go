package services

import (
	"context"
	"errors"
	errors2 "github.com/RadekKusiak71/places-app/internal/errors"
	"github.com/RadekKusiak71/places-app/internal/models"
	"github.com/RadekKusiak71/places-app/internal/stores"
)

type PlacesService struct {
	placesStore *stores.PlacesStore
}

func NewPlacesService(placesStore *stores.PlacesStore) *PlacesService {
	return &PlacesService{
		placesStore: placesStore,
	}
}

func (s *PlacesService) ListPlacesForUser(ctx context.Context, userID int) ([]models.ListPlaceResponse, error) {
	internalPlaces, err := s.placesStore.ListPlacesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]models.ListPlaceResponse, 0, len(internalPlaces))

	for _, place := range internalPlaces {
		response = append(response, models.ListPlaceResponse{
			ID:   place.ID,
			Name: place.Name,
			Lat:  place.Lat,
			Lon:  place.Lon,
		})
	}

	return response, nil
}

func (s *PlacesService) CreatePlaceForUser(ctx context.Context, userID int, req *models.CreatePlaceRequest) (*models.Place, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	place := &models.Place{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Lat:         *req.Lat,
		Lon:         *req.Lon,
	}

	if err := s.placesStore.Create(ctx, place); err != nil {
		return nil, err
	}

	return place, nil
}

func (s *PlacesService) GetPlaceByIDForUser(ctx context.Context, placeID, userID int) (*models.Place, error) {
	place, err := s.placesStore.GetByIDAndUserID(ctx, placeID, userID)
	if err != nil {
		if errors.Is(err, stores.ErrPlaceNotFound) {
			return nil, errors2.PlaceNotFoundError()
		}
		return nil, err
	}
	return place, nil
}

func (s *PlacesService) DeletePlaceByIDForUser(ctx context.Context, placeID, userID int) error {
	err := s.placesStore.DeleteByIDAndUserID(ctx, placeID, userID)
	if err != nil {
		if errors.Is(err, stores.ErrPlaceNotFound) {
			return errors2.PlaceNotFoundError()
		}
		return err
	}
	return nil
}
