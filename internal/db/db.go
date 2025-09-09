package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MAPiryazev/Wildberries_L0/internal/config"
	_ "github.com/lib/pq"
)

func InitPsqlDB(cfg *config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	psqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Ошибка подключения к БД: %w", err)
	}

	psqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	psqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	psqlDB.SetConnMaxLifetime(time.Duration(cfg.DBMaxConnLifeTime) * time.Minute)

	if err = psqlDB.Ping(); err != nil {
		psqlDB.Close()
		return nil, fmt.Errorf("Не удалось проверить соединение с БД: %w", err)
	}

	return psqlDB, nil
}
