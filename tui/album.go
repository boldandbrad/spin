package tui

import (
	"fmt"
	"time"

	"github.com/boldandbrad/spin/internal/api"
	"github.com/boldandbrad/spin/internal/profile"
	"github.com/boldandbrad/spin/internal/scrobble"
	"github.com/charmbracelet/huh"
)

func AlbumTUI(profileFlag string, dryrun bool) error {
	artist := ""
	album := ""
	timestampMode := 1
	dateInput := ""
	timeInput := ""

	inputForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Artist").Value(&artist).Placeholder("e.g., Radiohead"),
			huh.NewInput().Title("Album").Value(&album).Placeholder("e.g., OK Computer"),
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

	if artist == "" || album == "" {
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
	albumInfo, err := client.GetAlbumInfo(artist, album)
	if err != nil {
		return fmt.Errorf("failed to get album info: %w", err)
	}

	if len(albumInfo.Album.Tracks.Track) == 0 {
		fmt.Println("No tracks found.")
		return nil
	}

	username, err := profile.ResolveProfile(profileFlag)
	if err != nil {
		return err
	}

	var timestamp time.Time
	switch timestampMode {
	case 0:
		timestamp = time.Now()
	case 1:
		var totalDuration int
		for _, track := range albumInfo.Album.Tracks.Track {
			totalDuration += track.Duration
		}
		timestamp = time.Now().Add(-time.Duration(totalDuration) * time.Second)
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

	if dryrun {
		fmt.Printf("Would scrobble to %s:\n", username)
	} else {
		fmt.Printf("Scrobbled to %s:\n", username)
	}

	currentTimestamp := timestamp
	for i, track := range albumInfo.Album.Tracks.Track {
		tsFormatted := currentTimestamp.Format("2006-01-02 15:04")
		fmt.Printf("%2d. %s - %s (%s)\n", i+1, artist, track.Name, tsFormatted)
		currentTimestamp = currentTimestamp.Add(time.Duration(track.Duration) * time.Second)
	}

	if dryrun {
		return nil
	}

	cred, err := profile.GetCredentialForProfile(profileFlag)
	if err != nil {
		return err
	}

	currentTimestamp = timestamp
	for _, track := range albumInfo.Album.Tracks.Track {
		ts := scrobble.FormatTimestamp(currentTimestamp)
		if err := client.ScrobbleTrack(artist, track.Name, ts, cred.SessionKey); err != nil {
			fmt.Printf("Warning: failed to scrobble %s: %v\n", track.Name, err)
		}
		currentTimestamp = currentTimestamp.Add(time.Duration(track.Duration) * time.Second)
	}

	return nil
}
