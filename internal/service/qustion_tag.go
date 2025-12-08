package service

import (
	"database/sql"
)

// Структура для работы со всеми ф-ями service/qustion_tag.go.
type QuestionTagService struct {
	db *sql.DB
}

// Фунция  для создания объекта типа QuestionTagService.
func NewQuestionTagService(db *sql.DB) *QuestionTagService {
	return &QuestionTagService{db: db}
}

func (s *QuestionTagService) AddToQuestion(questionID, tagID int) error {

	// Выполнение функции, которая проводит sql запрос без возврата данных.
	_, err := s.db.Exec(`insert into questions_tags (question_id, tag_id) values ($1, $2)`, questionID, tagID)

	return err
}
