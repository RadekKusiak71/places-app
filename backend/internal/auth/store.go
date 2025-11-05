package auth

import (
	"database/sql"
	"errors"
)

var ErrRefreshTokenNotFound = errors.New("refresh token not found")

type AuthStore interface {
	CreateRefreshToken(token *RefreshToken) error
	GetRefreshTokenByID(ID string) (*RefreshToken, error)
	DeleteRefreshToken(ID string) error
	RotateRefreshToken(oldTokenID string, newRefreshToken *RefreshToken) error
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) AuthStore {
	return &Store{db: db}
}

func (s *Store) RotateRefreshToken(oldTokenID string, newRefreshToken *RefreshToken) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM refresh_tokens WHERE id = $1", oldTokenID)
	if err != nil {
		return err
	}

	row := tx.QueryRow(
		"INSERT INTO refresh_tokens (user_id, expires_at) VALUES ($1, $2) RETURNING id",
		newRefreshToken.UserID,
		newRefreshToken.ExpiresAt,
	)

	if err := row.Scan(&newRefreshToken.ID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteRefreshToken(ID string) error {
	res, err := s.db.Exec("DELETE FROM refresh_tokens WHERE id=$1", ID)

	if err != nil {
		return err
	}

	n, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if n == 0 {
		return ErrRefreshTokenNotFound
	}

	return nil
}

func (s *Store) GetRefreshTokenByID(ID string) (*RefreshToken, error) {
	row := s.db.QueryRow("SELECT * FROM refresh_tokens WHERE id = $1", ID)

	refreshToken := new(RefreshToken)

	if err := row.Scan(&refreshToken.ID, &refreshToken.UserID, &refreshToken.ExpiresAt, &refreshToken.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}

	return refreshToken, nil
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
