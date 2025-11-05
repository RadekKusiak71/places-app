package auth

import (
	"database/sql"
)

type AuthStore interface {
	CreateRefreshToken(token *RefreshToken) error
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) AuthStore {
	return &Store{db: db}
}

func (s *Store) CreateRefreshToken(token *RefreshToken) error {
	row := s.db.QueryRow(
		`INSERT INTO refresh_tokens (user_id,  expires_at) VALUES ($1, $2) RETURNING id, created_at`,
		token.UserID,
		token.ExpiresAt,
	)

	if err := row.Scan(&token.ID, &token.CreatedAt); err != nil {
		return err
	}

	return nil
}
