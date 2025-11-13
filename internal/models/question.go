package models

import "time"

type Question struct {
	ID           int       `db:"id" json:"id"`
	QuestionText string    `db:"question_text" json:"question_text"`
	TutorID      *int      `db:"tutor_id" json:"tutor_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	IsEdit       bool      `db:"is_edit" json:"is_edit"`
}
