package auth

import (
	"context"
	"database/sql"
	"errors"
)

var ErrRefreshTokenNotFound = errors.New("token not found")

type AuthStore interface {
	GetRefreshToken(ctx context.Context, tokenID string) (*RefreshToken, error)
	CreateRefreshToken(ctx context.Context, rt *RefreshToken) error
	RotateRefreshToken(ctx context.Context, oldRtID string, newRt *RefreshToken) error
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) AuthStore {
	return &Store{db: db}
}

func (s *Store) RotateRefreshToken(ctx context.Context, oldRtID string, newRt *RefreshToken) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM refresh_tokens WHERE id = $1`,
		oldRtID,
	)
	if err != nil {
		return err
	}

	row := tx.QueryRowContext(ctx, `
		INSERT INTO refresh_tokens (user_id, expires_at)
		VALUES ($1, $2) RETURNING id, created_at
	`, newRt.UserID, newRt.ExpiresAt)

	if err := row.Scan(&newRt.ID, &newRt.CreatedAt); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateRefreshToken(ctx context.Context, rt *RefreshToken) error {
	row := s.db.QueryRowContext(ctx, `
		INSERT INTO refresh_tokens (user_id, expires_at)
		VALUES ($1, $2) RETURNING id, created_at 
	`, rt.UserID, rt.ExpiresAt)

	return row.Scan(&rt.ID, &rt.CreatedAt)
}

func (s *Store) GetRefreshToken(ctx context.Context, tokenID string) (*RefreshToken, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, user_id, expires_at, created_at
		FROM refresh_tokens
		WHERE id = $1
	`, tokenID)

	var rt = new(RefreshToken)
	if err := row.Scan(&rt.ID, &rt.UserID, &rt.ExpiresAt, &rt.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return rt, nil
}
