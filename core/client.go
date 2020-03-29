package core

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// NotepadContent - whole notepad content
type NotepadContent struct {
	StudentsContent []StudentContent `xml:"STUDENT" json:"students"`
}

// StudentContent - represents one entity in the notepad
type StudentContent struct {
	Content   string `xml:"OBSAH" json:"content"`
	Uco       string `xml:"UCO" json:"uco"`
	ChangedBy string `xml:"ZMENIL" json:"changed_by"`
}

// CourseClient - crawls the is muni notepads
type CourseClient struct {
	URL       string
	Token     string
	FacultyID int
	Course    string
	DryRun    bool
}

// UnmarshalNotepadContent - unmarshal the notepad content
func UnmarshalNotepadContent(data []byte) (content NotepadContent, err error) {
	if err := xml.Unmarshal(data, &content); err != nil {
		return content, err
	}

	return content, nil
}

// NewCourseClient - Creates a new course client
func NewCourseClient(url string, token string, facultyID int, course string) CourseClient {
	return CourseClient{URL: url, Token: token, FacultyID: facultyID, Course: course, DryRun: false}
}

//GetNotepadContent - Gets a notepad content
func (client *CourseClient) GetNotepadContent(notepadCodename string) (NotepadContent, error) {
	data, err := client.GetNotepadContentData(notepadCodename)
	if err != nil {
		return NotepadContent{}, err
	}

	if client.DryRun {
		return NotepadContent{}, nil
	}

	return UnmarshalNotepadContent(data)
}

// GetNotepadContentData - Gets a raw notepad content data
func (client *CourseClient) GetNotepadContentData(notepadCodename string) ([]byte, error) {
	notepadURL := client.buildNotesURL(notepadCodename)

	log.WithField("url", notepadURL).Info("Using the notepad url")

	data, err := client.Fetch(notepadURL)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//Save the data to the XML file (caching)
func (client *CourseClient) Save(data []byte, name string) error {

	if client.DryRun {
		return nil
	}

	return ioutil.WriteFile(name, data, 0644)
}

// Fetch - fetches XML data
func (client *CourseClient) Fetch(url string) ([]byte, error) {
	log.WithField("url", url).Debug("Fetching data")

	if client.DryRun {
		return []byte{}, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		log.WithError(err).WithField("url", url).Error("Fetch failed")
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Status error: %d", resp.StatusCode)
		log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}
	return data, nil
}

func (client *CourseClient) buildNotesURL(notepadCodename string) string {
	return fmt.Sprintf(
		"%s/export/pb_blok_api?klic=%s;fakulta=%d;kod=%s;operace=blok-dej-obsah;zkratka=%s",
		client.URL, client.Token, client.FacultyID, client.Course, notepadCodename)
}
