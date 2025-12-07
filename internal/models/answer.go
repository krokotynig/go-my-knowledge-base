package models

import "time"

type Answer struct {
	ID          int       `db:"id" json:"id"`
	AnswersText string    `db:"answer_text" json:"answer_text"`
	TutorID     *int      `db:"tutor_id" json:"tutor_id"`
	QuestionID  int       `db:"question_id" json:"question_id"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	IsEdit      bool      `db:"is_edit" json:"is_edit"`
}

// Модель для swagger записи POST
type AnswersSwaggerRequestPostBody struct {
	AnswersText string `json:"answer_text"`
	TutorID     *int   `json:"tutor_id,omitempty"`
	QuestionID  int    `json:"question_id"`
}

// Модель для swagger записи PUT
type AnswersSwaggerRequestPutBody struct {
	AnswersText string `json:"answer_text"`
	TutorID     *int   `json:"tutor_id,omitempty"`
	QuestionID  int    `json:"question_id"`
	IsEdit      bool   `json:"is_edit"`
}
