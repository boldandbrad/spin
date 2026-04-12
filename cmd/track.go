package cmd

import (
	"fmt"
	"time"

	"github.com/boldandbrad/spin/internal/api"
	"github.com/boldandbrad/spin/internal/keyring"
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
		dateFlag, _ := cmd.Flags().GetString("date")
		timeFlag, _ := cmd.Flags().GetString("timestamp")
		dryrun, _ := cmd.Flags().GetBool("dryrun")

		if len(args) == 0 {
			return tui.TrackTUI(profileFlag, dryrun)
		}

		if len(args) != 2 {
			return fmt.Errorf("requires artist and track arguments")
		}

		artist := args[0]
		track := args[1]

		username := profileFlag
		if username == "" {
			var err error
			username, err = profile.GetActiveProfile()
			if err != nil {
				return err
			}
		}

		timestamp := time.Now()
		if timeFlag != "" {
			var err error
			timestamp, err = scrobble.ParseTimeOfDay(timeFlag)
			if err != nil {
				return err
			}
		}
		if dateFlag != "" {
			var err error
			timestamp, err = scrobble.ParseDateTime(dateFlag, timestamp.Format("15:04"))
			if err != nil {
				return err
			}
		}

		ts := scrobble.FormatTimestamp(timestamp)
		tsFormatted := timestamp.Format("2006-01-02 15:04")

		if dryrun {
			fmt.Printf("Would scrobble:\n")
		} else {
			fmt.Printf("Successfully scrobbled:\n")
		}
		fmt.Printf("  1. %s - %s (%s)\n", artist, track, tsFormatted)

		if dryrun {
			return nil
		}

		cred, err := keyring.GetCredential(username)
		if err != nil {
			return fmt.Errorf("failed to get credential: %w", err)
		}

		client := api.NewClient()
		if err := client.ScrobbleTrack(artist, track, ts, cred.SessionKey); err != nil {
			return fmt.Errorf("failed to scrobble: %w", err)
		}

		return nil
	},
}

func init() {
	trackCmd.Flags().StringP("date", "d", "", "date of listen (YYYY-MM-DD)")
	trackCmd.Flags().StringP("timestamp", "t", "", "time of listen (HH:MM)")
	trackCmd.Flags().StringP("profile", "p", "", "profile to scrobble with")
	trackCmd.Flags().Bool("dryrun", false, "show what would be scrobbled without submitting")
}
