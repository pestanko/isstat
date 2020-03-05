package parsers

import (
	"testing"
	"time"

	"github.com/pestanko/isstat/core"
)

/*
 === PARSE LINE TESTS
*/


func TestParseLine_4Elements(t *testing.T) {
	// GIVEN
	input := " 1  2020-02-18  08:45    *1 "
	var expected = core.Submission{
		Index: 1,
		DateTime: time.Date(2020, 02, 18, 8, 45, 0, 0, time.Local),
		Points: 1,
		Final: true,
	}

	// WHEN
	submission, err := parseSubmissionLine(input)

	// THEN
	if err != nil {
		t.Errorf("FAIL: Found error: %v", err)
	}

	assertSubmission(t, &expected, &submission)
}

func TestParseLine_4ElementsNonFinal(t *testing.T) {
	// GIVEN
	input := " 3  2020-02-18  08:45    10 "
	var expected = core.Submission{
		Index: 3,
		DateTime: time.Date(2020, 02, 18, 8, 45, 0, 0, time.Local),
		Points: 10,
		Final: false,
	}

	// WHEN
	submission, err := parseSubmissionLine(input)

	// THEN
	if err != nil {
		t.Errorf("FAIL: Found error: %v", err)
	}

	assertSubmission(t, &expected, &submission)
}

/*
 === PARSE NUMBER TESTS
*/
func TestParseNumber_OneFinal(t *testing.T) {
	//GIVEN
	input := "*1"

	var points int = 0
	var isFinal bool = false
	var err error

	// WHEN
	points, isFinal, err = parseNumberWithStar(input)

	// THEN
	if err != nil {
		t.Errorf("FAIL: Found error: %v", err)
	}

	if !isFinal {
		t.Error("FAIL: Should be final")
	}

	if points != 1 {
		t.Errorf("FAIL: points should 1, parsed: %d", points)
	}
}

func TestParseNumber_OneNormal(t *testing.T) {
	//GIVEN
	input := "1"

	var points int = 0
	var isFinal bool = false
	var err error

	// WHEN
	points, isFinal, err = parseNumberWithStar(input)

	// THEN
	if err != nil {
		t.Errorf("FAIL: Found error: %v", err)
	}

	if isFinal {
		t.Error("FAIL: Should not be final")
	}

	if points != 1 {
		t.Errorf("FAIL: points should be 1, parsed: %d", points)
	}
}

func TestParseNumber_NegativeNormal(t *testing.T) {
	//GIVEN
	input := "-1"

	var points int = 0
	var isFinal bool = false
	var err error

	// WHEN
	points, isFinal, err = parseNumberWithStar(input)

	// THEN
	if err != nil {
		t.Errorf("FAIL: Found error: %v", err)
	}

	if isFinal {
		t.Error("FAIL: Should not be final")
	}

	if points != -1 {
		t.Errorf("FAIL: points should be -1, parsed: %d", points)
	}
}

func TestParseNumber_NegativeFinal(t *testing.T) {
	//GIVEN
	input := "*-10"

	var points int = 0
	var isFinal bool = false
	var err error

	// WHEN
	points, isFinal, err = parseNumberWithStar(input)

	// THEN
	if err != nil {
		t.Errorf("FAIL: Found error: %v", err)
	}

	if !isFinal {
		t.Error("FAIL: Should be final")
	}

	if points != -10 {
		t.Errorf("FAIL: points should be -10, parsed: %d", points)
	}
}

func assertSubmission(t *testing.T, expected *core.Submission, provided *core.Submission) {
	if provided.Index != expected.Index {
		t.Errorf("FAIL: Submission index is %d, expected: %d", provided.Index, expected.Index)
	}

	if provided.Points != expected.Points {
		t.Errorf("FAIL: Submission Points is %d, expected: %d", provided.Points, expected.Points)
	}

	if provided.Final != expected.Final {
		t.Errorf("FAIL: Submission final is %v, expected: %v", provided.Final, expected.Final)
	}

	if provided.Bonus != provided.Bonus {
		t.Errorf("FAIL: Submission bonus is %v, expected: %v", provided.Bonus, expected.Bonus)
	}

	if provided.DateTime.Year() != expected.DateTime.Year() {
		t.Errorf("FAIL: Submission year is %d, expected: %d", provided.DateTime.Year(), expected.DateTime.Year())
	}

	if provided.DateTime.Month() != expected.DateTime.Month() {
		t.Errorf("FAIL: Submission month is %d, expected: %d", provided.DateTime.Month(), expected.DateTime.Month())
	}

	if provided.DateTime.Day() != expected.DateTime.Day() {
		t.Errorf("FAIL: Submission day is %d, expected: %d", provided.DateTime.Day(), expected.DateTime.Day())
	}

	if provided.DateTime.Hour() != expected.DateTime.Hour() {
		t.Errorf("FAIL: Submission hour is %d, expected: %d", provided.DateTime.Hour(), expected.DateTime.Hour())
	}

	if provided.DateTime.Minute() != expected.DateTime.Minute() {
		t.Errorf("FAIL: Submission minute is %d, expected: %d", provided.DateTime.Minute(), expected.DateTime.Minute())
	}
}