package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func buildTrackCLICommand(artist, track, album string, timestamp time.Time) string {
	cmd := fmt.Sprintf("spin track %q %q", artist, track)

	if album != "" {
		cmd += fmt.Sprintf(" --album %q", album)
	}

	cmd += fmt.Sprintf(" --timestamp %s", timestamp.Format("15:04"))
	cmd += fmt.Sprintf(" --date %s", timestamp.Format("2006-01-02"))

	return cmd
}

func buildAlbumCLICommand(artist, album string, timestamp time.Time) string {
	cmd := fmt.Sprintf("spin album %q %q", artist, album)

	cmd += fmt.Sprintf(" --timestamp %s", timestamp.Format("15:04"))
	cmd += fmt.Sprintf(" --date %s", timestamp.Format("2006-01-02"))

	return cmd
}

func askCopyToClipboard() bool {
	fmt.Print("Copy command to clipboard? [y/N] ")
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}

func copyToClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("bash", "-c", "echo -n "+fmt.Sprintf("%q", text)+" | xclip -selection clipboard")
	case "windows":
		cmd = exec.Command("cmd", "/c", "echo "+text+"| clip")
	default:
		return fmt.Errorf("unsupported platform")
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}
