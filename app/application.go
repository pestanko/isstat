package app

import (
	"encoding/json"
	"github.com/pestanko/isstat/core"
	"github.com/pestanko/isstat/parsers"
	log "github.com/sirupsen/logrus"
	"os"
	"sort"
)

// IsStatApp - Is MUNI Statistics application
type IsStatApp struct {
	Client  core.CourseClient
	Parser  parsers.Parser
	Results core.Results
	Config  *Config
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
		if resultItem.Name == "" {
			continue
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

func (app *IsStatApp) ConvertToCSV(patterns []string) ([]core.ResultItem, error) {
	var items []core.ResultItem

	fileNames := app.Results.GlobAll(patterns)
	log.WithField("filenames", fileNames).Info("found filenames")

	for _, notepad := range fileNames {
		resultItem, err := app.ConvertToCSVOne(notepad)
		if err != nil {
			log.WithField("notepad", notepad).WithError(err).Error("Unable to convert to CSV")
			continue
		}
		if resultItem.Name == "" {
			continue
		}
		items = append(items, resultItem)
	}
	return items, nil
}

func (app *IsStatApp) ConvertToCSVOne(notepad string) (core.ResultItem, error) {
	log.WithField("name", notepad).Info("Converting to csv")

	jsonItem := core.NewResultItemFromFullName(notepad)

	if jsonItem.Ext != "json" {
		return core.ResultItem{}, nil
	}

	csvContent, err := app.convertStudentInfo(&jsonItem)
	if err != nil {
		return jsonItem, err
	}

	csvItem := core.NewResultItem(jsonItem.Name, jsonItem.TimeStamp, "csv")

	if err := core.WriteStatisticsToCSVFile(app.Results.GetPath(&csvItem), csvContent); err != nil {
		return csvItem, err
	}

	return csvItem, err
}

func (app *IsStatApp) Parse(patterns []string) (map[string][]core.StudentInfo, error) {
	var items = make(map[string][]core.StudentInfo)
	log.WithField("patterns", patterns).Info("Parse notepads")

	fileNames := app.Results.GlobAll(patterns)
	log.WithField("filenames", fileNames).Info("found filenames")

	for _, notepad := range fileNames {
		info, err := app.ParseOne(notepad)
		if err != nil {
			log.WithError(err).WithField("notepad", notepad).Error("Error in parsing the notepad")
			continue
		}

		items[notepad] = info
	}
	return items, nil
}

func (app *IsStatApp) ParseOne(notepad string) ([]core.StudentInfo, error) {
	log.WithField("name", notepad).Info("Parsing notepad")

	resultItem := core.NewResultItemFromFullName(notepad)

	if resultItem.Ext != "xml" {
		return []core.StudentInfo{}, nil
	}

	info, err := app.parseResultItem(&resultItem)
	if err != nil {
		return info, err
	}

	jsonitem := core.NewResultItem(resultItem.Name, resultItem.TimeStamp, "json")

	data, err := json.Marshal(info)
	if err != nil {
		log.WithError(err).WithField("notepad", notepad).Error("Unable to marshall json with data")
		return info, err
	}

	jsonitem.Data = data

	if err := app.Results.Store(&jsonitem); err != nil {
		log.WithError(err).WithField("notepad", notepad).WithField("timestamp", jsonitem.TimeStamp).Error("Unable to store result")
		return info, err
	}
	return info, nil
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

func (app *IsStatApp) convertStudentInfo(item *core.ResultItem) ([]core.CSVStatistic, error) {
	fileContent, err := app.Results.GetContent(item)
	if err != nil {
		return []core.CSVStatistic{}, err
	}

	infoContent, err := core.UnmarshalStudentInfo(fileContent)
	if err != nil {
		return []core.CSVStatistic{}, err
	}

	return core.ConvertSubmissionsToCSVStatistics(infoContent), nil
}

func (app *IsStatApp) CleanResults(patterns []string, limit int) ([]core.ResultItem, error) {
	items := app.PatternsToResultItems(patterns)

	var removedItems []core.ResultItem

	ItemsSortByTimestamp(items)

	for i, item := range items {
		if i > limit {
			if err := os.Remove(app.Results.GetPath(&item)); err != nil {
				log.WithField("fullname", item.GetFullName()).WithError(err).Error("Unable to remove")
				continue
			}
			removedItems = append(removedItems, item)
		}
	}

	return removedItems, nil
}

func (app *IsStatApp) PatternsToResultItems(patterns []string) []core.ResultItem {
	fileNames := app.Results.GlobAll(patterns)
	log.WithField("filenames", fileNames).Info("found filenames")
	var items []core.ResultItem

	for _, fileName := range fileNames {
		item := core.NewResultItemFromFullName(fileName)
		if item.TimeStamp == "" {
			continue
		}
		items = append(items, item)
	}

	return items
}

func (app *IsStatApp) GetLatest() (result map[string]map[string]core.ResultItem) {
	result = make(map[string]map[string]core.ResultItem)
	resultItems := app.PatternsToResultItems([]string{"*"})
	ItemsSortByTimestamp(resultItems)
	categories := CategorizeResultItems(resultItems)

	if len(categories) == 0 {
		return result
	}

	for name, extensions := range categories {
		result[name] = make(map[string]core.ResultItem)
		for ext, values := range extensions {
			result[name][ext] = values[0]
		}
	}

	return result
}

func (app *IsStatApp) DumpLatest() (result []string, err error) {
	latest := app.GetLatest()
	if len(latest) == 0 {
		return result, nil
	}

	for _, extensions := range latest {
		for _, item := range extensions {
			data, err := app.Results.GetContent(&item)
			if err != nil {
				log.WithField("name", item.GetFullName()).Error("Unable to fetch data")
			}
			item.Data = data
			err = app.Results.StoreWithoutTimestamp(&item)
			result = append(result, item.GetFullName())
		}
	}
	return result, err
}

func CategorizeResultItems(items []core.ResultItem) (result map[string]map[string][]core.ResultItem) {
	result = make(map[string]map[string][]core.ResultItem)
	named := CategorizeByName(items)
	for key, value := range named {
		result[key] = CategorizeByExtension(value)
	}

	return result
}

func CategorizeByExtension(items []core.ResultItem) map[string][]core.ResultItem {
	var result = make(map[string][]core.ResultItem)

	for _, item := range items {
		result[item.Ext] = append(result[item.Ext], item)
	}

	return result
}

func CategorizeByName(items []core.ResultItem) (result map[string][]core.ResultItem) {
	result = make(map[string][]core.ResultItem)

	for _, item := range items {
		result[item.Name] = append(result[item.Name], item)
	}

	return result
}

// GetApplication - gets an application instance
func GetApplication(config *Config) (IsStatApp, error) {
	client := core.NewCourseClient(config.Muni.URL, config.Muni.Token, config.Muni.Faculty, config.Muni.Course)
	client.DryRun = config.DryRun

	register := parsers.GetParserRegister()
	register.Register("default", &parsers.KontrFunctionalityParser{})
	parser := register.GetOrDefault(config.Parser)
	basicParser := parsers.BasicParser{
		StudentsRegister:     core.NewStudentsRegister(),
		NotepadContentParser: parser,
	}

	return IsStatApp{Client: client, Parser: &basicParser, Results: core.NewResults(config.Results, config.WithoutTimestamp), Config: config}, nil
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

func ItemsSortByTimestamp(items []core.ResultItem) []core.ResultItem {
	sort.Slice(items[:], func(i, j int) bool {
		return items[i].TimeStamp > items[j].TimeStamp
	})
	return items
}
