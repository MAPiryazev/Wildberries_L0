package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type APIConfig struct {
	APIPort int
}

func LoadAPIConfig() *APIConfig {
	if err := godotenv.Load(".env"); err != nil {
		log.Printf(".env в корне не найден, пробуем ../environment/.env: %v", err)
		// пробуем ../environment/.env
		if err2 := godotenv.Load("../environment/.env"); err2 != nil {
			log.Fatalln("../environment/.env тоже не найден: ", err2)
		}
	}

	APIPort, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		log.Println("ошибка считывания порта API", err)
		APIPort = 8081
	}

	return &APIConfig{APIPort: APIPort}
}
