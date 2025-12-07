package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func SeedData(db *sql.DB) error {
	// Читаем SQL из файла миграций.
	sqlBytes, err := os.ReadFile("../../migrations/002_seed_data.sql")
	if err != nil {
		return fmt.Errorf("ошибка чтения файла сидов: %v", err)
	}

	// Выполняем SQL.
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("ошибка наполнения данными: %v", err)
	}

	log.Println("✅ Данные добавлены")
	return nil
}
