package parsers

import (
	"github.com/pestanko/isstat/core"
	log "github.com/sirupsen/logrus"

)

// NotepadContentParser - public parser interface
type NotepadContentParser interface {
	Parse(content string) ([]core.Submission, error)
}

// ParseNotepadContent - parses notepad content
func ParseNotepadContent(parser NotepadContentParser, content string) ([]core.Submission, error) {
	return parser.Parse(content)
}

// Parser - The main parser
type Parser interface {
	Parse(content *core.NotepadContent) ([]core.StudentInfo, error)
}

// BasicParser implementation
type BasicParser struct {
	StudentsRegister core.StudentsRegister
	NotepadContentParser NotepadContentParser
}

// Parse the provided student's content using the Basic parser
func (parser *BasicParser) Parse(content *core.NotepadContent) ([]core.StudentInfo, error)  {
	var students = make([]core.StudentInfo, len(content.StudentsContent))

	for i, student := range content.StudentsContent {
		var uid = parser.StudentsRegister.GetOrRegister(student.Uco)
		var err error

		students[i] = core.NewStudentSubmissions(uid)

		log.WithField("index", i).WithField("student_uco", student.Uco).WithField("content", student.Content).Debug("parsing content")
		students[i].Submissions, err = parser.NotepadContentParser.Parse(student.Content)

		if err != nil {
			log.WithField("content", student.Content).Error("Unable to parse submissions")
			continue
		}
	}

	return students, nil
}

