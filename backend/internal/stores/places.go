package stores

import (
	"context"
	"database/sql"
	"errors"
	"github.com/RadekKusiak71/places-app/internal/models"
)

var (
	ErrPlaceNotFound = errors.New("place not found")
)

type PlacesStore struct {
	db *sql.DB
}

func NewPlacesStore(db *sql.DB) *PlacesStore {
	return &PlacesStore{
		db: db,
	}
}

func (s *PlacesStore) GetByIDAndUserID(ctx context.Context, placeID, userID int) (*models.Place, error) {
	query := "SELECT id, user_id, name, description, lat, lon, created_at FROM places WHERE id = $1 AND user_id = $2"
	row := s.db.QueryRowContext(ctx, query, placeID, userID)

	place := new(models.Place)
	err := row.Scan(
		&place.ID,
		&place.UserID,
		&place.Name,
		&place.Description,
		&place.Lat,
		&place.Lon,
		&place.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPlaceNotFound
		}
		return nil, err
	}

	return place, nil
}

func (s *PlacesStore) DeleteByIDAndUserID(ctx context.Context, placeID, userID int) error {
	res, err := s.db.ExecContext(
		ctx,
		"DELETE FROM places WHERE id = $1 AND user_id = $2",
		placeID,
		userID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrPlaceNotFound
	}

	return err
}

func (s *PlacesStore) Create(ctx context.Context, place *models.Place) error {
	query := "INSERT INTO places (user_id, name, description, lat, lon) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at"
	row := s.db.QueryRowContext(
		ctx,
		query,
		place.UserID,
		place.Name,
		place.Description,
		place.Lat,
		place.Lon,
	)
	return row.Scan(&place.ID, &place.CreatedAt)
}

func (s *PlacesStore) ListPlacesByUserID(ctx context.Context, userID int) ([]models.Place, error) {
	query := "SELECT id, user_id, name, description, lat, lon ,created_at FROM places WHERE user_id = $1"
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var places []models.Place

	for rows.Next() {
		var place models.Place

		if err := rows.Scan(
			&place.ID,
			&place.UserID,
			&place.Name,
			&place.Description,
			&place.Lat,
			&place.Lon,
			&place.CreatedAt,
		); err != nil {
			return nil, err
		}

		places = append(places, place)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return places, nil
}
