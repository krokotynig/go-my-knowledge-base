package main

import (
	"knowledge-base/internal/database"
	"knowledge-base/internal/handler"
	"knowledge-base/internal/service"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	db := database.Connect()
	// закроет при выходе из программы
	defer db.Close() // defer отложенное действие к концу ф-и Close просто закрывает базу данных

	tutorService := service.NewTutor(db)
	tutorHandler := handler.NewTutorhandler(tutorService)

	http.HandleFunc("/", handler.StatusHandler)       // регистрация маршрутов в дефолтном роутере вроде , healthcheack
	http.HandleFunc("/status", handler.StatusHandler) // регистрация маршрутов , healthcheack
	http.HandleFunc("/tutors", tutorHandler.GetAllTutors)
	http.HandleFunc("/tutors/", tutorHandler.GetTutorByID)
	http.HandleFunc("/tutors/delete/", tutorHandler.DeleteTutorByID) // в REST операции определяются HTTP методами, а не путями

	err := database.RunMigrations(db) // вызов методов миграций через координатор
	if err != nil {
		log.Fatal("Ошибка миграций:", err)

	}

	log.Println(" ✅ База данных готова!")
	log.Println("🚀 Запуск сервера на http://localhost:2709")

	err = http.ListenAndServe(":2709", nil) // блокирующая функция, после нее программа ждет http запросы
	if err != nil {
		log.Printf("❌ Ошибка запускас ервера: %v\n", err)
	}
}
