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

var albumCmd = &cobra.Command{
	Use:   "album [artist] [album]",
	Short: "Scrobble an album",
	Long: `Scrobble an album to last.fm.

If no arguments are provided, launches TUI mode for interactive scrobbling.
If artist and album are provided, scrobbles directly (CLI mode).`,
	Args: cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileFlag, _ := cmd.Flags().GetString("profile")
		dateFlag, _ := cmd.Flags().GetString("date")
		timeFlag, _ := cmd.Flags().GetString("timestamp")
		dryrun, _ := cmd.Flags().GetBool("dryrun")

		if len(args) == 0 {
			return tui.AlbumTUI(profileFlag, dryrun)
		}

		if len(args) != 2 {
			return fmt.Errorf("requires artist and album arguments")
		}

		artist := args[0]
		album := args[1]

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

		client := api.NewClient()
		albumInfo, err := client.GetAlbumInfo(artist, album)
		if err != nil {
			return fmt.Errorf("failed to get album info: %w", err)
		}

		if len(albumInfo.Album.Tracks.Track) == 0 {
			return fmt.Errorf("no tracks found for %s - %s", artist, album)
		}

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
				fmt.Fprintf(cmd.OutOrStderr(), "Warning: failed to scrobble %s: %v\n", track.Name, err)
			}
		}

		return nil
	},
}

func init() {
	albumCmd.Flags().StringP("date", "d", "", "date of listen (YYYY-MM-DD)")
	albumCmd.Flags().StringP("timestamp", "t", "", "time of listen (HH:MM)")
	albumCmd.Flags().StringP("profile", "p", "", "profile to scrobble with")
	albumCmd.Flags().Bool("dryrun", false, "show what would be scrobbled without submitting")
}
