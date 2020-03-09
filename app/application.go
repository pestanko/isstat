package app

import (
	"github.com/pestanko/isstat/core"
	"github.com/pestanko/isstat/parsers"
)

// IsStatApp - Is MUNI Statistics application
type IsStatApp struct {
	Client *core.CourseClient
	Parser parsers.Parser
}

// NewBasicApp - create a new basic application
func NewBasicApp(client *core.CourseClient, parser parsers.NotepadContentParser) IsStatApp {
	return IsStatApp {Client: client, Parser: parser}
}

// GetApplication - gets an application instance
func GetApplication(config *Config) (IsStatApp, error) {
	client := core.NewCourseClient(config.IsMuni.URL, config.IsMuni.Token, config.IsMuni.FacultyID, config.IsMuni.Course)
	parser, err := parsers.GetParserRegister().Get(config.Parser)
	if err != nil {
		return IsStatApp{}, err
	}

	return NewBasicApp(client, parser), nil
}
