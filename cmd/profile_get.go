package cmd

import (
	"fmt"

	"github.com/boldandbrad/spin/internal/profile"
	"github.com/spf13/cobra"
)

var profileGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the active profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, err := profile.GetActiveProfile()
		if err != nil {
			return err
		}

		fmt.Println(username)
		return nil
	},
}
