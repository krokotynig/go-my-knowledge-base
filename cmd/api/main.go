package main

import (
	_ "knowledge-base/docs"
	"knowledge-base/internal/app"
	"knowledge-base/internal/database"
	"knowledge-base/internal/router"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// @title Knowledge Base API 📚
// @version 1.0
// @description API для базы знаний с вопросами и ответами
// @host localhost:2709
// @BasePath /
func main() {

	// Подключение к базе данных.
	db := database.Connect()
	defer db.Close()

	// Выполнение миграций.
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Ошибка миграций:", err)
	}

	// Создание контейнера зависимостей.
	container := app.NewContainer(db)

	// Настройка маршрутизатора.
	router := router.Setup(container)

	// Запуск сервера.
	log.Println(" ✅ База данных готова!")
	log.Println(" ✅ API готово!")
	log.Println("🚀 Запуск сервера на http://localhost:2709")
	log.Println("📚 Swagger UI доступен на http://localhost:2709/swagger/index.html")

	if err := http.ListenAndServe(":2709", router); err != nil {
		log.Printf("❌ Ошибка запуска сервера: %v\n", err)
	}
}
