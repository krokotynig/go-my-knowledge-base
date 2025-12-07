package service

import (
	"database/sql"
	"fmt"
	"knowledge-base/internal/models"
)

// Структура для работы со всеми ф-ями service/tutor.go.
type TutorService struct {
	db *sql.DB
}

// Фунция  для создания объекта типа TutorService.
func NewTutor(db *sql.DB) *TutorService {
	return &TutorService{db: db}
}

func (tutorService *TutorService) GetAll() ([]models.Tutor, error) {

	//Создание sql запроса для получения данных по всем тьюторам.
	var query string = `select id, full_name, email from tutors order by id`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу с несколькими строками.
	rows, err := tutorService.db.Query(query)
	if err != nil {
		return nil, err
	}

	// Закрытие, чтобы не было уттечки соединений, надо закрыть.
	defer rows.Close()

	var tutors []models.Tutor

	// Запись полученных данных из БД в массив формата []models.Tutor.
	for rows.Next() {
		var tutor models.Tutor
		err := rows.Scan(&tutor.ID, &tutor.FullName, &tutor.Email)
		if err != nil {
			return nil, err
		}
		tutors = append(tutors, tutor)
	}

	return tutors, nil
}

func (tutorService *TutorService) GetByID(id int) (models.Tutor, error) {

	//Создание sql запроса для получения данных по одному конкретному тьютору.
	var query string = `select id, full_name, email from tutors where id = $1`

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу из одной строки.
	row := tutorService.db.QueryRow(query, id)

	var tutor models.Tutor

	// Запись полученных данных из БД в перемнную типа models.Tutor.
	err := row.Scan(&tutor.ID, &tutor.FullName, &tutor.Email)
	if err != nil {
		return models.Tutor{}, err
	}

	return tutor, nil
}

func (tutorService *TutorService) DeleteByID(id int) error {

	//Создание sql запроса для удаления данных одного кокретного тьютора.
	query := `delete from tutors where id = $1`

	// Выполнение функции, которая проводит sql запрос без возврата данных.
	result, err := tutorService.db.Exec(query, id)
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
		return fmt.Errorf("tutor with id %d not found", id)
	}

	return nil
}

func (tutorService *TutorService) PostString(fullName string, email string) (int, error) {

	//Создание sql запроса для появления новой записи в таблице тьюторов.
	query := `insert into tutors (full_name, email) values
			($1,$2) returning id`

	var id int

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу из одной строки.
	row := tutorService.db.QueryRow(query, fullName, email)

	// Получение id созданной записи.
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (tutorService *TutorService) PutString(fullName string, email string, id int) (models.Tutor, error) {

	//Создание sql запроса для обновления данных конкретного тьютора.
	query := `update tutors 
			set 
			full_name = $1, email = $2
			where id = $3
			returning full_name, email`

	var tutor models.Tutor

	// Выполнение функции, которая проводит sql запрос и возвращает таблицу из одной строки. Заполнение полей переменной типа models.Tutor.
	err := tutorService.db.QueryRow(query, fullName, email, id).Scan(&tutor.FullName, &tutor.Email)
	if err != nil {
		return models.Tutor{}, err
	}

	return tutor, nil
}
