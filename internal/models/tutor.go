package models

type Tutor struct {
	ID       int    `db:"id" json:"id"`
	FullName string `db:"full_name" json:"full_name"`
	Email    string `db:"email" json:"email"`
}

type TutorSwaggerRequestBody struct { // Нужен для запросов в swagger
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}
