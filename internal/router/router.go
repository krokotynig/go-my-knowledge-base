package router

import (
	"knowledge-base/internal/app"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"knowledge-base/internal/handler"
)

// Setup создает и настраивает роутер со всеми маршрутами.
func Setup(handlers *app.Handlers) *mux.Router {
	router := mux.NewRouter()

	// Базовые маршруты (статус)
	registerCommonRoutes(router)

	// API маршруты
	registerTutorRoutes(router, handlers.Tutor)
	registerQuestionRoutes(router, handlers.Question)
	registerAnswerRoutes(router, handlers.Answer)
	registerTagRoutes(router, handlers.Tag)
	registerQuestionVersionRoutes(router, handlers.QuestionVersion)
	registerAnswerVersionRoutes(router, handlers.AnswerVersion)
	registerQuestionTagRoutes(router, handlers.QuestionTag)
	registerSearchRoutes(router, handlers.SimpleSearch)

	// Документация
	registerSwaggerRoutes(router)

	return router
}

// registerCommonRoutes регистрирует общие маршруты.
func registerCommonRoutes(router *mux.Router) {
	router.HandleFunc("/", handler.StatusHandler).Methods("GET")
	router.HandleFunc("/status", handler.StatusHandler).Methods("GET")
}

// Регистрирует маршруты для тьюторов.
func registerTutorRoutes(router *mux.Router, handler *handler.TutorHandler) {
	subrouter := router.PathPrefix("/tutors").Subrouter()

	subrouter.HandleFunc("", handler.GetAllTutors).Methods("GET")
	subrouter.HandleFunc("/{id}", handler.GetTutorByID).Methods("GET")
	subrouter.HandleFunc("/{id}", handler.DeleteTutorByID).Methods("DELETE")
	subrouter.HandleFunc("", handler.PostTutorString).Methods("POST")
	subrouter.HandleFunc("/{id}", handler.PutTutorString).Methods("PUT")
}

// Регистрирует маршруты для вопросов.
func registerQuestionRoutes(router *mux.Router, handler *handler.QuestionHandler) {
	subrouter := router.PathPrefix("/questions").Subrouter()

	subrouter.HandleFunc("", handler.GetAllQuestions).Methods("GET")
	subrouter.HandleFunc("/{id}", handler.GetQuestionByID).Methods("GET")
	subrouter.HandleFunc("/{id}/delete-by-tutor/{tutor_id}", handler.DeleteQuestionByID).Methods("DELETE")
	subrouter.HandleFunc("", handler.PostQuestionString).Methods("POST")
	subrouter.HandleFunc("/{id}", handler.PutQuestionString).Methods("PUT")
}

// Регистрирует маршруты для ответов.
func registerAnswerRoutes(router *mux.Router, handler *handler.AnswerHandler) {
	subrouter := router.PathPrefix("/answers").Subrouter()

	subrouter.HandleFunc("", handler.GetAllAnswers).Methods("GET")
	subrouter.HandleFunc("/{id}", handler.GetAnswerByID).Methods("GET")
	subrouter.HandleFunc("/{id}/delete-by-tutor/{tutor_id}", handler.DeleteAnswerByID).Methods("DELETE")
	subrouter.HandleFunc("", handler.PostAnswerString).Methods("POST")
	subrouter.HandleFunc("/{id}", handler.PutAnswerString).Methods("PUT")
}

// Регистрирует маршруты для тегов.
func registerTagRoutes(router *mux.Router, handler *handler.TagHandler) {
	subrouter := router.PathPrefix("/tags").Subrouter()

	subrouter.HandleFunc("", handler.GetAllTags).Methods("GET")
	subrouter.HandleFunc("/{id}", handler.GetTagByID).Methods("GET")
	subrouter.HandleFunc("/name/{name}", handler.GetTagByName).Methods("GET")
	subrouter.HandleFunc("/{id}", handler.DeleteTagByID).Methods("DELETE")
	subrouter.HandleFunc("", handler.PostTagString).Methods("POST")
}

// Регистрирует маршруты для версий вопросов.
func registerQuestionVersionRoutes(router *mux.Router, h *handler.QuestionVersionHandler) {
	router.HandleFunc("/question-versions/{id}", h.GetAllQuestionVersionsByID).Methods("GET")
}

// Регистрирует регистрирует маршруты для версий ответов.
func registerAnswerVersionRoutes(router *mux.Router, handler *handler.AnswerVersionHandler) {
	router.HandleFunc("/answer-versions/{id}", handler.GetAllAnswerVersionsByID).Methods("GET")
}

// Регистрирует регистрирует маршруты для связи вопросов и тегов.
func registerQuestionTagRoutes(router *mux.Router, handler *handler.QuestionTagHandler) {
	subrouter := router.PathPrefix("/question-tags").Subrouter()

	subrouter.HandleFunc("/{question_id}/{tag_id}", handler.AddTagToQuestion).Methods("POST")
	subrouter.HandleFunc("", handler.GetAllQuestionTagRelations).Methods("GET")
	subrouter.HandleFunc("/by-tag/{tag_id}", handler.GetAllQuestionTagRelationsByTagID).Methods("GET")
	subrouter.HandleFunc("/{question_id}/{tag_id}", handler.DeleteQuestionTagRelationByID).Methods("DELETE")
}

// Регистрирует регистрирует маршруты поиска.
func registerSearchRoutes(router *mux.Router, handler *handler.SimpleSearchHandler) {
	router.HandleFunc("/simple-search/{name}", handler.SearchHandler).Methods("GET")
}

// Регистрирует регистрирует маршруты для Swagger.
func registerSwaggerRoutes(route *mux.Router) {
	route.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
