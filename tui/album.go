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

func AlbumTUI(profileFlag string, dryrun bool) error {
	artist := ""
	album := ""
	dateInput := ""
	timeInput := ""

	inputForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Artist").Value(&artist).Placeholder("e.g., Radiohead"),
			huh.NewInput().Title("Album").Value(&album).Placeholder("e.g., OK Computer"),
		),
		huh.NewGroup(
			huh.NewInput().Title("Date (YYYY-MM-DD)").Value(&dateInput).Placeholder("optional"),
			huh.NewInput().Title("Time (HH:MM)").Value(&timeInput).Placeholder("optional"),
		),
	)

	if err := inputForm.Run(); err != nil {
		return err
	}

	if artist == "" || album == "" {
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
	albumInfo, err := client.GetAlbumInfo(artist, album)
	if err != nil {
		return fmt.Errorf("failed to get album info: %w", err)
	}

	if len(albumInfo.Album.Tracks.Track) == 0 {
		fmt.Println("No tracks found.")
		return nil
	}

	ts := scrobble.FormatTimestamp(timestamp)
	tsFormatted := timestamp.Format("2006-01-02 15:04")

	if dryrun {
		fmt.Printf("Would scrobble:\n")
	} else {
		fmt.Printf("Successfully scrobbled:\n")
	}
	for i, track := range albumInfo.Album.Tracks.Track {
		fmt.Printf("%2d. %s - %s (%s)\n", i+1, artist, track.Name, tsFormatted)
	}

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

	for _, track := range albumInfo.Album.Tracks.Track {
		if err := client.ScrobbleTrack(artist, track.Name, ts, cred.SessionKey); err != nil {
			fmt.Printf("Warning: failed to scrobble %s: %v\n", track.Name, err)
		}
	}

	return nil
}
