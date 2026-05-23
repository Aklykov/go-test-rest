package db

import (
	"database/sql"
	"fmt"
	"go-rest-crud/config"
	"log"

	_ "github.com/lib/pq"
)

// NewPostgresDB инициализирует соединение с базой на основе конфигурации
func NewPostgresDB() (*sql.DB, error) {
	cfg, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("ошибка загрузки конфига: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть соединение: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе: %w", err)
	}

	return db, nil
}
