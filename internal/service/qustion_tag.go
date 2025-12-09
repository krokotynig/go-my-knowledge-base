package service

import (
	"database/sql"
	"fmt"
	"knowledge-base/internal/models"
)

// Структура для работы со всеми ф-ями service/qustion_tag.go.
type QuestionTagService struct {
	db *sql.DB
}

// Фунция  для создания объекта типа QuestionTagService.
func NewQuestionTagService(db *sql.DB) *QuestionTagService {
	return &QuestionTagService{db: db}
}

func (questionTagService *QuestionTagService) AddToQuestion(questionID, tagID int) error {

	//Создание sql запроса для прикрепления тега к вопросу.
	var query string = `insert into questions_tags (question_id, tag_id) values ($1, $2)`

	// Выполнение функции, которая проводит sql запрос без возврата данных.
	_, err := questionTagService.db.Exec(query, questionID, tagID)

	return err
}

func (questionTagService *QuestionTagService) GetAllRelations() ([]models.QuestionTag, error) {

	//Создание sql запроса для получения данных по всем связям.
	var query string = `select question_id, tag_id from questions_tags order by question_id, tag_id`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу с несколькими строками.
	rows, err := questionTagService.db.Query(query)
	if err != nil {
		return nil, err
	}

	// Закрытие, чтобы не было утечки соединений, надо закрыть.
	defer rows.Close()

	var relations []models.QuestionTag

	// Запись полученных данных из БД в массив формата []models.QuestionTag.
	for rows.Next() {
		var relation models.QuestionTag
		err := rows.Scan(&relation.QuestionID, &relation.TagID)
		if err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}

	return relations, nil
}

func (questionTagService *QuestionTagService) GetAllRelationsByTagID(tagID int) ([]models.QuestionTag, error) {

	//Создание sql запроса для получения данных по всем связям.
	var query string = `select tag_id, question_id from questions_tags where tag_id = $1`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу с несколькими строками.
	rows, err := questionTagService.db.Query(query, tagID)
	if err != nil {
		return nil, err
	}

	// Закрытие, чтобы не было утечки соединений, надо закрыть.
	defer rows.Close()

	var relations []models.QuestionTag

	// Запись полученных данных из БД в массив формата []models.QuestionTag.
	for rows.Next() {
		var relation models.QuestionTag
		err := rows.Scan(&relation.TagID, &relation.QuestionID)
		if err != nil {
			return nil, err
		}
		relations = append(relations, relation)
	}

	return relations, nil
}

func (questionTagService *QuestionTagService) DeleteRelationByID(questionID int, tagID int) error {

	//Создание sql запроса для удаления данных одной конкретной связи.
	var queryDelete string = `delete from questions_tags where question_id = $1 and tag_id = $2`

	// Выполнение функции, которая проводит sql запрос без возврата данных.
	result, err := questionTagService.db.Exec(queryDelete, questionID, tagID)
	if err != nil {
		return err
	}

	// Выполнение функции, которая возаращает количество удаленных строк.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	//Проверка было ли удаление строки. Если rowsAffected = 0, то не было.
	if rowsAffected == 0 {
		return fmt.Errorf("question with id %d,%d not found", questionID, tagID)
	}
	return nil
}
