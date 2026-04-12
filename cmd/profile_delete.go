package cmd

import (
	"fmt"

	"github.com/boldandbrad/spin/internal/profile"
	"github.com/spf13/cobra"
)

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <lastfm-username>",
	Short: "Delete a profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]

		if err := profile.DeleteProfile(username); err != nil {
			return err
		}

		fmt.Printf("Profile %s deleted.\n", username)
		return nil
	},
}
