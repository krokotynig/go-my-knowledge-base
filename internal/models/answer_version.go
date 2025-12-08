package models

import "time"

type AnswerVersion struct {
	ID            int       `db:"id" json:"id"`
	AnswerID      int       `db:"answer_id" json:"answer_id"`
	AnswerText    string    `db:"answer_text" json:"answer_text"`
	QuestionID    int       `db:"question_id" json:"question_id"`
	TutorID       *int      `db:"tutor_id" json:"tutor_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	VersionNumber int       `db:"version_number" json:"version_number"`
	IsDelete      bool      `db:"is_delete" json:"is_delete"`
	DeleteByTutor *int      `db:"delete_by_tutor" json:"delete_by_tutor"`
}
