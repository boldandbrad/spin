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
		startTimeFlag, _ := cmd.Flags().GetString("start-time")
		endTimeFlag, _ := cmd.Flags().GetString("end-time")
		dateFlag, _ := cmd.Flags().GetString("date")
		dryrun, _ := cmd.Flags().GetBool("dryrun")

		if startTimeFlag != "" && endTimeFlag != "" {
			return fmt.Errorf("cannot specify both --start-time and --end-time")
		}

		if len(args) == 0 {
			return tui.TrackTUI(profileFlag, dryrun)
		}

		if len(args) != 2 {
			return fmt.Errorf("requires artist and track arguments")
		}

		artist := args[0]
		track := args[1]

		client := api.NewClient()
		timestamp := time.Now()

		if startTimeFlag != "" {
			t, err := scrobble.ParseTimeOfDay(startTimeFlag)
			if err != nil {
				return fmt.Errorf("invalid --start-time: %w", err)
			}
			timestamp = t
		} else if endTimeFlag != "" {
			t, err := scrobble.ParseTimeOfDay(endTimeFlag)
			if err != nil {
				return fmt.Errorf("invalid --end-time: %w", err)
			}
			duration, err := client.GetTrackInfo(artist, track)
			if err != nil {
				return fmt.Errorf("failed to get track info: %w", err)
			}
			if duration == 0 {
				return fmt.Errorf("track duration unknown, cannot use --end-time")
			}
			timestamp = t.Add(-time.Duration(duration) * time.Millisecond)
		}

		timestamp, err := scrobble.ApplyDateToTimestamp(dateFlag, timestamp)
		if err != nil {
			return fmt.Errorf("invalid --date: %w", err)
		}

		username, err := profile.ResolveProfile(profileFlag)
		if err != nil {
			return err
		}

		ts := scrobble.FormatTimestamp(timestamp)
		tsFormatted := timestamp.Format("2006-01-02 15:04")

		if dryrun {
			fmt.Printf("Would scrobble to %s:\n", username)
		} else {
			fmt.Printf("Scrobbled to %s:\n", username)
		}
		fmt.Printf("  1. %s - %s (%s)\n", artist, track, tsFormatted)

		if dryrun {
			return nil
		}

		cred, err := profile.GetCredentialForProfile(profileFlag)
		if err != nil {
			return err
		}

		if err := client.ScrobbleTrack(artist, track, ts, cred.SessionKey); err != nil {
			return fmt.Errorf("failed to scrobble: %w", err)
		}

		return nil
	},
}

func init() {
	trackCmd.Flags().String("start-time", "", "start time of listen (HH:MM)")
	trackCmd.Flags().String("end-time", "", "end time of listen (HH:MM)")
	trackCmd.Flags().String("date", "", "date of listen (YYYY-MM-DD)")
	trackCmd.Flags().StringP("profile", "p", "", "profile to scrobble with")
	trackCmd.Flags().Bool("dryrun", false, "show what would be scrobbled without submitting")
}
