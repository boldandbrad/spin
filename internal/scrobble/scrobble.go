package scrobble

import (
	"fmt"
	"time"
)

type TimeMode int

const (
	TimeModeStartNow TimeMode = iota // 0
	TimeModeEndNow                   // 1
	TimeModeCustom                   // 2
)

func ResolveTimestampFromMode(mode TimeMode, customTime, customDate string, totalDurationSec int) (time.Time, error) {
	now := time.Now()

	switch mode {
	case TimeModeStartNow:
		return now, nil
	case TimeModeEndNow:
		return now.Add(-time.Duration(totalDurationSec) * time.Second), nil
	case TimeModeCustom:
		timestamp, err := parseTimeOfDay(customTime)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid time format '%s', use HH:MM (e.g., 14:30)", customTime)
		}
		if customDate != "" {
			parsedDate, err := time.Parse("2006-01-02", customDate)
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid date format '%s', use YYYY-MM-DD (e.g., 2026-04-18)", customDate)
			}
			timestamp = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second(), timestamp.Nanosecond(), timestamp.Location())
		}
		return timestamp, nil
	default:
		return now, nil
	}
}

func FormatTimestamp(t time.Time) string {
	return fmt.Sprintf("%d", t.Unix())
}

func parseTimeOfDay(timeStr string) (time.Time, error) {
	today := time.Now().Format("2006-01-02")
	combined := fmt.Sprintf("%s %s", today, timeStr)
	return time.Parse("2006-01-02 15:04", combined)
}
