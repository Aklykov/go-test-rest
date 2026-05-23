package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port int
}

func LoadAppConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("ошибка загрузки .env: %w", err)
	}

	portStr := os.Getenv("APP_PORT")
	// Валидация порта (должно быть числом)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("неверный порт приложения: %w", err)
	}

	return &AppConfig{
		Port: port,
	}, nil
}

func (c AppConfig) GetPort() string {
	return fmt.Sprintf("%d", c.Port)
}

// DBConfig хранит настройки подключения к базе
type DBConfig struct {
	Host     string
	Port     int // int, а не string, для валидации
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Load загружает конфигурацию из .env или окружения
func LoadDBConfig() (*DBConfig, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("ошибка загрузки .env: %w", err)
	}

	// Считываем переменные
	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	// Валидация порта (должно быть числом)
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("неверный порт БД: %w", err)
	}

	return &DBConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
		SSLMode:  sslMode,
	}, nil
}

// DSN возвращает строку подключения для lib/pq
func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}
