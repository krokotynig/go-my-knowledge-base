package app

import (
	"database/sql"
	"knowledge-base/internal/handler"
	"knowledge-base/internal/service"
)

// Services содержит все сервисы.
type services struct {
	Tutor           *service.TutorService
	Question        *service.QuestionService
	Answer          *service.AnswerService
	Tag             *service.TagService
	QuestionVersion *service.QuestionVersionService
	AnswerVersion   *service.AnswerVersionService
	QuestionTag     *service.QuestionTagService
	SimpleSearch    *service.SimpleSearchService
}

// Handlers содержит все хэндлеры.
type Handlers struct {
	Tutor           *handler.TutorHandler
	Question        *handler.QuestionHandler
	Answer          *handler.AnswerHandler
	Tag             *handler.TagHandler
	QuestionVersion *handler.QuestionVersionHandler
	AnswerVersion   *handler.AnswerVersionHandler
	QuestionTag     *handler.QuestionTagHandler
	SimpleSearch    *handler.SimpleSearchHandler
}

// Создает и инициализирует все зависимости.
func NewContainer(db *sql.DB) *Handlers {

	// Инициализация всех сервисов.
	services := services{
		Tutor:           service.NewTutor(db),
		Question:        service.NewQuestionService(db),
		Answer:          service.NewAnswerService(db),
		Tag:             service.NewTagService(db),
		QuestionVersion: service.NewQuestionVersionService(db),
		AnswerVersion:   service.NewAnswerVersionService(db),
		QuestionTag:     service.NewQuestionTagService(db),
		SimpleSearch:    service.NewSimpleSearchService(db),
	}

	// Инициализация всех хэндлеров с соответствующими сервисами.
	handlers := &Handlers{
		Tutor:           handler.NewTutorhandler(services.Tutor),
		Question:        handler.NewQuestionHandler(services.Question),
		Answer:          handler.NewAnswerHandler(services.Answer),
		Tag:             handler.NewTagHandler(services.Tag),
		QuestionVersion: handler.NewQuestionVersionHandler(services.QuestionVersion),
		AnswerVersion:   handler.NewAnswerVersionHandler(services.AnswerVersion),
		QuestionTag:     handler.NewQuestionTagHandler(services.QuestionTag),
		SimpleSearch:    handler.NewSimpleSearchHandler(services.SimpleSearch),
	}

	return handlers
}
