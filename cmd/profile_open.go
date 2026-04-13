package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/boldandbrad/spin/internal/profile"
	"github.com/spf13/cobra"
)

var profileOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open profile in browser",
	Long:  `Open the active profile's Last.fm page in your default browser.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		username, err := profile.GetActiveProfile()
		if err != nil {
			return err
		}

		url := fmt.Sprintf("https://www.last.fm/user/%s", username)

		if err := openBrowser(url); err != nil {
			return fmt.Errorf("failed to open browser: %w", err)
		}

		fmt.Printf("Opening %s\n", url)
		return nil
	},
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", url).Run()
	case "linux":
		err = exec.Command("xdg-open", url).Run()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Run()
	default:
		err = fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return err
}
