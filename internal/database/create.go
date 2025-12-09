package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func CreateTables(db *sql.DB) error {
	// Читаем SQL из файла миграций.

	sqlBytes, err := os.ReadFile("./migrations/001_create_tables.sql")
	if err != nil {
		// Пробуем путь для локальной разработки.
		sqlBytes, err = os.ReadFile("../../migrations/001_create_tables.sql")
		if err != nil {
			return fmt.Errorf("ошибка чтения файла миграции: %v", err)
		}
	}

	// Выполняем SQL запросы.
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("ошибка создания таблиц: %v", err)
	}

	log.Println("✅ Таблицы созданы")
	return nil
}
