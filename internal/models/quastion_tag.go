package models

type QuestionTag struct {
	QuestionID int `db:"question_id" json:"question_id"`
	TagID      int `db:"tag_id" json:"tag_id"`
}
