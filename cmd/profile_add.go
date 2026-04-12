package cmd

import (
	"fmt"
	"os"

	"github.com/boldandbrad/spin/internal/api"
	"github.com/boldandbrad/spin/internal/profile"
	"github.com/spf13/cobra"
)

var profileAddCmd = &cobra.Command{
	Use:   "add <lastfm-username>",
	Short: "Add a last.fm profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]

		fmt.Printf("Adding profile for %s...\n", username)

		client := api.NewClient()

		fmt.Print("Enter last.fm password: ")
		var password string
		fmt.Fscanln(os.Stdin, &password)

		if password == "" {
			return fmt.Errorf("password is required")
		}

		sessionKey, err := client.GetSessionKey(username, password)
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}

		if err := profile.AddProfile(username, sessionKey); err != nil {
			return fmt.Errorf("failed to add profile: %w", err)
		}

		fmt.Printf("Profile %s added successfully!\n", username)
		return nil
	},
}

func init() {
	profileAddCmd.Flags().StringP("password", "p", "", "last.fm password (will prompt if not provided)")
}
