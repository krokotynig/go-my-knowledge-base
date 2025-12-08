package main

import (
	_ "knowledge-base/docs"
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

	// Закрытие соединения при выходе из программы.
	defer db.Close()

	// Создание объектов, отвечающих за работу получения и добовалнея данных в БД - service и за логику handlers.
	tutorService := service.NewTutor(db)
	tutorHandler := handler.NewTutorhandler(tutorService)

	questionService := service.NewQuestionService(db)
	questionHandler := handler.NewQuestionHandler(questionService)

	answerService := service.NewAnswerService(db)
	answerHandler := handler.NewAnswerHandler(answerService)

	tagService := service.NewTagService(db)
	tagHandler := handler.NewTagHandler(tagService)

	questionVersionService := service.NewQuestionVersionService(db)
	questionVersionHandler := handler.NewQuestionVersionHandler(questionVersionService)

	answerVersionService := service.NewAnswerVersionService(db)
	answerVersionHandler := handler.NewAnswerVersionHandler(answerVersionService)

	questionTagService := service.NewQuestionTagService(db)
	questionTagHandler := handler.NewQuestionTagHandler(questionTagService)

	// в REST операции определяются HTTP методами, а не путями.

	//Создание явновго роутера из пакета gorilla/mux.
	r := mux.NewRouter()

	//Регистрация базового маршрута.
	r.HandleFunc("/", handler.StatusHandler).Methods("GET")
	r.HandleFunc("/status", handler.StatusHandler).Methods("GET")

	//Регистрация муршрута tutors.
	r.HandleFunc("/tutors", tutorHandler.GetAllTutors).Methods("GET")
	r.HandleFunc("/tutors/{id}", tutorHandler.GetTutorByID).Methods("GET")
	r.HandleFunc("/tutors/{id}", tutorHandler.DeleteTutorByID).Methods("DELETE")
	r.HandleFunc("/tutors", tutorHandler.PostTutorString).Methods("POST")
	r.HandleFunc("/tutors/{id}", tutorHandler.PutTutorString).Methods("PUT")

	//Регистрация муршрута questions.
	r.HandleFunc("/questions", questionHandler.GetAllQuestions).Methods("GET")
	r.HandleFunc("/questions/{id}", questionHandler.GetQuestionByID).Methods("GET")
	r.HandleFunc("/questions/{id}", questionHandler.DeleteQuestionByID).Methods("DELETE")
	r.HandleFunc("/questions", questionHandler.PostQuestionString).Methods("POST")
	r.HandleFunc("/questions/{id}", questionHandler.PutQuestionString).Methods("PUT")

	//Регистрация муршрута answers.
	r.HandleFunc("/answers", answerHandler.GetAllAnswers).Methods("GET")
	r.HandleFunc("/answers/{id}", answerHandler.GetAnswerByID).Methods("GET")
	r.HandleFunc("/answers/{id}", answerHandler.DeleteAnswerByID).Methods("DELETE")
	r.HandleFunc("/answers", answerHandler.PostAnswerString).Methods("POST")
	r.HandleFunc("/answers/{id}", answerHandler.PutAnswerString).Methods("PUT")

	//Регистрация муршрута tags.
	r.HandleFunc("/tags", tagHandler.GetAllTags).Methods("GET")
	r.HandleFunc("/tags/{id}", tagHandler.GetTagByID).Methods("GET")
	r.HandleFunc("/tags/{id}", tagHandler.DeleteTagByID).Methods("DELETE")
	r.HandleFunc("/tags", tagHandler.PostTagString).Methods("POST")

	//Регистрация маршрута question-versions.
	r.HandleFunc("/question-versions/{id}", questionVersionHandler.GetAllQuestionVersionsByID).Methods("GET")

	//Регистрация маршрута answer-versions.
	r.HandleFunc("/answer-versions/{id}", answerVersionHandler.GetAllAnswerVersionsByID).Methods("GET")

	//Регистрация маршрутов questions_tags.
	r.HandleFunc("/questions/{question_id}/tags/{tag_id}", questionTagHandler.AddTagToQuestion).Methods("POST")

	//Регистрация муршрута swagger.
	r.HandleFunc("/swagger/{any}", httpSwagger.WrapHandler).Methods("GET")

	//Вызов методов миграций через координатор "RunMigrations.
	err := database.RunMigrations(db)
	if err != nil {
		log.Fatal("Ошибка миграций:", err)

	}

	log.Println(" ✅ База данных готова!")
	log.Println("🚀 Запуск сервера на http://localhost:2709")

	log.Println("📚 Swagger UI доступен на http://localhost:2709/swagger/index.html")

	//Блокирующая функция, после нее программа ждет http запросы.
	err = http.ListenAndServe(":2709", r)
	if err != nil {
		log.Printf("❌ Ошибка запускас ервера: %v\n", err)
	}
}
