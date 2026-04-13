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

var albumCmd = &cobra.Command{
	Use:   "album [artist] [album]",
	Short: "Scrobble an album",
	Long: `Scrobble an album to last.fm.

If no arguments are provided, launches TUI mode for interactive scrobbling.
If artist and album are provided, scrobbles directly (CLI mode).`,
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
			return tui.AlbumTUI(profileFlag, dryrun)
		}

		if len(args) != 2 {
			return fmt.Errorf("requires artist and album arguments")
		}

		artist := args[0]
		album := args[1]

		client := api.NewClient()
		albumInfo, err := client.GetAlbumInfo(artist, album)
		if err != nil {
			return fmt.Errorf("failed to get album info: %w", err)
		}

		if len(albumInfo.Album.Tracks.Track) == 0 {
			return fmt.Errorf("no tracks found for %s - %s", artist, album)
		}

		var timestamp time.Time
		var totalDuration int
		for _, track := range albumInfo.Album.Tracks.Track {
			totalDuration += track.Duration
		}

		if endTimeFlag != "" {
			t, err := scrobble.ParseTimeOfDay(endTimeFlag)
			if err != nil {
				return fmt.Errorf("invalid --end-time: %w", err)
			}
			timestamp = t.Add(-time.Duration(totalDuration) * time.Second)
		} else if startTimeFlag != "" {
			t, err := scrobble.ParseTimeOfDay(startTimeFlag)
			if err != nil {
				return fmt.Errorf("invalid --start-time: %w", err)
			}
			timestamp = t
		} else {
			timestamp = time.Now()
		}

		timestamp, err = scrobble.ApplyDateToTimestamp(dateFlag, timestamp)
		if err != nil {
			return fmt.Errorf("invalid --date: %w", err)
		}

		username, err := profile.ResolveProfile(profileFlag)
		if err != nil {
			return err
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
				fmt.Fprintf(cmd.OutOrStderr(), "Warning: failed to scrobble %s: %v\n", track.Name, err)
			}
			currentTimestamp = currentTimestamp.Add(time.Duration(track.Duration) * time.Second)
		}

		return nil
	},
}

func init() {
	albumCmd.Flags().String("start-time", "", "start time of listen (HH:MM)")
	albumCmd.Flags().String("end-time", "", "end time of listen (HH:MM)")
	albumCmd.Flags().String("date", "", "date of listen (YYYY-MM-DD)")
	albumCmd.Flags().StringP("profile", "p", "", "profile to scrobble with")
	albumCmd.Flags().Bool("dryrun", false, "show what would be scrobbled without submitting")
}
