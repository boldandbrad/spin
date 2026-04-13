package tui

import (
	"fmt"
	"os"

	"charm.land/huh/v2"
	"github.com/boldandbrad/spin/internal/scrobble"
)

type ScrobbleInput struct {
	Artist   string
	Name     string
	TimeMode scrobble.TimeMode
	Date     string
	Time     string
}

func CollectInput(isAlbum bool) (*ScrobbleInput, error) {
	artist := ""
	name := ""
	timeMode := scrobble.TimeModeEndNow

	artistField := huh.NewInput().
		Title("Artist").
		Value(&artist).
		Placeholder("e.g., Radiohead").
		Validate(func(s string) error {
			if s == "" {
				return fmt.Errorf("artist is required")
			}
			return nil
		})

	nameField := huh.NewInput().
		Title(func() string {
			if isAlbum {
				return "Album"
			}
			return "Track"
		}()).
		Value(&name).
		Placeholder(func() string {
			if isAlbum {
				return "e.g., OK Computer"
			}
			return "e.g., Paranoid Android"
		}()).
		Validate(func(s string) error {
			if s == "" {
				if isAlbum {
					return fmt.Errorf("album is required")
				}
				return fmt.Errorf("track is required")
			}
			return nil
		})

	form := huh.NewForm(
		huh.NewGroup(
			artistField,
			nameField,
		),
		huh.NewGroup(
			huh.NewSelect[scrobble.TimeMode]().
				Title("When did you listen?").
				Options(
					huh.NewOption("Starting now", scrobble.TimeModeStartNow),
					huh.NewOption("Ending now", scrobble.TimeModeEndNow),
					huh.NewOption("Custom start time", scrobble.TimeModeCustom),
				).
				Value(&timeMode),
		),
	)

	os.Setenv("TEA_LOG", "")
	if err := form.Run(); err != nil {
		return nil, err
	}

	date := ""
	timeStr := ""
	if timeMode == scrobble.TimeModeCustom {
		timeForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Date (YYYY-MM-DD)").Value(&date).Placeholder("e.g., 2026-04-12 (optional)"),
				huh.NewInput().Title("Time (HH:MM)").Value(&timeStr).Placeholder("e.g., 15:00").Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("time is required")
					}
					return nil
				}),
			),
		)
		if err := timeForm.Run(); err != nil {
			return nil, err
		}
	}

	return &ScrobbleInput{
		Artist:   artist,
		Name:     name,
		TimeMode: timeMode,
		Date:     date,
		Time:     timeStr,
	}, nil
}

func CollectTrackInput() (*ScrobbleInput, error) {
	return CollectInput(false)
}

func CollectAlbumInput() (*ScrobbleInput, error) {
	return CollectInput(true)
}
