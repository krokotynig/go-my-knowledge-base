package models

import "time"

type Question struct {
	ID           int       `db:"id" json:"id"`
	QuestionText string    `db:"question_text" json:"question_text"`
	TutorID      *int      `db:"tutor_id" json:"tutor_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	IsEdit       bool      `db:"is_edit" json:"is_edit"`
}

type QuestionsSwaggerRequestBody struct {
	QuestionText string `json:"question_text"`
	TutorID      *int   `json:"tutor_id"`
	IsEdit       bool   `json:"is_edit"`
}
