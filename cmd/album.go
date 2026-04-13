package cmd

import (
	"fmt"
	"time"

	"github.com/boldandbrad/spin/internal/api"
	"github.com/boldandbrad/spin/internal/profile"
	"github.com/boldandbrad/spin/internal/scrobble"
	"github.com/boldandbrad/spin/tui"
	"github.com/spf13/cobra"
)

type trackResult struct {
	Name     string
	Duration int
}

func getAlbumTracks(artist, album string) ([]trackResult, error) {
	albumInfo, err := api.NewClient().GetAlbumInfo(artist, album)
	if err != nil {
		return nil, err
	}
	tracks := make([]trackResult, len(albumInfo.Album.Tracks.Track))
	for i, t := range albumInfo.Album.Tracks.Track {
		tracks[i] = trackResult{Name: t.Name, Duration: t.Duration}
	}
	return tracks, nil
}

type scrobbleInput struct {
	Artist    string
	Name      string
	Tracks    []trackResult
	Timestamp time.Time
	Profile   string
	Dryrun    bool
}

var albumCmd = &cobra.Command{
	Use:   "album [artist] [album]",
	Short: "Scrobble an album",
	Long: `Scrobble an album to last.fm.

If no arguments are provided, launches TUI mode for interactive scrobbling.
If artist and album are provided, scrobbles directly (CLI mode).`,
	Args: cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileFlag, _ := cmd.Flags().GetString("profile")
		endNow, _ := cmd.Flags().GetBool("end-now")
		dateFlag, _ := cmd.Flags().GetString("date")
		timestampFlag, _ := cmd.Flags().GetString("timestamp")
		dryrun, _ := cmd.Flags().GetBool("dryrun")

		var artist, name string
		var timeMode scrobble.TimeMode
		var customDate, customTime string

		if len(args) == 0 {
			input, err := tui.CollectAlbumInput()
			if err != nil {
				return err
			}
			if input == nil {
				return nil
			}
			artist = input.Artist
			name = input.Name
			timeMode = input.TimeMode
			customDate = input.Date
			customTime = input.Time
		} else {
			if len(args) != 2 {
				return fmt.Errorf("requires artist and album arguments")
			}
			artist = args[0]
			name = args[1]

			if endNow {
				timeMode = scrobble.TimeModeEndNow
			} else if timestampFlag != "" || dateFlag != "" {
				timeMode = scrobble.TimeModeCustom
				customTime = timestampFlag
				customDate = dateFlag
			}
		}

		tracks, err := getAlbumTracks(artist, name)
		if err != nil {
			return fmt.Errorf("failed to get album info: %w", err)
		}
		if len(tracks) == 0 {
			return fmt.Errorf("no tracks found for %s - %s", artist, name)
		}

		var totalDuration int
		for _, t := range tracks {
			totalDuration += t.Duration
		}

		timestamp, err := scrobble.ResolveTimestampFromMode(timeMode, customTime, customDate, totalDuration)
		if err != nil {
			return err
		}

		input := &scrobbleInput{
			Artist:    artist,
			Name:      name,
			Tracks:    tracks,
			Timestamp: timestamp,
			Profile:   profileFlag,
			Dryrun:    dryrun,
		}

		return scrobbleAlbum(input)
	},
}

func scrobbleAlbum(input *scrobbleInput) error {
	username, err := profile.ResolveProfile(input.Profile)
	if err != nil {
		return err
	}

	if input.Dryrun {
		fmt.Printf("Would scrobble to %s:\n", username)
	} else {
		fmt.Printf("Scrobbled to %s:\n", username)
	}

	currentTs := input.Timestamp
	for i, track := range input.Tracks {
		tsFormatted := currentTs.Format("2006-01-02 15:04")
		fmt.Printf("%2d. %s - %s (%s)\n", i+1, input.Artist, track.Name, tsFormatted)
		currentTs = currentTs.Add(time.Duration(track.Duration) * time.Second)
	}

	if input.Dryrun {
		return nil
	}

	p, err := profile.GetCredential(username)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}

	client := api.NewClient()
	currentTs = input.Timestamp
	for _, track := range input.Tracks {
		ts := scrobble.FormatTimestamp(currentTs)
		if err := client.ScrobbleTrack(input.Artist, track.Name, ts, p.SessionKey); err != nil {
			fmt.Printf("Warning: failed to scrobble %s: %v\n", track.Name, err)
		}
		currentTs = currentTs.Add(time.Duration(track.Duration) * time.Second)
	}

	return nil
}

func init() {
	albumCmd.Flags().Bool("end-now", false, "calculate start time from album duration")
	albumCmd.Flags().String("date", "", "date of listen (YYYY-MM-DD)")
	albumCmd.Flags().String("timestamp", "", "time of listen (HH:MM)")
	albumCmd.Flags().StringP("profile", "p", "", "profile to scrobble with")
	albumCmd.Flags().Bool("dryrun", false, "show what would be scrobbled without submitting")
}
