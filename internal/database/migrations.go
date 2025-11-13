package database

import "database/sql"

func RunMigrations(db *sql.DB) error { // координатор
	if err := CreateTables(db); err != nil {
		return err
	}

	if err := SeedData(db); err != nil {
		return err
	}

	return nil
}
