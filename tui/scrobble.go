package tui

import (
	"fmt"

	"charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/boldandbrad/spin/internal/scrobble"
)

type Input struct {
	Artist   string
	Track    string
	Album    string
	TimeMode scrobble.TimeMode
	Date     string
	Time     string
}

func CollectInput(isAlbum bool) (*Input, error) {
	artist := ""
	track := ""
	album := ""
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

	trackField := huh.NewInput().
		Title("Track").
		Value(&track).
		Placeholder("e.g., Paranoid Android").
		Validate(func(s string) error {
			if s == "" && !isAlbum {
				return fmt.Errorf("track is required")
			}
			return nil
		})

	albumField := huh.NewInput().
		Title("Album").
		Value(&album).
		Placeholder("e.g., OK Computer").
		Validate(func(s string) error {
			if s == "" && isAlbum {
				return fmt.Errorf("album is required")
			}
			return nil
		})

	formFields := []huh.Field{artistField}
	if isAlbum {
		formFields = append(formFields, albumField)
	} else {
		formFields = append(formFields, trackField)
		albumField.Title("Album (optional)")
		formFields = append(formFields, albumField)
	}

	form := huh.NewForm(
		huh.NewGroup(formFields...),
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
	).WithProgramOptions(tea.WithoutCatchPanics())

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
		).WithProgramOptions(tea.WithoutCatchPanics())
		if err := timeForm.Run(); err != nil {
			return nil, err
		}
	}

	return &Input{
		Artist:   artist,
		Track:    track,
		Album:    album,
		TimeMode: timeMode,
		Date:     date,
		Time:     timeStr,
	}, nil
}

func CollectTrackInput() (*Input, error) {
	return CollectInput(false)
}

func CollectAlbumInput() (*Input, error) {
	return CollectInput(true)
}
