package config

import (
	"os"
	"strconv"
)

type EnvVariables struct {
	PORT                    string
	DB_HOST                 string
	DB_PORT                 string
	DB_USER                 string
	DB_PASSWORD             string
	DB_NAME                 string
	JWT_SECRET              string
	JWT_ACCESS_EXP_SECONDS  int
	JWT_REFRESH_EXP_SECONDS int
}

var Config EnvVariables

func InitConfig() {
	Config = loadEnvVariables()
}

func loadEnvVariables() EnvVariables {
	return EnvVariables{
		PORT:                    getVar("GO_PORT", "8080"),
		DB_HOST:                 getVar("DB_HOST", "localhost"),
		DB_PORT:                 getVar("DB_PORT", "5432"),
		DB_USER:                 getVar("DB_USER", "user"),
		DB_PASSWORD:             getVar("DB_PASSWORD", "password"),
		DB_NAME:                 getVar("DB_NAME", "dbname"),
		JWT_SECRET:              getVar("GO_JWT_SECRET_KEY", "secret"),
		JWT_ACCESS_EXP_SECONDS:  getVarAsInt("GO_JWT_ACCESS_EXP_SECONDS", 900),
		JWT_REFRESH_EXP_SECONDS: getVarAsInt("GO_JWT_REFRESH_EXP_SECONDS", 900),
	}
}

func getVar(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getVarAsInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
		return fallback
	}
	return fallback
}
