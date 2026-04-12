package scrobble

import (
	"fmt"
	"strconv"
	"time"
)

func FormatTimestamp(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

func ParseTimestamp(ts string) (time.Time, error) {
	unix, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timestamp: %w", err)
	}
	return time.Unix(unix, 0), nil
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
