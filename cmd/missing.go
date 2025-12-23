package cmd

import (
	"fmt"
	"km/pkg/word"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var missingCmd = &cobra.Command{
	Use:     "missing",
	Aliases: []string{"m"},
	Short:   "Show missing words | Onyesha maneno yaliyokosekana",
	Long:    "Display the most frequently searched words that were not found in the dictionary",
	Example: `  km missing
  km missing --limit 20
  km m`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")

		missing, err := word.GetMissingWords(limit)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(missing) == 0 {
			fmt.Println("Hakuna maneno yaliyokosekana | No missing words recorded yet.")
			return
		}

		fmt.Printf("\nManeno yaliyotafutwa lakini hayakupatikana | Words searched but not found:\n\n")

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NENO\tIDADI\tMUDA")
		fmt.Fprintln(w, "WORD\tCOUNT\tLAST SEARCHED")
		fmt.Fprintln(w, "----\t-----\t-------------")

		for _, item := range missing {
			fmt.Fprintf(w, "%s\t%d\t%s\n",
				item.Word,
				item.SearchCount,
				item.LastSearched.Format("2006-01-02 15:04:05"))
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(missingCmd)
	missingCmd.Flags().IntP("limit", "l", 10, "Number of missing words to display")
}
