package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RadekKusiak71/places-app/internal/models"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type RefreshTokenStore struct {
	db *sql.DB
}

func NewRefreshTokenStore(db *sql.DB) *RefreshTokenStore {
	return &RefreshTokenStore{
		db: db,
	}
}

func (s *RefreshTokenStore) Rotate(ctx context.Context, oldID string, rt *models.RefreshToken) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
       DELETE FROM refresh_tokens
       WHERE id = $1
    `, oldID)

	if err != nil {
		return err
	}

	row := tx.QueryRowContext(ctx, `
       INSERT INTO refresh_tokens (user_id, expires_at)
       VALUES ($1, $2)
       RETURNING id, created_at
    `, rt.UserID, rt.ExpiresAt)

	if err := row.Scan(&rt.ID, &rt.CreatedAt); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *RefreshTokenStore) Get(ctx context.Context, id string) (*models.RefreshToken, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, expires_at, created_at
		FROM refresh_tokens
		WHERE id = $1
	`, id)

	rt := new(models.RefreshToken)
	if err := row.Scan(&rt.ID, &rt.UserID, &rt.ExpiresAt, &rt.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return rt, nil
}

func (s *RefreshTokenStore) Delete(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM refresh_tokens
		WHERE id = $1
	`, id)

	return err
}

func (s *RefreshTokenStore) Create(ctx context.Context, rt *models.RefreshToken) error {
	row := s.db.QueryRowContext(ctx, `
		INSERT INTO refresh_tokens (user_id, expires_at)
		VALUES ($1, $2)
		RETURNING id, created_at
	`, rt.UserID, rt.ExpiresAt)

	return row.Scan(&rt.ID, &rt.CreatedAt)
}
