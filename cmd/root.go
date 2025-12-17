package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var version = "v0.1.0"

var rootCmd = &cobra.Command{
	Use:     "km",
	Version: version,
	Short:   "CLI to search word meanings",
	Long:    "SEARCH FOR MEANING OF KISWAHILI WORDS | TAFUTA MAANA YA MANENO YA KISWAHILI",
	Example: `  km search funua
  km search shamba
  km s chakula`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func init() {
	// Add subcommands here
}
