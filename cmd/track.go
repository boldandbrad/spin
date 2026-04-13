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

var trackCmd = &cobra.Command{
	Use:   "track [artist] [track]",
	Short: "Scrobble a track",
	Long: `Scrobble a track to last.fm.

If no arguments are provided, launches TUI mode for interactive scrobbling.
If artist and track are provided, scrobbles directly (CLI mode).`,
	Args: cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileFlag, _ := cmd.Flags().GetString("profile")
		endNow, _ := cmd.Flags().GetBool("end-now")
		dateFlag, _ := cmd.Flags().GetString("date")
		timestampFlag, _ := cmd.Flags().GetString("timestamp")
		dryrun, _ := cmd.Flags().GetBool("dryrun")

		var artist, track string
		var timeMode scrobble.TimeMode
		var customDate, customTime string

		if len(args) == 0 {
			input, err := tui.CollectTrackInput()
			if err != nil {
				return err
			}
			if input == nil {
				return nil
			}
			artist = input.Artist
			track = input.Name
			timeMode = input.TimeMode
			customDate = input.Date
			customTime = input.Time
		} else {
			if len(args) != 2 {
				return fmt.Errorf("requires artist and track arguments")
			}
			artist = args[0]
			track = args[1]

			if endNow {
				timeMode = scrobble.TimeModeEndNow
			} else if timestampFlag != "" || dateFlag != "" {
				timeMode = scrobble.TimeModeCustom
				customTime = timestampFlag
				customDate = dateFlag
			} else {
				timeMode = scrobble.TimeModeStartNow
			}
		}

		totalDuration := 0
		if timeMode == scrobble.TimeModeEndNow {
			client := api.NewClient()
			durationMs, err := client.GetTrackInfo(artist, track)
			if err != nil {
				return fmt.Errorf("failed to get track info: %w", err)
			}
			if durationMs == 0 {
				return fmt.Errorf("track duration unknown, cannot use --end-now")
			}
			totalDuration = durationMs / 1000
		}

		timestamp, err := scrobble.ResolveTimestampFromMode(timeMode, customTime, customDate, totalDuration)
		if err != nil {
			return err
		}

		return scrobbleTrack(artist, track, timestamp, profileFlag, dryrun)
	},
}

func scrobbleTrack(artist, track string, timestamp time.Time, profileFlag string, dryrun bool) error {
	ts := scrobble.FormatTimestamp(timestamp)
	tsFormatted := timestamp.Format("2006-01-02 15:04")

	username, err := profile.ResolveProfile(profileFlag)
	if err != nil {
		return err
	}

	if dryrun {
		fmt.Printf("Would scrobble to %s:\n", username)
		fmt.Printf("  1. %s - %s (%s)\n", artist, track, tsFormatted)
		return nil
	}

	p, err := profile.GetCredential(username)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}

	client := api.NewClient()
	if err := client.ScrobbleTrack(artist, track, ts, p.SessionKey); err != nil {
		return fmt.Errorf("  1. %s - %s: failed to scrobble: %w", artist, track, err)
	}

	fmt.Printf("Scrobbled to %s:\n", username)
	fmt.Printf("  1. %s - %s (%s)\n", artist, track, tsFormatted)

	return nil
}

func init() {
	trackCmd.Flags().Bool("end-now", false, "calculate start time from track duration")
	trackCmd.Flags().String("date", "", "date of listen (YYYY-MM-DD)")
	trackCmd.Flags().String("timestamp", "", "time of listen (HH:MM)")
	trackCmd.Flags().StringP("profile", "p", "", "profile to scrobble with")
	trackCmd.Flags().Bool("dryrun", false, "show what would be scrobbled without submitting")
}
