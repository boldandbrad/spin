package cmd

import (
	"fmt"

	"charm.land/bubbletea/v2"
	"charm.land/huh/v2"
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

		var password string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Password").
					Placeholder("last.fm password").
					Password(true).
					Value(&password).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("password is required")
						}
						return nil
					}),
			),
		).WithProgramOptions(tea.WithoutCatchPanics())

		if err := form.Run(); err != nil {
			return err
		}

		client := api.NewClient()

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
}
