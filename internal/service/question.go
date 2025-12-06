package service

import (
	"database/sql"
	"fmt"
	"knowledge-base/internal/models"
)

type QuestionService struct {
	db *sql.DB
}

func NewQuestionServicer(db *sql.DB) *QuestionService {
	return &QuestionService{db: db}
}

func (questionService *QuestionService) GetAll() ([]models.Question, error) {
	var query string = `select id, question_text, tutor_id, created_at, is_edit from questions order by id`

	rows, err := questionService.db.Query(query) // для получения нескольких строк
	if err != nil {
		return nil, err
	}
	defer rows.Close() // вроде как, чтобы не было уттечки соединений, надо закрыть. Держет Бд открытой пока читаешь

	var questions []models.Question
	for rows.Next() {
		var question models.Question
		err := rows.Scan(&question.ID, &question.QuestionText, &question.TutorID, &question.CreatedAt, &question.IsEdit) // заполнения полей данынми
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil // slice ссылочный тип

}

func (questionService *QuestionService) GetByID(id int) (models.Question, error) {
	var query string = `select id, question_text, tutor_id, created_at, is_edit from questions where id = $1` // $1 - placeholder

	row := questionService.db.QueryRow(query, id) // для получения одной строки. Вроде как автоматически закрывает соеединение.
	var question models.Question

	err := row.Scan(&question.ID, &question.QuestionText, &question.TutorID, &question.CreatedAt, &question.IsEdit) // возвращает ошибку при select и handler видитЮ что нужно отдать текст на клиент
	if err != nil {
		return models.Question{}, err
	}

	return question, nil
}

func (questionService *QuestionService) DeleteByID(id int) error {
	query := `delete from questions where id = $1`

	result, err := questionService.db.Exec(query, id) // Exec длля операций не возвращающих данные
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected() // возвращение кол-ва удаленных строк, проверяет, что запись существовала. На клиенте не нужно. Для 404
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("question with id %d not found", id) // для того, чтобы не вернуть err, которая nil будет, для того, чтобы можно было лооги внутри сервера посмотреть. Будет nil потому что удаление id, которого нет все равно проходит.
	}

	return nil
}

func (questionService *QuestionService) PostString(questionText string, tutorId *int, isEdit bool) (int, error) {
	query := `insert into questions (question_text, tutor_id, is_edit) 
              values ($1, $2, $3) returning id`

	var id int

	row := questionService.db.QueryRow(query, questionText, tutorId, isEdit)

	err := row.Scan(&id) // заполним id
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (questionService *QuestionService) PutString(questionText string, tutorId *int, isEdit bool, id int) (models.Question, error) {
	query := `update questions 
              set question_text = $1, tutor_id = $2, is_edit = $3
              where id = $4
              returning question_text, tutor_id, created_at, is_edit`

	var question models.Question
	err := questionService.db.QueryRow(
		query, questionText, tutorId, isEdit, id).Scan(&question.QuestionText, &question.TutorID, &question.CreatedAt, &question.IsEdit)

	if err != nil {
		return models.Question{}, err
	}

	question.ID = id //Мутки мутные. надо думать

	return question, nil
}
