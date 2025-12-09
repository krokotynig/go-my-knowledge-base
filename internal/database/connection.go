package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // драйвер дработы c postgres. Не нужен при работе с "database/sql" в других пакетах.
)

func Connect() *sql.DB {

	if os.Getenv("DB_HOST") != "knowledge_db" {
		// Локальная разработка - грузим .env.
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Локальная разработка: не найден .env файл")
		}
	}

	// Получение env в память переменных без hardcode.
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost" // fallback для локальной разработки.
	}

	var dbName string = os.Getenv("POSTGRES_DB")
	var dbUser string = os.Getenv("POSTGRES_USER")
	var dbPassword string = os.Getenv("POSTGRES_PASSWORD")
	var dbPort string = os.Getenv("POSTGRES_PORT")

	if dbPort == "" {
		dbPort = "5432"
	}

	// Url, строка для подключения к БД. Потом ее парсит "database/sql".
	var dsn string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	//db это результат работы метода Open. Он как бы создает DB. Ключевая сущность проекта, нужная для работы с BD.
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	// Ping проверяет, что соединение с базой данных все еще работает, и при необходимости устанавливает соединение.
	err = db.Ping()
	if err != nil {
		log.Fatal("БД не отвечает:", err)
	}

	fmt.Println("✅ Успешно подключились к PostgreSQL!")

	return db
}
