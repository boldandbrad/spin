package scrobble

import (
	"fmt"
	"time"
)

func FormatTimestamp(t time.Time) string {
	return fmt.Sprintf("%d", t.Unix())
}

func ParseTimeOfDay(timeStr string) (time.Time, error) {
	today := time.Now().Format("2006-01-02")
	combined := fmt.Sprintf("%s %s", today, timeStr)
	return time.Parse("2006-01-02 15:04", combined)
}

func ParseDateTime(dateStr, timeStr string) (time.Time, error) {
	combined := fmt.Sprintf("%s %s", dateStr, timeStr)
	return time.Parse("2006-01-02 15:04", combined)
}

func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func ApplyDateToTimestamp(dateStr string, timestamp time.Time) (time.Time, error) {
	if dateStr == "" {
		return timestamp, nil
	}
	t, err := ParseDate(dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), t.Day(), timestamp.Hour(), timestamp.Minute(), 0, 0, timestamp.Location()), nil
}
