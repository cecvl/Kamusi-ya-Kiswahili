package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "km",
    Short: "CLI to search word meanings",
    Long:  "KAMUSI YA KISWAHILI : MAANA YA MANENO YA KISWAHILI",
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        log.Fatalf("Error: %v\n", err)
    }
}

func init() {
    // Add subcommands here
}
