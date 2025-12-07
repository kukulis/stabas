package my_tests

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	timeStr := "2025-12-05T15:04:05Z"
	//timeStr := "2006-01-02T15:04:05Z0300

	parsedTime, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		t.Error("Error parsing time " + err.Error())
	}

	got := parsedTime.Format(time.RFC3339)
	if got != timeStr {
		t.Errorf("Parsed time %s does not match expected time %s", got, timeStr)
	}
}

// TODO test parse in location

//func TestParseCarbon(t *testing.T) {
//	time2Str := "2025-12-05T15:04:05Z03:00"
//
//	time2 := carbon.ParseByFormat(time2Str, carbon.RFC3339Format)
//
//	got := time2.Format(carbon.RFC3339Format)
//	if got != time2Str {
//		t.Errorf("Parsed time %s does not match expected time %s", got, time2Str)
//	}
//}
