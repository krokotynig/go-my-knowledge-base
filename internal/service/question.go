package service

import (
	"database/sql"
	"fmt"
	"knowledge-base/internal/models"
)

// Структура для работы со всеми ф-ями service/question.go.
type QuestionService struct {
	db *sql.DB
}

// Фунция  для создания объекта типа QuestionService.
func NewQuestionService(db *sql.DB) *QuestionService {
	return &QuestionService{db: db}
}

func (questionService *QuestionService) GetAll() ([]models.Question, error) {

	//Создание sql запроса для получения данных по всем вопросам.
	var query string = `select id, question_text, tutor_id, created_at, is_edit from questions order by id`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу с несколькими строками.
	rows, err := questionService.db.Query(query)
	if err != nil {
		return nil, err
	}

	// Закрытие, чтобы не было уттечки соединений, надо закрыть.
	defer rows.Close()

	var questions []models.Question

	// Запись полученных данных из БД в массив формата []models.Question.
	for rows.Next() {
		var question models.Question
		err := rows.Scan(&question.ID, &question.QuestionText, &question.TutorID, &question.CreatedAt, &question.IsEdit)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil

}

func (questionService *QuestionService) GetByID(id int) (models.Question, error) {

	//Создание sql запроса для получения данных по одному конкретному вопросу.
	var query string = `select id, question_text, tutor_id, created_at, is_edit from questions where id = $1`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу из одной строки.
	row := questionService.db.QueryRow(query, id)
	var question models.Question

	// Запись полученных данных из БД в перемнную типа models.Question.
	err := row.Scan(&question.ID, &question.QuestionText, &question.TutorID, &question.CreatedAt, &question.IsEdit)
	if err != nil {
		return models.Question{}, err
	}

	return question, nil
}

func (questionService *QuestionService) DeleteByID(id int) error {

	//Создание sql запроса для удаления данных одного кокретного вопроса.
	query := `delete from questions where id = $1`

	// Выполнение функции, которая проводит sql запрос без возврата данных.
	result, err := questionService.db.Exec(query, id)
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
		return fmt.Errorf("question with id %d not found", id)
	}

	return nil
}

func (questionService *QuestionService) PostString(questionText string, tutorId *int) (int, error) {

	//Создание sql запроса для появления новой записи в таблице вопросов.
	query := `insert into questions (question_text, tutor_id) 
              values ($1, $2) returning id`

	var id int

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу из одной строки.
	row := questionService.db.QueryRow(query, questionText, tutorId)

	// Получение id созданной записи.
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (questionService *QuestionService) PutString(questionText string, tutorId *int, isEdit bool, id int) (models.Question, error) {

	//Создание sql запроса для обновления данных конкретного вопроса.
	query := `update questions 
              set question_text = $1, tutor_id = $2, is_edit = $3
              where id = $4
              returning question_text, tutor_id, created_at, is_edit`

	var question models.Question

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу из одной строки. Заполнение полей переменной типа models.Question.
	err := questionService.db.QueryRow(
		query, questionText, tutorId, isEdit, id).Scan(&question.QuestionText, &question.TutorID, &question.CreatedAt, &question.IsEdit)

	if err != nil {
		return models.Question{}, err
	}

	return question, nil
}
