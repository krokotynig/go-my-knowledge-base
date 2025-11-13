package models

import "time"

type QuestionVersion struct {
	ID            int       `db:"id" json:"id"`
	QuestionID    int       `db:"question_id" json:"question_id"`
	QuestionText  string    `db:"question_text" json:"question_text"`
	TutorID       *int      `db:"tutor_id" json:"tutor_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	VersionNumber int       `db:"version_number" json:"version_number"`
}
