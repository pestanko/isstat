package app

import (
	"encoding/json"
	"github.com/pestanko/isstat/core"
	"github.com/pestanko/isstat/parsers"
	log "github.com/sirupsen/logrus"
	"os"
)

// IsStatApp - Is MUNI Statistics application
type IsStatApp struct {
	Client  core.CourseClient
	Parser  parsers.Parser
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

	for _, notepad := range notepads {
		resultItem, err2 := app.FetchOne(notepad, timestamp)
		if err2 != nil {
			return items, err2
		}
		items = append(items, resultItem)
	}
	return items, nil
}

func (app *IsStatApp) FetchOne(notepad string, timestamp string) (core.ResultItem, error) {
	log.WithField("name", notepad).Info("Fetching notepad")
	data, err := app.Client.GetNotepadContentData(notepad)

	if err != nil {
		log.WithField("name", notepad).WithError(err).Error("Unable to fetch the data")
		return core.ResultItem{}, err
	}

	resultItem := core.NewResultItem(notepad, timestamp, "xml")
	resultItem.Data = data

	if err := app.Results.Store(&resultItem); err != nil {
		log.WithError(err).WithField("notepad", notepad).WithField("timestamp", timestamp).Error("Unable to store result")
		return core.ResultItem{}, err
	}
	return resultItem, nil
}

func (app *IsStatApp) Parse(notepads []string) (map[string][]core.StudentInfo, error) {
	var items map[string][]core.StudentInfo = make(map[string][]core.StudentInfo)

	for i, notepad := range notepads {
		log.WithField("index", i).WithField("name", notepads).Info("Parsing notepad")

		resultItem := core.NewResultItemFromFullName(notepad)

		info, err := app.parseResultItem(&resultItem)
		if err != nil {
			return items, err
		}

		jsonitem := core.NewResultItem(resultItem.Name, resultItem.TimeStamp, "json")

		data, err := json.Marshal(info)
		if err != nil {
			log.WithError(err).WithField("notepad", notepad).Error("Unable to marshall json with data")
			return items, err
		}

		jsonitem.Data = data

		if err := app.Results.Store(&jsonitem); err != nil {
			log.WithError(err).WithField("notepad", notepad).WithField("timestamp", jsonitem.TimeStamp).Error("Unable to store result")
			return items, err
		}

		items[notepad] = info
	}
	return items, nil
}

func (app *IsStatApp) parseResultItem(item *core.ResultItem) ([]core.StudentInfo, error) {
	fileContent, err := app.Results.GetContent(item)
	if err != nil {
		return []core.StudentInfo{}, err
	}

	notepadContent, err := core.UnmarshalNotepadContent(fileContent)
	if err != nil {
		return []core.StudentInfo{}, err
	}

	return app.Parser.Parse(&notepadContent)
}

// GetApplication - gets an application instance
func GetApplication(config *Config) (IsStatApp, error) {
	client := core.NewCourseClient(config.Muni.URL, config.Muni.Token, config.Muni.Faculty, config.Muni.Course)
	client.DryRun = config.DryRun

	register := parsers.GetParserRegister()
	register.Register("default", &parsers.KontrFunctionalityParser{})
	parser := register.GetOrDefault(config.Parser)
	basicParser := parsers.BasicParser{
		StudentsRegister: core.NewStudentsRegister(),
		NotepadContentParser: parser,
	}

	return IsStatApp{Client: client, Parser: &basicParser, Results: core.NewResults(config.Results), DryRun: config.DryRun}, nil
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
