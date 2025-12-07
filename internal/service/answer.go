package service

import (
	"database/sql"
	"fmt"
	"knowledge-base/internal/models"
)

// Структура для работы со всеми ф-ями service/answer.go.
type AnswerService struct {
	db *sql.DB
}

// Фунция  для создания объекта типа AnswerService.
func NewAnswerService(db *sql.DB) *AnswerService {
	return &AnswerService{db: db}
}

func (answerService *AnswerService) GetAll() ([]models.Answer, error) {

	//Создание sql запроса для получения данных по всем ответам.
	var query string = `select id, answer_text, tutor_id, question_id, created_at, is_edit from answers order by id`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу с несколькими строками.
	rows, err := answerService.db.Query(query)
	if err != nil {
		return nil, err
	}

	// Закрытие, чтобы не было уттечки соединений, надо закрыть.
	defer rows.Close()

	// Запись полученных данных из БД в массив формата []models.Answer.
	var answers []models.Answer
	for rows.Next() {
		var answer models.Answer
		err := rows.Scan(&answer.ID, &answer.AnswersText, &answer.TutorID, &answer.QuestionID, &answer.CreatedAt, &answer.IsEdit)
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil

}

func (answerService *AnswerService) GetByID(id int) (models.Answer, error) {

	//Создание sql запроса для получения данных по одному конкретному ответу.
	var query string = `select id, answer_text, tutor_id, question_id, created_at, is_edit from answers where id = $1`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу из одной строки.
	row := answerService.db.QueryRow(query, id)

	var answer models.Answer

	// Запись полученных данных из БД в перемнную типа models.Answer.
	err := row.Scan(&answer.ID, &answer.AnswersText, &answer.TutorID, &answer.QuestionID, &answer.CreatedAt, &answer.IsEdit)
	if err != nil {
		return models.Answer{}, err
	}

	return answer, nil
}

func (answerService *AnswerService) DeleteByID(id int) error {

	//Создание sql запроса для удаления данных одного кокретного овтета.
	query := `delete from answers where id = $1`

	// Выполнение функции, которая проводит sql запрос без возврата данных.
	result, err := answerService.db.Exec(query, id)
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
		return fmt.Errorf("answer with id %d not found", id)
	}

	return nil
}

func (answerService *AnswerService) PostString(answerText string, tutorId *int, questionId int) (int, error) {

	//Создание sql запроса для появления новой записи в таблице овтетов.
	query := `insert into answers (answer_text, tutor_id, question_id) 
              values ($1, $2, $3) returning id`

	var id int

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу из одной строки.
	row := answerService.db.QueryRow(query, answerText, tutorId, questionId)

	// Получение id созданной записи.
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (answerService *AnswerService) PutString(answerText string, tutorId *int, questionId int, isEdit bool, id int) (models.Answer, error) {

	//Создание sql запроса для обновления данных конкретного вопроса.
	query := `update answers 
              set answer_text = $1, tutor_id = $2, question_id = $3, is_edit = $4
              where id = $5
              returning answer_text, tutor_id, question_id, created_at, is_edit`

	var answer models.Answer

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу из одной строки. Заполнение полей переменной типа models.Answer.
	err := answerService.db.QueryRow(
		query, answerText, tutorId, questionId, isEdit, id).Scan(&answer.AnswersText, &answer.TutorID, &answer.QuestionID, &answer.CreatedAt, &answer.IsEdit)

	if err != nil {
		return models.Answer{}, err
	}

	return answer, nil
}
