package config

import (
	"os"
)

type EnvVariables struct {
	PORT        string
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
}

var Config EnvVariables

func InitConfig() {
	Config = loadEnvVariables()
}

func loadEnvVariables() EnvVariables {
	return EnvVariables{
		PORT:        getVar("GO_PORT", "8080"),
		DB_HOST:     getVar("DB_HOST", "localhost"),
		DB_PORT:     getVar("DB_PORT", "5432"),
		DB_USER:     getVar("DB_USER", "user"),
		DB_PASSWORD: getVar("DB_PASSWORD", "password"),
		DB_NAME:     getVar("DB_NAME", "dbname"),
	}
}

func getVar(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
