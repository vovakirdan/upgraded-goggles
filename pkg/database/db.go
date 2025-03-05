package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

// DB является глобальным пулом подключений к базе данных
var DB *sql.DB

// InitDB инициализирует подключение к базе данных.
// dataSourceName содержит строку подключения:
// "postgres://user:password@localhost:5432/dbname?sslmode=disable"
func InitDB(dataSourceName string) error {
	var err error

	// Открываем подключение к базе данных
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Проверяем подключение
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}
