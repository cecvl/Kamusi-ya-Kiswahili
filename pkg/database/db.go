package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db   *sql.DB
	once sync.Once
)

// getDBPath determines the database file path in the following order:
// 1. Environment variable KAMUSI_DB_PATH
// 2. ./data/kamusi.db (current directory)
// 3. Executable directory/data/kamusi.db
// 4. ~/.kamusi/kamusi.db (user home directory)
func getDBPath() (string, error) {
	// 1. Check environment variable
	if dbPath := os.Getenv("KAMUSI_DB_PATH"); dbPath != "" {
		return dbPath, nil
	}

	// 2. Check current directory
	localPath := "data/kamusi.db"
	if _, err := os.Stat(localPath); err == nil {
		return localPath, nil
	}

	// 3. Check executable directory
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		exeDBPath := filepath.Join(exeDir, "data", "kamusi.db")
		if _, err := os.Stat(exeDBPath); err == nil {
			return exeDBPath, nil
		}
	}

	// 4. Use user home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %v", err)
	}

	kamuiDir := filepath.Join(homeDir, ".kamusi")
	dbPath := filepath.Join(kamuiDir, "kamusi.db")

	// Create directory if it doesn't exist
	if err := os.MkdirAll(kamuiDir, 0755); err != nil {
		return "", fmt.Errorf("could not create kamusi directory: %v", err)
	}

	return dbPath, nil
}

// GetDB returns a singleton database connection with connection pooling
func GetDB() (*sql.DB, error) {
	var err error
	once.Do(func() {
		var dbPath string
		dbPath, err = getDBPath()
		if err != nil {
			return
		}

		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			return
		}

		// Configure connection pool for better concurrency
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)

		err = db.Ping()
	})
	return db, err
}

// InitDB creates the database schema
func InitDB() error {
	database, err := GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database: %v", err)
	}

	// Create words table
	wordsTable := `
	CREATE TABLE IF NOT EXISTS words (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word TEXT NOT NULL,
		meaning TEXT NOT NULL,
		synonyms TEXT,
		conjugation TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	-- Index for faster word lookups
	CREATE INDEX IF NOT EXISTS idx_word ON words(word COLLATE NOCASE);
	`

	// Create missing_words table to track words that couldn't be found
	missingWordsTable := `
	CREATE TABLE IF NOT EXISTS missing_words (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word TEXT NOT NULL,
		search_count INTEGER DEFAULT 1,
		first_searched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_searched_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(word COLLATE NOCASE)
	);
	
	-- Index for faster lookups
	CREATE INDEX IF NOT EXISTS idx_missing_word ON missing_words(word COLLATE NOCASE);
	`

	// Execute schema creation
	if _, err := database.Exec(wordsTable); err != nil {
		return fmt.Errorf("failed to create words table: %v", err)
	}

	if _, err := database.Exec(missingWordsTable); err != nil {
		return fmt.Errorf("failed to create missing_words table: %v", err)
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
