package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	log "github.com/sirupsen/logrus"
)

// Results - stucture to hold the result
type Results struct {
	ResultsDir string
}

// ResultItem - represent one item in results
type ResultItem struct {
	Name      string `json:"name"`
	TimeStamp string `json:"timestamp"`
	Ext       string `json:"ext"`
	Data      []byte `json:"data"`
}

// NewResultItem - Creates a new result item
func NewResultItem(name, timestamp, ext string) ResultItem {
	return ResultItem{Name: name, TimeStamp: timestamp, Ext: ext}
}

func NewResultItemFromFullName(fullName string) ResultItem {
	item := ResultItem{}

	_, err := fmt.Sscanf(fullName, "%s.%s.%s", &item.Name, &item.TimeStamp, &item.Ext)
	if err != nil {
		log.WithField("fullname", fullName).Warning("Unable to parse the result item from the provided full name")
	}

	return item
}

// GetFullName for the item
func (item *ResultItem) GetFullName() string {
	return fmt.Sprintf("%s.%s.%s", item.Name, item.TimeStamp, item.Ext)
}

func (item *ResultItem) getLogEntry() *log.Entry {
	return log.WithField("name", item.Name).
		WithField("timestamp", item.TimeStamp).
		WithField("ext", item.Ext).
		WithField("fullname", item.GetFullName)
}

// NewResults - Creates a new result holder
func NewResults(resultsDir string) Results {
	var err error
	if resultsDir == "" {
		resultsDir, err = os.Getwd()
		if err != nil {
			log.WithError(err).Warning("Unable to get current working directory")
		}
	}

	log.WithField("dir", resultsDir).Info("Results dir location")
	return Results{ResultsDir: resultsDir}
}

// Store - store the content to the file
func (results *Results) Store(item *ResultItem) error {
	if item.TimeStamp == "" {
		item.TimeStamp = GetCurrentTimestamp()
	}

	fullPath := path.Join(results.ResultsDir, item.GetFullName())
	item.getLogEntry().WithField("path", fullPath).Info("Storing result")

	return ioutil.WriteFile(fullPath, item.Data, 0644)
}

// Get - get item's content
func (results *Results) Get(item *ResultItem) (*ResultItem, error) {
	fullpath := path.Join(results.ResultsDir, item.GetFullName())

	entry := log.WithField("fullpath", fullpath)

	data, err := ioutil.ReadFile(fullpath)

	if err != nil {
		entry.WithError(err).Error("Unable to read a file")
		return nil, err
	}

	item.Data = data

	return item, nil
}

// List all Result entries
func (results *Results) List() (items []ResultItem, err error) {
	files, err := results.ListPaths()

	if err != nil {
		return nil, err
	}
	for _, fpath := range files {
		fname := path.Base(fpath)
		log.WithField("path", fpath).Debug("Processing path")
		var name, datetime, ext string
		if n, err := fmt.Sscanf(fname, "%s.%s.%s", name, datetime, ext); err != nil {
			log.WithError(err).
				WithField("n", n).
				WithField("filename", fname).
				Error("Unable to parse file name")
				return nil, err
		}
		items = append(items, NewResultItem(name, datetime, ext))
	}
	return items, nil
}

// ListPaths all path in the results dir
func (results *Results) ListPaths() (paths []string, err error) {
	files, err := ioutil.ReadDir(results.ResultsDir)

    if err != nil {
		return paths, err
	}

	for _, f := range files {

		fp := path.Join(results.ResultsDir, f.Name())
		paths = append(paths, fp)
	}
	return paths, nil
}

// GetPath - gets a full result path
func (results *Results) GetPath(item *ResultItem) string {
	return path.Join(results.ResultsDir, item.GetFullName())
}

// Gets content as bytes
func (results *Results) GetContent(item *ResultItem) ([]byte, error) {
	fp := path.Join(results.ResultsDir, item.GetFullName())
	return ioutil.ReadFile(fp)
}


// GetCurrentTimestamp - Gets a current timestamp
func GetCurrentTimestamp() string {
	return time.Now().Format("2006-01-02T15-04-05")
}
