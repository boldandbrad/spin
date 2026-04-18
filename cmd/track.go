package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/boldandbrad/spin/internal/api"
	"github.com/boldandbrad/spin/internal/profile"
	"github.com/boldandbrad/spin/internal/scrobble"
	"github.com/boldandbrad/spin/tui"
	"github.com/spf13/cobra"
)

var trackCmd = &cobra.Command{
	Use:           "track [artist] [track]",
	Short:         "Scrobble a track",
	SilenceUsage:  true,
	SilenceErrors: true,
	Long: `Scrobble a track to last.fm.

If no arguments are provided, launches TUI mode for interactive scrobbling.
If artist and track are provided, scrobbles directly (CLI mode).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		profileFlag, _ := cmd.Flags().GetString("profile")
		endNow, _ := cmd.Flags().GetBool("end-now")
		dateFlag, _ := cmd.Flags().GetString("date")
		timestampFlag, _ := cmd.Flags().GetString("timestamp")
		dryrun, _ := cmd.Flags().GetBool("dryrun")

		var artist, track, album string
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
			track = input.Track
			album = input.Album
			timeMode = input.TimeMode
			customDate = input.Date
			customTime = input.Time
		} else {
			if len(args) != 2 {
				fmt.Fprintf(os.Stderr, "Error: requires artist and track arguments\n\n")
				cmd.Usage()
				return nil
			}
			artist = args[0]
			track = args[1]
			albumFlag, _ := cmd.Flags().GetString("album")
			album = albumFlag

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

		client := api.NewClient()
		trackMetadata, err := client.GetTrackInfo(artist, track)
		if err != nil {
			return err
		}

		totalDuration := 0
		if timeMode == scrobble.TimeModeEndNow {
			if trackMetadata.Duration == 0 {
				return fmt.Errorf("last.fm doesn't have duration for this track, cannot use --end-now")
			}
			totalDuration = trackMetadata.Duration / 1000
		}

		timestamp, err := scrobble.ResolveTimestampFromMode(timeMode, customTime, customDate, totalDuration)
		if err != nil {
			return err
		}

		albumName := ""
		if trackMetadata != nil {
			artist = trackMetadata.Artist
			track = trackMetadata.Track

			if album != "" {
				if validatedAlbum, _ := client.ValidateAlbumForTrack(artist, track, album); validatedAlbum != "" {
					albumName = validatedAlbum
				} else {
					albumName = trackMetadata.Album
				}
			} else {
				albumName = trackMetadata.Album
			}
		}

		return scrobbleTrack(artist, track, albumName, timestamp, profileFlag, dryrun)
	},
}

func printTrack(artist, track, album, timestamp string) {
	if album != "" {
		fmt.Printf("%2d. %s - %s (%s) (%s)\n", 1, artist, track, album, timestamp)
	} else {
		fmt.Printf("%2d. %s - %s (%s)\n", 1, artist, track, timestamp)
	}
}

func scrobbleTrack(artist, track, album string, timestamp time.Time, profileFlag string, dryrun bool) error {
	ts := scrobble.FormatTimestamp(timestamp)
	tsFormatted := timestamp.Format("2006-01-02 15:04")

	username, err := profile.ResolveProfile(profileFlag)
	if err != nil {
		return err
	}

	if dryrun {
		cliCmd := buildTrackCLICommand(artist, track, album, timestamp)
		fmt.Printf("Would scrobble to %s:\n\n", username)
		printTrack(artist, track, album, tsFormatted)

		fmt.Printf("\nRun this command to scrobble:\n  %s\n\n", cliCmd)

		if askCopyToClipboard() {
			if err := copyToClipboard(cliCmd); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to copy: %v\n", err)
			} else {
				fmt.Println("Command copied to clipboard!")
			}
		}

		return nil
	}

	p, err := profile.GetCredential(username)
	if err != nil {
		return fmt.Errorf("failed to get credential: %w", err)
	}

	client := api.NewClient()
	if err := client.ScrobbleTrack(artist, track, ts, p.SessionKey, album); err != nil {
		return fmt.Errorf("%2d. %s - %s (%s): failed to scrobble: %w", 1, artist, track, album, err)
	}

	fmt.Printf("Scrobbled to %s:\n\n", username)
	printTrack(artist, track, album, tsFormatted)

	return nil
}

func init() {
	trackCmd.Flags().Bool("end-now", false, "calculate start time from track duration")
	trackCmd.Flags().String("date", "", "date of listen (YYYY-MM-DD)")
	trackCmd.Flags().String("timestamp", "", "time of listen (HH:MM)")
	trackCmd.Flags().StringP("profile", "p", "", "profile to scrobble with")
	trackCmd.Flags().String("album", "", "album to scrobble with (optional)")
	trackCmd.Flags().Bool("dryrun", false, "show what would be scrobbled without submitting")
}
