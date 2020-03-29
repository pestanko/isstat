package core

import (
	"os"

	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

// CSVStatistic - representation of the CSV statistics
type CSVStatistic struct {
	StudentID string  `csv:"student_id"`
	DateTime  string  `csv:"datetime"`
	Index     int     `json:"index" csv:"index"`
	Points    float64 `json:"points" csv:"points"`
	Final     bool    `json:"final" csv:"final"`
	Bonus     float64 `json:"bonus" csv:"bonus"`
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
				DateTime:  submission.DateTime.String(),
				Index:     submission.Index,
				Points:    submission.Points,
				Final:     submission.Final,
				Bonus:     submission.Bonus,
			}
			stats = append(stats, stat)
		}
	}
	return stats
}
