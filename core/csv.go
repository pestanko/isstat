package core

import (
	"os"

	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

// CSVStatistic - representation of the CSV statistics
type CSVStatistic struct {
	StudentID string  `csv:"student_id"`
	Index     int     `json:"index" csv:"index"`
	DateTime  string  `csv:"datetime"`
	Date      string  `csv:"date"`
	Time      string  `csv:"time"`
	Points    float64 `json:"points" csv:"points"`
	Bonus     float64 `json:"bonus" csv:"bonus"`
	Final     bool    `json:"final" csv:"final"`
}

// WriteStatisticsToCSVFile - writes statistics to the CSV file
func WriteStatisticsToCSVFile(file string, statistics []CSVStatistic) error {
	csvFile, err := os.Create(file)
	if err != nil {
		log.WithField("file", file).WithError(err).Error("Unable to create file")
		return err
	}

	defer csvFile.Close()

	return gocsv.MarshalFile(statistics, csvFile)
}

// ConvertSubmissionsToCSVStatistics - Converter
func ConvertSubmissionsToCSVStatistics(students []StudentInfo) []CSVStatistic {
	var stats []CSVStatistic

	for _, student := range students {
		for _, submission := range student.Submissions {
			stat := CSVStatistic{
				StudentID: student.ID.String(),
				Index:     submission.Index,
				DateTime:  submission.DateTime.Format("2006-01-02T15:04"),
				Date:      submission.DateTime.Format("2006-01-02"),
				Time:      submission.DateTime.Format("15-04"),
				Points:    submission.Points,
				Final:     submission.Final,
				Bonus:     submission.Bonus,
			}
			stats = append(stats, stat)
		}
	}

	return stats
}
