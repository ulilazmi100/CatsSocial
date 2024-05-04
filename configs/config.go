package configs

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DbName     string
	DbPort     string
	DbHost     string
	DbUsername string
	DbPassword string

	APPPort string
	ENV     string

	JWTSecret  string
	BcryptSalt int
}

func LoadConfig() (Config, error) {
	config := Config{
		DbName:     os.Getenv("DB_NAME"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
		DbUsername: os.Getenv("DB_USERNAME"),
		DbPassword: os.Getenv("DB_PASSWORD"),

		APPPort: os.Getenv("APP_PORT"),
		ENV:     os.Getenv("ENV"),

		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	salt, err := strconv.Atoi(os.Getenv("BCRYPT_SALT"))
	if err != nil {
		return Config{}, fmt.Errorf("failed to get bcrypt salt %v", err)
	}

	if os.Getenv("APP_PORT") == "" {
		config.APPPort = "8000"
	}

	config.BcryptSalt = salt

	return config, nil
}
