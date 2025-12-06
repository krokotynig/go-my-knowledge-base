package service

import (
	"database/sql" // Postgre драйвер по всей видимости не нужен. Потому что он для подлючения, а не логики service
	"fmt"
	"knowledge-base/internal/models"
)

type Tutor struct {
	db *sql.DB
}

func NewTutor(db *sql.DB) *Tutor { // для Main
	return &Tutor{db: db}
}

func (tutorService *Tutor) GetAll() ([]models.Tutor, error) { //(s *Tutor) - ресивер, указатель, что метод принадлежит структуре, далее метод вызывается из нее
	var query string = `SELECT id, full_name, email FROM tutors ORDER BY id`

	rows, err := tutorService.db.Query(query) // для получения нескольких строк
	if err != nil {
		return nil, err
	}
	defer rows.Close() // вроде как, чтобы не было уттечки соединений, надо закрыть. Держет Бд открытой пока читаешь

	var tutors []models.Tutor
	for rows.Next() {
		var tutor models.Tutor
		err := rows.Scan(&tutor.ID, &tutor.FullName, &tutor.Email) // заполнения полей данынми
		if err != nil {
			return nil, err
		}
		tutors = append(tutors, tutor)
	}

	return tutors, nil // slice ссылочный тип
}

func (tutorService *Tutor) GetByID(id int) (models.Tutor, error) {
	var query string = `SELECT id, full_name, email FROM tutors WHERE id = $1` // $1 - placeholder

	row := tutorService.db.QueryRow(query, id) // для получения одной строки. Вроде как автоматически закрывает соеединение
	var tutor models.Tutor

	err := row.Scan(&tutor.ID, &tutor.FullName, &tutor.Email)
	if err != nil {
		return models.Tutor{}, err
	}

	return tutor, nil
}

func (tutorService *Tutor) DeleteByID(id int) error {
	query := `DELETE FROM tutors WHERE id = $1`

	result, err := tutorService.db.Exec(query, id) // Exec длля операций не возвращающих данные
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected() // возвращение кол-ва удаленных строк, проверяет, что запись существовала. На клиенте не нужно. Для 404
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tutor with id %d not found", id)
	}

	return nil
}

func (tutorService *Tutor) PostString(fullName string, email string) (int, error) {
	query := `insert into tutors (full_name, email) values
			($1,$2) RETURNING id`

	var id int

	row := tutorService.db.QueryRow(query, fullName, email)

	err := row.Scan(&id) // заполним id
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (tutorService *Tutor) PutString(fullName string, email string, id int) (models.Tutor, error) {
	query := `update tutors 
			set 
			full_name = $1, email = $2
			where id = $3
			RETURNING id, full_name, email`

	var tutor models.Tutor
	err := tutorService.db.QueryRow(query, fullName, email, id).Scan(
		&tutor.ID, &tutor.FullName, &tutor.Email,
	)

	if err != nil {
		return models.Tutor{}, err
	}

	return tutor, nil
}
