package tui

import (
	"fmt"
	"time"

	"github.com/boldandbrad/spin/internal/api"
	"github.com/boldandbrad/spin/internal/profile"
	"github.com/boldandbrad/spin/internal/scrobble"
	"github.com/charmbracelet/huh"
)

func TrackTUI(profileFlag string, dryrun bool) error {
	artist := ""
	track := ""
	timestampMode := 1
	dateInput := ""
	timeInput := ""

	inputForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Artist").Value(&artist).Placeholder("e.g., Radiohead"),
			huh.NewInput().Title("Track").Value(&track).Placeholder("e.g., Paranoid Android"),
		),
		huh.NewGroup(
			huh.NewSelect[int]().
				Title("When did you listen?").
				Options(
					huh.NewOption("Starting now", 0),
					huh.NewOption("Ending now", 1),
					huh.NewOption("Custom start time", 2),
				).
				Value(&timestampMode),
		),
	)

	if err := inputForm.Run(); err != nil {
		return err
	}

	if artist == "" || track == "" {
		return nil
	}

	if timestampMode == 2 {
		timeForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Date (YYYY-MM-DD)").Value(&dateInput).Placeholder("e.g., 2026-04-12"),
				huh.NewInput().Title("Time (HH:MM)").Value(&timeInput).Placeholder("e.g., 15:00"),
			),
		)
		if err := timeForm.Run(); err != nil {
			return err
		}
	}

	client := api.NewClient()
	searchResults, err := client.SearchTrack(artist, track)
	if err != nil {
		return fmt.Errorf("failed to search tracks: %w", err)
	}

	if len(searchResults) == 0 {
		fmt.Println("No tracks found.")
		return nil
	}

	selectedTrack := searchResults[0]
	artistName := selectedTrack.Artist
	if artistName == "" {
		artistName = artist
	}

	var timestamp time.Time
	switch timestampMode {
	case 0:
		timestamp = time.Now()
	case 1:
		timestamp = time.Now()
		duration, err := client.GetTrackInfo(artistName, selectedTrack.Name)
		if err == nil && duration > 0 {
			timestamp = timestamp.Add(-time.Duration(duration) * time.Millisecond)
		}
	case 2:
		timestamp = time.Now()
		if timeInput != "" {
			t, err := scrobble.ParseTimeOfDay(timeInput)
			if err == nil {
				timestamp = t
			}
		}
		if dateInput != "" {
			t, err := scrobble.ParseDateTime(dateInput, timestamp.Format("15:04"))
			if err == nil {
				timestamp = t
			}
		}
	}

	ts := scrobble.FormatTimestamp(timestamp)
	tsFormatted := timestamp.Format("2006-01-02 15:04")

	username, err := profile.ResolveProfile(profileFlag)
	if err != nil {
		return err
	}

	if dryrun {
		fmt.Printf("Would scrobble to %s:\n", username)
	} else {
		fmt.Printf("Scrobbled to %s:\n", username)
	}
	fmt.Printf("  1. %s - %s (%s)\n", artistName, selectedTrack.Name, tsFormatted)

	if dryrun {
		return nil
	}

	cred, err := profile.GetCredentialForProfile(profileFlag)
	if err != nil {
		return err
	}

	if err := client.ScrobbleTrack(artistName, selectedTrack.Name, ts, cred.SessionKey); err != nil {
		return fmt.Errorf("failed to scrobble: %w", err)
	}

	return nil
}
