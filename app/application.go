package app

import (
	"github.com/pestanko/isstat/core"
	"github.com/pestanko/isstat/parsers"
	log "github.com/sirupsen/logrus"
	"os"
)

// IsStatApp - Is MUNI Statistics application
type IsStatApp struct {
	Client  core.CourseClient
	Parser  parsers.NotepadContentParser
	Results core.Results
	DryRun  bool
}

// Fetch - fetches the notepads content
func (app *IsStatApp) Fetch(notepads []string) ([]core.ResultItem, error) {
	timestamp := core.GetCurrentTimestamp()
	return app.FetchWithTimestamp(notepads, timestamp)
}

// FetchWithTimestamp - fetches the notepads content
func (app *IsStatApp) FetchWithTimestamp(notepads []string, timestamp string) ([]core.ResultItem, error) {
	var items []core.ResultItem

	for i, notepad := range notepads {
		log.WithField("index", i).WithField("name", notepads).Info("Fetching notepad")
		data, err := app.Client.GetNotepadContentData(notepad)

		if err != nil {
			log.WithField("name", notepad).WithError(err).Error("Unable to fetch the data")
			return items, err
		}

		resultItem := core.NewResultItem(notepad, timestamp, "xml")
		resultItem.Data = data

		items = append(items, resultItem)

		if err := app.Results.Store(&resultItem); err != nil {
			log.WithError(err).WithField("notepad", notepad).WithField("timestamp", timestamp).Error("Unable to store result")
			return items, err
		}
	}
	return items, nil
}

// GetApplication - gets an application instance
func GetApplication(config *Config) (IsStatApp, error) {
	client := core.NewCourseClient(config.Muni.URL, config.Muni.Token, config.Muni.Faculty, config.Muni.Course)
	client.DryRun = config.DryRun

	register := parsers.GetParserRegister()
	register.Register("default", &parsers.KontrFunctionalityParser{})
	parser := register.GetOrDefault(config.Parser)

	return IsStatApp{Client: client, Parser: parser, Results: core.NewResults(config.Results), DryRun: config.DryRun}, nil
}

func SetupLogger(loggingLevel string) {
	if loggingLevel == "" {
		loggingLevel = os.Getenv("LOG_LEVEL")
		if loggingLevel == "" {
			loggingLevel = "warning"
		}
	}

	level, err := log.ParseLevel(loggingLevel)
	if err != nil {
		log.WithError(err).WithField("level", loggingLevel).Warning("Unable to parse the log level")
		level = log.WarnLevel
	}

	log.SetLevel(level)
	log.SetOutput(os.Stderr)
}
