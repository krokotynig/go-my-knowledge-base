package database

import "database/sql"

// Создание координатора, кторый сразу вызовет 2 пакета запросов на создание таблиц и их заполнение.
func RunMigrations(db *sql.DB) error {
	if err := CreateTables(db); err != nil {
		return err
	}

	if err := SeedData(db); err != nil {
		return err
	}

	return nil
}
