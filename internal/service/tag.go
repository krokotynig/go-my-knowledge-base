package service

import (
	"database/sql"
	"fmt"
	"knowledge-base/internal/models"
)

// Структура для работы со всеми ф-ями service/tag.go.
type TagService struct {
	db *sql.DB
}

// Фунция для создания объекта типа TagService.
func NewTagService(db *sql.DB) *TagService {
	return &TagService{db: db}
}

func (tagService *TagService) GetAll() ([]models.Tag, error) {

	//Создание sql запроса для получения данных по всем тегам.
	var query string = `select id, tutor_id, tag from tags order by tag`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу с несколькими строками.
	rows, err := tagService.db.Query(query)
	if err != nil {
		return nil, err
	}

	// Закрытие, чтобы не было утечки соединений, надо закрыть.
	defer rows.Close()

	var tags []models.Tag

	// Запись полученных данных из БД в массив формата []models.Tag.
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.ID, &tag.TutorID, &tag.Tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (tagService *TagService) GetByID(id int) (models.Tag, error) {

	//Создание sql запроса для получения данных по одному конкретному тегу.
	var query string = `select id, tutor_id, tag from tags where id = $1`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу из одной строки.
	row := tagService.db.QueryRow(query, id)

	var tag models.Tag

	// Запись полученных данных из БД в перемнную типа models.Tag.
	err := row.Scan(&tag.ID, &tag.TutorID, &tag.Tag)
	if err != nil {
		return models.Tag{}, err
	}

	return tag, nil
}

func (tagService *TagService) DeleteByID(id int) error {

	//Создание sql запроса для удаления данных одного кокретного тега.
	query := `delete from tags where id = $1`

	// Выполнение функции, которая проводит sql запрос без возврата данных.
	result, err := tagService.db.Exec(query, id)
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
		return fmt.Errorf("tag with id %d not found", id)
	}

	return nil
}

func (tagService *TagService) PostString(tag string, tutorID *int) (int, error) {

	//Создание sql запроса для появления новой записи в таблице тегов.
	query := `insert into tags (tag, tutor_id) values
			($1,$2) returning id`

	var id int

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу из одной строки.
	row := tagService.db.QueryRow(query, tag, tutorID)

	// Получение id созданной записи.
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
