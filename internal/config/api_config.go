package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	ErrEnvNotFound = errors.New("Env файл не был найден")
)

type APIConfig struct {
	APIPort string
}

func LoadAPIConfig() (*APIConfig, error) {
	if err := godotenv.Load(".env"); err != nil {
		if err2 := godotenv.Load("../environment/.env"); err2 != nil {
			return nil, fmt.Errorf("%w: %v", ErrEnvNotFound, err2)
		}
	}

	APIPort := os.Getenv("API_PORT")
	if APIPort == "" {
		log.Println("ошибка считывания порта API, используем порт 8081")
		APIPort = "8081"
	}

	return &APIConfig{APIPort: APIPort}, nil
}
