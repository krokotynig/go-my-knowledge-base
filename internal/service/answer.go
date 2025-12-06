package service

import (
	"database/sql"
	"fmt"
	"knowledge-base/internal/models"
)

type AnswerService struct {
	db *sql.DB
}

func NewAnswerService(db *sql.DB) *AnswerService {
	return &AnswerService{db: db}
}

func (answerService *AnswerService) GetAll() ([]models.Answer, error) {
	var query string = `select id, answer_text, tutor_id, question_id, created_at, is_edit from answers order by id`

	rows, err := answerService.db.Query(query) // для получения нескольких строк
	if err != nil {
		return nil, err
	}
	defer rows.Close() // вроде как, чтобы не было уттечки соединений, надо закрыть. Держет Бд открытой пока читаешь

	var answers []models.Answer
	for rows.Next() {
		var answer models.Answer
		err := rows.Scan(&answer.ID, &answer.AnswersText, &answer.TutorID, &answer.QuestionID, &answer.CreatedAt, &answer.IsEdit) // заполнения полей данынми
		if err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil // slice ссылочный тип

}

func (answerService *AnswerService) GetByID(id int) (models.Answer, error) {
	var query string = `select id, answer_text, tutor_id, question_id, created_at, is_edit from answers where id = $1` // $1 - placeholder

	row := answerService.db.QueryRow(query, id) // для получения одной строки. Вроде как автоматически закрывает соеединение.
	var answer models.Answer

	err := row.Scan(&answer.ID, &answer.AnswersText, &answer.TutorID, &answer.QuestionID, &answer.CreatedAt, &answer.IsEdit) // возвращает ошибку при select и handler видитЮ что нужно отдать текст на клиент
	if err != nil {
		return models.Answer{}, err
	}

	return answer, nil
}

func (answerService *AnswerService) DeleteByID(id int) error {
	query := `delete from answers where id = $1`

	result, err := answerService.db.Exec(query, id) // Exec длля операций не возвращающих данные
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected() // возвращение кол-ва удаленных строк, проверяет, что запись существовала. На клиенте не нужно. Для 404
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("answer with id %d not found", id) // для того, чтобы не вернуть err, которая nil будет, для того, чтобы можно было лооги внутри сервера посмотреть. Будет nil потому что удаление id, которого нет все равно проходит.
	}

	return nil
}

func (answerService *AnswerService) PostString(answerText string, tutorId *int, questionId int) (int, error) {
	query := `insert into answers (answer_text, tutor_id, question_id) 
              values ($1, $2, $3) returning id` //без is_edit, надо чтобы он fals был при создании

	var id int
	row := answerService.db.QueryRow(query, answerText, tutorId, questionId)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (answerService *AnswerService) PutString(answerText string, tutorId *int, questionId int, isEdit bool, id int) (models.Answer, error) {
	// тут с is_edit, чтобы было true по умолчанию, по умолчанию в теле запроса
	query := `update answers 
              set answer_text = $1, tutor_id = $2, question_id = $3, is_edit = $4
              where id = $5
              returning answer_text, tutor_id, question_id, created_at, is_edit`

	var answer models.Answer
	err := answerService.db.QueryRow(
		query, answerText, tutorId, questionId, isEdit, id).Scan(&answer.AnswersText, &answer.TutorID, &answer.QuestionID, &answer.CreatedAt, &answer.IsEdit)

	if err != nil {
		return models.Answer{}, err
	}

	answer.ID = id //Мутки мутные. Надо думать

	return answer, nil
}
