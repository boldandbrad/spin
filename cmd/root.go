package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	debug   bool
)

var rootCmd = &cobra.Command{
	Use:   "spin",
	Short: "A command line last.fm scrobbler for techies",
	Long: `Interactively or programmatically scrobble tracks and albums to last.fm from the
terminal.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "enable debug logging")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			log.SetOutput(os.Stderr)
		} else {
			log.SetOutput(io.Discard)
		}
	}

	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(trackCmd)
	rootCmd.AddCommand(albumCmd)
	rootCmd.AddCommand(historyCmd)
}
