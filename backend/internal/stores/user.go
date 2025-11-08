package stores

import (
	"context"
	"database/sql"
	"errors"

	"github.com/RadekKusiak71/places-app/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Get(ctx context.Context, userID int) (*models.User, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, username, password, created_at FROM users WHERE id = $1`, userID)

	user := new(models.User)
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserStore) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, username, password, created_at FROM users WHERE username = $1`, username)

	user := new(models.User)
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Create(ctx context.Context, user *models.User) error {
	row := s.db.QueryRowContext(ctx,
		`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, created_at`,
		user.Username, user.Password,
	)
	return row.Scan(&user.ID, &user.CreatedAt)
}
