package core

import uuid "github.com/google/uuid"

import "time"


// StudentInfo - representation of the student submissions
type StudentInfo struct {
	ID uuid.UUID `json:"uid"`
	Submissions []Submission `json:"submissions"`
}

// Submission - representation of the one student submission
type Submission struct {
	DateTime time.Time `json:"datetime"`
	Index int `json:"index"`
	Points float64 `json:"points"`
	Final bool `json:"final"`
	Bonus float64 `json:"bonus"`
}

// NewStundentSubmissions creates instance for the stundet submissions
func NewStundentSubmissions(uid uuid.UUID) StudentInfo {
	return StudentInfo{ID: uid, Submissions: []Submission{}}
}
