package cmd

import (
	"fmt"

	"github.com/boldandbrad/spin/internal/profile"
	"github.com/spf13/cobra"
)

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, err := profile.ListProfiles()
		if err != nil {
			return err
		}

		if len(profiles) == 0 {
			fmt.Println("No profiles found. Add one with: spin profile add <username>")
			return nil
		}

		activeProfile, _ := profile.GetActiveProfile()

		fmt.Println("Profiles:\n")
		for _, p := range profiles {
			marker := "  "
			if p.Username == activeProfile {
				marker = "* "
			}
			fmt.Printf("  %s%s\n", marker, p.Username)
		}

		return nil
	},
}
