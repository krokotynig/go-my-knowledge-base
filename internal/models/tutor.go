package models

type Tutor struct {
	ID       int    `db:"id" json:"id"`
	Fullname string `db:"full_name" json:"full_name"`
	Email    string `db:"email" json:"email"`
}
