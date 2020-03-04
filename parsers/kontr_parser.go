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
			submission, err := parseSubmission(line)
			if err != nil {
				log.Errorf("Unable to parse the submission: %v", err)
			}
			submissions = append(submissions, *submission)
		}
	}
	return submissions, nil
}

/**
%%       datum    cas  body
 1  2020-02-18  08:45    *1
 */
func parseSubmission(line string) (*core.Submission, error) {
	submission := &core.Submission{ Bonus:0, Final: false}
	words := strings.Fields(line)

	wordsCount := len(words)

	if wordsCount == 0 {
		return nil, fmt.Errorf("not enought line parts - %d found", wordsCount)
	}

	index, err := strconv.Atoi(words[0])
	if err != nil {
		return nil, err
	}
	submission.Index = index

	if wordsCount <= 2 {
		return submission, fmt.Errorf("not enought line parts - %d found - unable to parse datetime", wordsCount)
	}

	datetime, err := parseDateTime(words[1], words[2])
	if err != nil {
		return submission, err
	}
	submission.DateTime = datetime

	if wordsCount <= 3 {
		return submission, fmt.Errorf("not enought line parts - %d found - unable to parse points", wordsCount)
	}

	points, isFinal, err := parseNumberWithStar(words[3])
	if err != nil {
		return submission, err
	}

	submission.Points = points
	submission.Final = isFinal

	if wordsCount == 4 {
		return submission, nil
	}

	bonus, err := strconv.Atoi(words[4])
	if err != nil {
		return submission, err
	}
	submission.Bonus = bonus

	if wordsCount == 5 {
		return submission, nil
	}

	_, isFinal, err = parseNumberWithStar(words[5])

	if err != nil {
		return submission, err
	}

	if !submission.Final && isFinal {
		submission.Final = true
	}

	return submission, nil
}

func parseDateTime(date string, tPart string) (time.Time, error) {
	full := date + " " + tPart

	return time.Parse("2006-01-02 15:04", full)
}

func parseNumberWithStar(s string) (int, bool, error) {
	var isFinal = false
	if strings.HasPrefix(s, "*") {
		isFinal = true
		s = s[1:]
	}

	points, err := strconv.Atoi(s)

	if err != nil {
		return 0, isFinal, err
	}

	return points, isFinal, nil
}
