package word

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

//data types in json key value pairs
type WordEntry struct {
    Word         string   `json:"Word"`
    Meaning      string   `json:"Meaning"`
    Synonyms     *string  `json:"Synonyms"`
    Conjugation  *string  `json:"Conjugation"`
}

// LoadData reads the JSON data from the file and unmarshals it into a map with numbered keys.
func LoadData() (map[string]WordEntry, error) {
    // file location hard coded
	file, err := os.Open("C:\\Users\\user\\km\\data\\words.json")
    if err != nil {
        return nil, fmt.Errorf("failed to open data file: %v", err)
    }
    defer file.Close()

    var entries map[string]WordEntry
    byteValue, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("failed to read data file: %v", err)
    }

    if err := json.Unmarshal(byteValue, &entries); err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %v", err)
    }
    return entries, nil
}

func Search(word string) (*WordEntry, error) {
    entries, err := LoadData()
    if err != nil {
        return nil, err
    }

    for _, entry := range entries {
        if entry.Word == word {
            return &entry, nil
        }
    }
    return nil, fmt.Errorf("word '%s' not found", word)
}
