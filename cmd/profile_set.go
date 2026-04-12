package cmd

import (
	"fmt"

	"github.com/boldandbrad/spin/internal/profile"
	"github.com/spf13/cobra"
)

var profileSetCmd = &cobra.Command{
	Use:   "set <lastfm-username>",
	Short: "Set the active profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		username := args[0]

		if err := profile.SetActiveProfile(username); err != nil {
			return err
		}

		fmt.Printf("Active profile set to: %s\n", username)
		return nil
	},
}
