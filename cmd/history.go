package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/boldandbrad/spin/internal/api"
	"github.com/boldandbrad/spin/internal/profile"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View recent scrobbles",
	Long:  `View recent scrobbles from the active profile.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		limitFlag, _ := cmd.Flags().GetInt("limit")

		username, err := profile.GetActiveProfile()
		if err != nil {
			return err
		}

		client := api.NewClient()
		tracks, err := client.GetRecentTracks(username, limitFlag)
		if err != nil {
			return fmt.Errorf("failed to get history: %w", err)
		}

		if len(tracks) == 0 {
			fmt.Println("No recent scrobbles found.")
			return nil
		}

		fmt.Printf("Recent scrobbles for %s:\n\n", username)
		for i, track := range tracks {
			dateStr := ""
			if track.Date.UTS != "" {
				ts, _ := strconv.ParseInt(track.Date.UTS, 10, 64)
				date := time.Unix(ts, 0)
				dateStr = date.Format("2006-01-02 15:04")
			}
			artist := track.Artist.Text
			if artist == "" {
				artist = "Unknown Artist"
			}
			if dateStr != "" {
				fmt.Printf("%2d. %s - %s (%s)\n", i+1, artist, track.Name, dateStr)
			} else {
				fmt.Printf("%2d. %s - %s\n", i+1, artist, track.Name)
			}
		}

		return nil
	},
}

func init() {
	historyCmd.Flags().IntP("limit", "n", 10, "number of results to show")
}
