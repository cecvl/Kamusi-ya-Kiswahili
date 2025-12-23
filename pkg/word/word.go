package word

import (
	"database/sql"
	"fmt"
	"km/pkg/database"
	"time"
)

// WordEntry represents a word entry in the dictionary
type WordEntry struct {
	ID          int64
	Word        string
	Meaning     string
	Synonyms    *string
	Conjugation *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Search performs a case-insensitive search for a word in the database
// This function is thread-safe and can be called concurrently
func Search(word string) (*WordEntry, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %v", err)
	}

	var entry WordEntry
	var synonyms, conjugation sql.NullString

	// Use parameterized query with COLLATE NOCASE for case-insensitive search
	query := `
		SELECT id, word, meaning, synonyms, conjugation, created_at, updated_at
		FROM words
		WHERE word = ? COLLATE NOCASE
		LIMIT 1
	`

	err = db.QueryRow(query, word).Scan(
		&entry.ID,
		&entry.Word,
		&entry.Meaning,
		&synonyms,
		&conjugation,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Track missing word
		if trackErr := trackMissingWord(word); trackErr != nil {
			// Log but don't fail the search
			fmt.Printf("Warning: failed to track missing word: %v\n", trackErr)
		}
		return nil, fmt.Errorf("word '%s' not found", word)
	}

	if err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}

	if synonyms.Valid {
		entry.Synonyms = &synonyms.String
	}
	if conjugation.Valid {
		entry.Conjugation = &conjugation.String
	}

	return &entry, nil
}

// SearchMultiple searches for multiple words concurrently
// Returns a map of word -> result (or error)
func SearchMultiple(words []string) map[string]*WordEntry {
	results := make(map[string]*WordEntry)
	resultChan := make(chan struct {
		word  string
		entry *WordEntry
	}, len(words))

	// Search words concurrently
	for _, w := range words {
		go func(word string) {
			entry, _ := Search(word)
			resultChan <- struct {
				word  string
				entry *WordEntry
			}{word, entry}
		}(w)
	}

	// Collect results
	for i := 0; i < len(words); i++ {
		result := <-resultChan
		results[result.word] = result.entry
	}

	return results
}

// trackMissingWord records a word that couldn't be found
// This function uses INSERT OR UPDATE pattern for concurrency safety
func trackMissingWord(word string) error {
	db, err := database.GetDB()
	if err != nil {
		return err
	}

	query := `
		INSERT INTO missing_words (word, search_count, first_searched_at, last_searched_at)
		VALUES (?, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(word) DO UPDATE SET
			search_count = search_count + 1,
			last_searched_at = CURRENT_TIMESTAMP
	`

	_, err = db.Exec(query, word)
	return err
}

// GetMissingWords returns the most frequently searched missing words
func GetMissingWords(limit int) ([]struct {
	Word         string
	SearchCount  int
	LastSearched time.Time
}, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT word, search_count, last_searched_at
		FROM missing_words
		ORDER BY search_count DESC
		LIMIT ?
	`

	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		Word         string
		SearchCount  int
		LastSearched time.Time
	}

	for rows.Next() {
		var item struct {
			Word         string
			SearchCount  int
			LastSearched time.Time
		}
		if err := rows.Scan(&item.Word, &item.SearchCount, &item.LastSearched); err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, rows.Err()
}

// FuzzySearch performs a fuzzy search for words that contain the search term
// This is useful when exact match is not found
func FuzzySearch(searchTerm string, limit int) ([]*WordEntry, error) {
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, word, meaning, synonyms, conjugation, created_at, updated_at
		FROM words
		WHERE word LIKE ? COLLATE NOCASE
		LIMIT ?
	`

	rows, err := db.Query(query, "%"+searchTerm+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*WordEntry
	for rows.Next() {
		var entry WordEntry
		var synonyms, conjugation sql.NullString

		err := rows.Scan(
			&entry.ID,
			&entry.Word,
			&entry.Meaning,
			&synonyms,
			&conjugation,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if synonyms.Valid {
			entry.Synonyms = &synonyms.String
		}
		if conjugation.Valid {
			entry.Conjugation = &conjugation.String
		}

		entries = append(entries, &entry)
	}

	return entries, rows.Err()
}
