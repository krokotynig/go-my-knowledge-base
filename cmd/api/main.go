package main

import (
	_ "knowledge-base/docs" // импорт docs
	"knowledge-base/internal/database"
	"knowledge-base/internal/handler"
	"knowledge-base/internal/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Knowledge Base API
// @version 1.0
// @description API для базы знаний с вопросами и ответами
// @host localhost:2709
// @BasePath /
func main() {

	db := database.Connect()
	// закроет при выходе из программы
	defer db.Close() // defer отложенное действие к концу ф-и Close просто закрывает базу данных

	tutorService := service.NewTutor(db)
	tutorHandler := handler.NewTutorhandler(tutorService)
	questionService := service.NewQuestionServicer(db)
	questionHandler := handler.NewQuestionHandler(questionService)

	r := mux.NewRouter() // создание новго явновго роутера из пакета gorilla/mux

	r.HandleFunc("/", handler.StatusHandler).Methods("GET")       // регистрация маршрутов в дефолтном роутере вроде , healthcheack
	r.HandleFunc("/status", handler.StatusHandler).Methods("GET") // регистрация маршрутов , healthcheack

	r.HandleFunc("/tutors", tutorHandler.GetAllTutors).Methods("GET")
	r.HandleFunc("/tutors/{id}", tutorHandler.GetTutorByID).Methods("GET")
	r.HandleFunc("/tutors/{id}", tutorHandler.DeleteTutorByID).Methods("DELETE") // в REST операции определяются HTTP методами, а не путями
	r.HandleFunc("/tutors", tutorHandler.PostTutorString).Methods("POST")
	r.HandleFunc("/tutors/{id}", tutorHandler.PutTutorString).Methods("PUT")

	r.HandleFunc("/questions", questionHandler.GetAllQuestions).Methods("GET")
	r.HandleFunc("/questions/{id}", questionHandler.GetQuestionByID).Methods("GET")
	r.HandleFunc("/questions/{id}", questionHandler.DeleteQuestionByID).Methods("DELETE") // в REST операции определяются HTTP методами, а не путями
	r.HandleFunc("/questions", questionHandler.PostQuestionString).Methods("POST")
	r.HandleFunc("/questions/{id}", questionHandler.PutQuestionString).Methods("PUT")

	r.HandleFunc("/swagger/{any}", httpSwagger.WrapHandler).Methods("GET")

	err := database.RunMigrations(db) // вызов методов миграций через координатор
	if err != nil {
		log.Fatal("Ошибка миграций:", err)

	}

	log.Println(" ✅ База данных готова!")
	log.Println("🚀 Запуск сервера на http://localhost:2709")

	log.Println("📚 Swagger UI доступен на http://localhost:2709/swagger/index.html")

	err = http.ListenAndServe(":2709", r) // блокирующая функция, после нее программа ждет http запросы
	if err != nil {
		log.Printf("❌ Ошибка запускас ервера: %v\n", err)
	}
}
