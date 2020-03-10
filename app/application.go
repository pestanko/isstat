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
	Results core.Results
}

// Fetch - fetches the notepads content
func (app *IsStatApp) Fetch(notepads []string) error {
	timestamp := core.GetCurrentTimestamp()

	for i, notepad := range notepads {
		log.WithField("index", i).WithField("name", notepads).Info("Fetching notepad")
		data, err := app.Client.GetNotepadContentData(notepad)

		if err != nil {
			log.WithField("name", notepad).WithError(err).Error("Unable to fetch the data")
			return err
		}

		if err := app.Results.StoreWithTimestamp(notepad, timestamp, data); err != nil {
			log.WithError(err).WithField("notepad", notepad).WithField("timestamp", timestamp).Error("Unable to store result")
			return err
		}
	}
	return nil
}


// GetApplication - gets an application instance
func GetApplication(config *Config) (IsStatApp, error) {
	client := core.NewCourseClient(config.IsMuni.URL, config.IsMuni.Token, config.IsMuni.FacultyID, config.IsMuni.Course)
	parser, err := parsers.GetParserRegister().Get(config.Parser)
	if err != nil {
		return IsStatApp{}, err
	}

	return IsStatApp{ Client: client, Parser: parser, Results: core.NewResults(config.ResultsDir)}, nil
}
