package tui

import (
	"fmt"
	"time"

	"github.com/boldandbrad/spin/internal/api"
	"github.com/boldandbrad/spin/internal/keyring"
	"github.com/boldandbrad/spin/internal/profile"
	"github.com/boldandbrad/spin/internal/scrobble"
	"github.com/charmbracelet/huh"
)

func TrackTUI(profileFlag string, dryrun bool) error {
	artist := ""
	track := ""
	dateInput := ""
	timeInput := ""

	inputForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Artist").Value(&artist).Placeholder("e.g., Radiohead"),
			huh.NewInput().Title("Track").Value(&track).Placeholder("e.g., Paranoid Android"),
		),
		huh.NewGroup(
			huh.NewInput().Title("Date (YYYY-MM-DD)").Value(&dateInput).Placeholder("optional"),
			huh.NewInput().Title("Time (HH:MM)").Value(&timeInput).Placeholder("optional"),
		),
	)

	if err := inputForm.Run(); err != nil {
		return err
	}

	if artist == "" || track == "" {
		return nil
	}

	timestamp := time.Now()
	if dateInput != "" || timeInput != "" {
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

	ts := scrobble.FormatTimestamp(timestamp)
	tsFormatted := timestamp.Format("2006-01-02 15:04")

	if dryrun {
		fmt.Printf("Would scrobble:\n")
	} else {
		fmt.Printf("Successfully scrobbled:\n")
	}
	fmt.Printf("  1. %s - %s (%s)\n", artistName, selectedTrack.Name, tsFormatted)

	if dryrun {
		return nil
	}

	username := profileFlag
	if username == "" {
		username, err = profile.GetActiveProfile()
		if err != nil {
			return err
		}
	}

	cred, err := keyring.GetCredential(username)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}

	if err := client.ScrobbleTrack(artistName, selectedTrack.Name, ts, cred.SessionKey); err != nil {
		return fmt.Errorf("failed to scrobble: %w", err)
	}

	return nil
}
