package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // драйвер
)

func Connect() *sql.DB {

	// типа процедура, но не процедура, тут все ф-я
	err := godotenv.Load("../../.env") // относительный
	if err != nil {
		log.Fatal("Ошибка загрузки .env", err)
	}

	var dbHost string = "localhost"
	var dbName string = os.Getenv("POSTGRES_DB")
	var dbUser string = os.Getenv("POSTGRES_USER")
	var dbPassword string = os.Getenv("POSTGRES_PASSWORD")
	var dbPort string = os.Getenv("POSTGRES_PORT")

	var dsn string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName) // url, строка подключения к БД  потом ее парсит "database/sql"

	//*sql.DB типа указатель на поля структуры, без создания ее копий
	db, err := sql.Open("postgres", dsn) //db это результат работы метода Open. Он как бы создает DB
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	// err уже объявлена, поэтому "="
	err = db.Ping() // Ping проверяет, что соединение с базой данных все еще работает, и при необходимости устанавливает соединение.
	if err != nil {
		log.Fatal("БД не отвечает:", err)
	}

	fmt.Println("✅ Успешно подключились к PostgreSQL!")

	return db
}
