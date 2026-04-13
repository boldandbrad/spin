package cmd

import (
	"fmt"
	"time"
)

func ParseTimeArg(arg string) (time.Time, error) {
	t, err := time.Parse("2006-01-02 15:04", arg)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse("15:04", arg)
	if err == nil {
		today := time.Now().Format("2006-01-02")
		return time.Parse("2006-01-02 15:04", today+" "+arg)
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s (expected YYYY-MM-DD or HH:MM)", arg)
}
