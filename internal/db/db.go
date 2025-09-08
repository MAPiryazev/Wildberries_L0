package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/MAPiryazev/Wildberries_L0/internal/config"
	_ "github.com/lib/pq"
)

func InitPsqlDB(cfg *config.DBConfig) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	psqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД: ", err)
	}

	psqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	psqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	psqlDB.SetConnMaxLifetime(time.Duration(cfg.DBMaxConnLifeTime) * time.Minute)

	if err = psqlDB.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к БД: ", err)
	}

	log.Println("Подключение к БД успешно создано")
	return psqlDB
}
