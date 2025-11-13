package models

type Tag struct {
	ID      int    `db:"id" json:"id"`
	TutorID *int   `db:"tutor_id" json:"tutor_id"`
	Tag     string `db:"tag" json:"tag"`
}
