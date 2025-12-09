package service

import (
	"database/sql"
	"fmt"
	"knowledge-base/internal/models"
)

// Структура для работы со всеми ф-ями service/simple_search.go.
type SimpleSearchService struct {
	db *sql.DB
}

// Фунция  для создания объекта типа QuestionTagService.
func NewSimpleSearchService(db *sql.DB) *SimpleSearchService {
	return &SimpleSearchService{db: db}
}

func (simpleSearchService SimpleSearchService) SearchLogic(name string) ([]models.Question, error) {
	var query string = `select q.*
        from public.questions q
        inner join public.questions_tags qt on q.id = qt.question_id
        inner join public.tags t on qt.tag_id = t.id
        where lower(trim(t.tag)) = lower(trim($1))
        order by q.created_at desc`

	rows, err := simpleSearchService.db.Query(query, name)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска: %v", err)
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var question models.Question
		err := rows.Scan(
			&question.ID, &question.QuestionText,
			&question.TutorID, &question.CreatedAt, &question.IsEdit)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil

}
