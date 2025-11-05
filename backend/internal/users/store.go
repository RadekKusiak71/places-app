package users

import (
	"database/sql"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type UserStore interface {
	CreateUser(user *User) error
	GetUser(username string) (*User, error)
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) UserStore {
	return &Store{db: db}
}

func (s *Store) GetUser(username string) (*User, error) {
	row := s.db.QueryRow(
		`SELECT id, username, password, created_at FROM users WHERE username = $1`,
		username,
	)
	user := new(User)
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *Store) CreateUser(user *User) error {
	row := s.db.QueryRow(
		`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, created_at`,
		user.Username,
		user.Password,
	)
	return row.Scan(&user.ID, &user.CreatedAt)
}
