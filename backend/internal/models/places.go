package models

import (
	"github.com/RadekKusiak71/places-app/internal/errors"
	"strings"
	"time"
)

type Place struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Lat         float64    `json:"lat"`
	Lon         float64    `json:"lon"`
	CreatedAt   *time.Time `json:"created_at"`
}

type ListPlaceResponse struct {
	ID   int     `json:"id"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

type CreatePlaceRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Lat         *float64 `json:"lat"`
	Lon         *float64 `json:"lon"`
}

func (r *CreatePlaceRequest) Validate() error {
	validationErrors := errors.NewValidationError()

	r.Name = strings.TrimSpace(r.Name)
	r.Description = strings.TrimSpace(r.Description)

	if r.Name == "" {
		validationErrors.Add("name", []string{"name is required"})
	}

	if r.Lat == nil {
		validationErrors.Add("lat", []string{"lat is required"})
	} else {
		if *r.Lat < -90 || *r.Lat > 90 {
			validationErrors.Add("lat", []string{"lat must be between -90 and 90"})
		}
	}

	if r.Lon == nil {
		validationErrors.Add("lon", []string{"lon is required"})
	} else {
		if *r.Lon < -180 || *r.Lon > 180 {
			validationErrors.Add("lon", []string{"lon must be between -180 and 180"})
		}
	}

	if validationErrors.HasErrors() {
		return validationErrors
	}

	return nil
}
