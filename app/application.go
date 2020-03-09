package app

import (
	"github.com/pestanko/isstat/core"
	"github.com/pestanko/isstat/parsers"
	log "github.com/sirupsen/logrus"
)

// IsStatApp - Is MUNI Statistics application
type IsStatApp struct {
	Client core.CourseClient
	Parser parsers.NotepadContentParser
	Cache core.Cache
}

func (app *IsStatApp) Fetch(notepads []string) error {
	for i, notepad := range notepads {
		log.WithField("index", i).WithField("name", notepads).Info("Fetching notepad")
		data, err := app.Client.GetNotepadContentData(notepad)

		if err != nil {
			log.WithField("name", notepad).WithError(err).Error("Unable to fetch the data")
			return err
		}

		serialized, err := core.UnmarshalNotepadContent(data)
		if(err != nil) {

		}
	}
}

// NewBasicApp - create a new basic application
func NewBasicApp(client core.CourseClient, parser parsers.NotepadContentParser, ) IsStatApp {
	return IsStatApp {Client: client, Parser: parser}
}

// GetApplication - gets an application instance
func GetApplication(config *Config) (IsStatApp, error) {
	client := core.NewCourseClient(config.IsMuni.URL, config.IsMuni.Token, config.IsMuni.FacultyID, config.IsMuni.Course)
	parser, err := parsers.GetParserRegister().Get(config.Parser)
	if err != nil {
		return IsStatApp{}, err
	}

	return IsStatApp{ Client: client, Parser: parser, Cache: core.NewCache(config.Cache.Directory)}, nil
}
