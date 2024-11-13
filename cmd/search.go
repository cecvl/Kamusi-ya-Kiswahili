package cmd

import (
	"fmt"
	"km/pkg/word"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
    Use:   "s [word]",
    Short: "Search for a word",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        result, err := word.Search(args[0])
        if err != nil {
            fmt.Println("Error:", err)
            return
        }
        //fmt.Printf("Word: %s\nMeaning: %s\nSynonyms: %v\nConjugations: %v\n", 
          //  result.Word, result.Meaning, result.Synonyms, result.Conjugation)
		  fmt.Printf("Maana yake: %s\n", result.Meaning)
    },
}

func init() {
    rootCmd.AddCommand(searchCmd)
}
