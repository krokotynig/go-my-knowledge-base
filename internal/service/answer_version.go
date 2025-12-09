package service

import (
	"database/sql"
	"knowledge-base/internal/models"
)

// Структура для работы со всеми ф-ями service/answer_version.go.
type AnswerVersionService struct {
	db *sql.DB
}

// Функция для создания объекта типа AnswerVersionService.
func NewAnswerVersionService(db *sql.DB) *AnswerVersionService {
	return &AnswerVersionService{db: db}
}

func (answerVersionService *AnswerVersionService) GetAllByID(id int) ([]models.AnswerVersion, error) {

	// Создание sql запроса для получения данных о версиях конкретного ответа.
	var query string = `select id, answer_id, answer_text, question_id, tutor_id, created_at, version_number, is_delete, delete_by_tutor from answer_versions where answer_id = $1 order by version_number`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу с несколькими строками.
	rows, err := answerVersionService.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	// Закрытие, чтобы не было утечки соединений, надо закрыть.
	defer rows.Close()

	var answerVersions []models.AnswerVersion

	// Запись полученных данных из БД в массив формата []models.AnswerVersion.
	for rows.Next() {
		var answerVersion models.AnswerVersion
		err := rows.Scan(&answerVersion.ID, &answerVersion.AnswerID, &answerVersion.AnswerText, &answerVersion.QuestionID, &answerVersion.TutorID, &answerVersion.CreatedAt, &answerVersion.VersionNumber, &answerVersion.IsDelete, &answerVersion.DeleteByTutor)
		if err != nil {
			return nil, err
		}
		answerVersions = append(answerVersions, answerVersion)
	}

	return answerVersions, nil
}
