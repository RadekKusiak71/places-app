package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RadekKusiak71/places-app/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func New() *sql.DB {
	db, err := sql.Open(
		"pgx",
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Config.DB_HOST,
			config.Config.DB_PORT,
			config.Config.DB_USER,
			config.Config.DB_PASSWORD,
			config.Config.DB_NAME,
		),
	)

	if err != nil {
		log.Fatalf("Couldn't open database: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Couldn't establish connection with database: %s", err.Error())
	}

	return db
}
