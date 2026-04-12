package cmd

import (
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage last.fm profiles",
	Long:  `Add, list, set, get, and delete last.fm profiles.`,
}

func init() {
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileSetCmd)
	profileCmd.AddCommand(profileGetCmd)
	profileCmd.AddCommand(profileDeleteCmd)
}
