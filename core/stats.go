package core

import uuid "github.com/google/uuid"

import "time"


// StudentSubmissions - representation of the student submissions
type StudentSubmissions struct {
	ID uuid.UUID `json:"uid"`
	Submissions []Submission `json:"submissions"`
}

// Submission - representation of the one student submission
type Submission struct {
	DateTime time.Time `json:"datetime"`
	Index int `json:"index"`
	Points int `json:"points"`
	Final bool `json:"final"`
	Bonus int `json:"bonus"`
}


