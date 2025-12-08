package service

import (
	"database/sql"
	"knowledge-base/internal/models"
)

// Структура для работы со всеми ф-ями service/question-version.go.
type QuestionVersionService struct {
	db *sql.DB
}

// Фунция  для создания объекта типа QuestionVersionService.
func NewQuestionVersionService(db *sql.DB) *QuestionVersionService {
	return &QuestionVersionService{db: db}
}

func (questionVersionService *QuestionVersionService) GetAllByID(id int) ([]models.QuestionVersion, error) {

	//Создание sql запроса для получения данных о версиях конкретного вопроса.
	var query string = `select id, question_id, question_text, tutor_id, created_at, version_number, is_delete, delete_by_tutor from question_versions where question_id = $1 order by version_number`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу с несколькими строками.
	rows, err := questionVersionService.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	// Закрытие, чтобы не было утечки соединений, надо закрыть.
	defer rows.Close()

	var questionVersions []models.QuestionVersion

	// Запись полученных данных из БД в массив формата []models.QuestionVersion.
	for rows.Next() {
		var questionVersion models.QuestionVersion
		err := rows.Scan(&questionVersion.ID, &questionVersion.QuestionID, &questionVersion.QuestionText, &questionVersion.TutorID, &questionVersion.CreatedAt, &questionVersion.VersionNumber, &questionVersion.IsDelete, &questionVersion.DeleteByTutor)
		if err != nil {
			return nil, err
		}
		questionVersions = append(questionVersions, questionVersion)
	}

	return questionVersions, nil
}
