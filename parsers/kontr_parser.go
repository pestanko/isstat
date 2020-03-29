package parsers

import (
	"fmt"
	"github.com/pestanko/isstat/core"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

// KontrFunctionalityParser - parses functionality points
type KontrFunctionalityParser struct {
}

/*
Parse the notepad points

Format:
-----
# zapsáno z Kontru 2020-02-18 08:45, v2.2.1

%%       datum    cas  body
 1  2020-02-18  08:45    *1

# POZOR: Tento blok NEUPRAVUJTE!

# Kontr může veškeré změny kdykoliv přepsat.
# Poznámky k odevzdání a hodnocení čistoty pište
# do bloku určeného pro tyto účely.
*/
func (parser *KontrFunctionalityParser) Parse(content string) ([]core.Submission, error) {
	lines := strings.Split(content, "\n")

	var foundHeader = false
	var submissions []core.Submission

	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "%%") {
			foundHeader = true
			continue
		}
		if foundHeader {
			submission, err := parseSubmissionLine(line)
			if err != nil {
				log.Errorf("Unable to parse the submission: %v", err)
			}
			submissions = append(submissions, submission)
		}
	}
	return submissions, nil
}

/*
ParseSubmissionLine - Parses one sumission line

%%       datum    cas  body
 1  2020-02-18  08:45    *1
 */
func parseSubmissionLine(line string) (core.Submission, error) {
	submission := core.Submission{ Bonus:0, Final: false}
	fields := strings.Fields(line)

	wordsCount := len(fields)

	if wordsCount == 0 {
		return submission, fmt.Errorf("not enought line parts - %d found", wordsCount)
	}

	var err error

	submission.Index, err = strconv.Atoi(fields[0])
	if err != nil {
		return submission, err
	}

	if wordsCount == 2 {
		log.Infof("Found just %d parts", wordsCount)
		return submission, nil
	}

	submission.DateTime, err = time.Parse("2006-01-02 15:04", fields[1] + " " + fields[2])
	if err != nil {
		return submission, err
	}

	if wordsCount <= 3 {
		log.Infof("Found just %d parts: %v", wordsCount, submission)
		return submission, nil
	}

	submission.Points, submission.Final, err = parseNumberWithStar(fields[3])
	if err != nil {
		return submission, err
	}

	if wordsCount == 4 {
		return submission, nil
	}

	submission.Bonus , err = strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return submission, err
	}

	if wordsCount == 5 {
		return submission, nil
	}

	_, isFinal, err := parseNumberWithStar(fields[5])

	if err != nil {
		return submission, err
	}

	if !submission.Final && isFinal {
		submission.Final = true
	}

	return submission, nil
}


func parseNumberWithStar(s string) (float64, bool, error) {
	var isFinal = false
	if strings.HasPrefix(s, "*") {
		isFinal = true
		s = s[1:]
	}

	points, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return 0, isFinal, err
	}

	return points, isFinal, nil
}
