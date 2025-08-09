package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBMaxConnLifeTime int
}

func LoadDBConfig() *DBConfig {

	if err := godotenv.Load(".env"); err != nil {
		log.Printf(".env в корне не найден, пробуем ../environment/.env: %v", err)
		// пробуем ../environment/.env
		if err2 := godotenv.Load("../environment/.env"); err2 != nil {
			log.Fatalln("../environment/.env тоже не найден: ", err2)
		}
	}

	MaxOpenConns, error := strconv.Atoi(os.Getenv("POSTGRES_MAX_OPEN_CONNS"))
	if error != nil {
		log.Println("ошибка при считывании POSTGRES_MAX_OPEN_CONNS из env", error)
		MaxOpenConns = 10
	}
	DBMaxIdleConns, error := strconv.Atoi(os.Getenv("POSTGRES_MAX_IDLE_CONNS"))
	if error != nil {
		log.Println("ошибка при считывании POSTGRES_MAX_IDLE_CONNS из env", error)
		DBMaxIdleConns = 10
	}
	DBMaxConnLifeTime, error := strconv.Atoi(os.Getenv("POSTGRES_CONN_MAX_LIFETIME"))
	if error != nil {
		log.Println("ошибка при считывании POSTGRES_CONN_MAX_LIFETIME из env", error)
		DBMaxConnLifeTime = 10
	}

	return &DBConfig{
		DBHost:            os.Getenv("POSTGRES_HOST"),
		DBPort:            os.Getenv("POSTGRES_PORT"),
		DBUser:            os.Getenv("POSTGRES_USER"),
		DBPassword:        os.Getenv("POSTGRES_PASSWORD"),
		DBName:            os.Getenv("POSTGRES_NAME"),
		DBSSLMode:         os.Getenv("POSTGRES_SSLMODE"),
		DBMaxOpenConns:    MaxOpenConns,
		DBMaxIdleConns:    DBMaxIdleConns,
		DBMaxConnLifeTime: DBMaxConnLifeTime,
	}
}
