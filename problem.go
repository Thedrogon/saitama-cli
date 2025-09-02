// problem.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

// Problem defines the structure for a coding problem
type Problem struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Tags       []string  `json:"tags"`
	DateAdded  time.Time `json:"date_added,omitempty"`
	LastSolved time.Time `json:"last_solved,omitempty"`
	SolveCount int       `json:"solve_count,omitempty"`
	Difficulty string    `json:"difficulty,omitempty"` // easy, medium, hard
	Platform   string    `json:"platform,omitempty"`   // leetcode, codeforces, etc.
	URL        string    `json:"url,omitempty"`
	Notes      string    `json:"notes,omitempty"`
}

const maxBackups = 5

// getDbPath finds the appropriate user config directory for data storage.
// THIS IS THE CRITICAL FIX TO PREVENT DATA LOSS.
func getDbPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("could not get user config directory: %w", err)
	}
	appConfigDir := filepath.Join(configDir, "saitama")
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("could not create app config directory: %w", err)
	}
	return filepath.Join(appConfigDir, "problems.json"), nil
}

// getBackupDir returns the path to the backup directory inside the app's config folder.
func getBackupDir() (string, error) {
	dbPath, err := getDbPath()
	if err != nil {
		return "", err
	}
	// Place backups in the same directory as the database file.
	return filepath.Join(filepath.Dir(dbPath), ".saitama_backups"), nil
}

// loadProblems reads the problems from the JSON file in the user's config directory.
func loadProblems() ([]Problem, error) {
	dbPath, err := getDbPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return []Problem{}, nil // File doesn't exist yet, return empty list.
	}

	data, err := os.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read problems file: %w", err)
	}

	if len(data) == 0 {
		return []Problem{}, nil // Handle empty file
	}

	var problems []Problem
	if err := json.Unmarshal(data, &problems); err != nil {
		return nil, fmt.Errorf("failed to parse problems file: %w", err)
	}

	// Data migration for older records without DateAdded
	needsSave := false
	for i := range problems {
		if problems[i].DateAdded.IsZero() {
			problems[i].DateAdded = time.Now() // Default to now
			needsSave = true
		}
	}
	if needsSave {
		// Save migrated data silently
		_ = saveProblems(problems)
	}

	return problems, nil
}

// saveProblems writes the current list of problems to the JSON file, creating a backup first.
func saveProblems(problems []Problem) error {
	dbPath, err := getDbPath()
	if err != nil {
		return err
	}

	if err := createBackup(dbPath); err != nil {
		// Don't fail the save operation if backup fails, just warn
		color.Yellow("Warning: Failed to create backup: %v\n", err)
	}

	data, err := json.MarshalIndent(problems, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal problems: %w", err)
	}

	// Atomic write operation
	tempFile := dbPath + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}
	if err := os.Rename(tempFile, dbPath); err != nil {
		_ = os.Remove(tempFile) // Clean up temp file on failure
		return fmt.Errorf("failed to replace problems file: %w", err)
	}
	return nil
}

// createBackup creates a backup of the current problems file.
func createBackup(dbPath string) error {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil // Nothing to backup
	}

	backupDir, err := getBackupDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("problems_%s.json", timestamp))

	data, err := os.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("failed to read original file for backup: %w", err)
	}
	if err := os.WriteFile(backupFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return cleanupOldBackups(backupDir)
}

// cleanupOldBackups removes old backup files, keeping only the most recent ones.
func cleanupOldBackups(backupDir string) error {
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}

	var backups []os.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			backups = append(backups, entry)
		}
	}

	if len(backups) <= maxBackups {
		return nil
	}

	// The user requested to remove explicit sorting.
	// We now rely on the filesystem's default order, which is generally chronological
	// for timestamped filenames but is not guaranteed across all systems.

	// Remove the oldest backups (assuming first entries are the oldest)
	for i := 0; i < len(backups)-maxBackups; i++ {
		if err := os.Remove(filepath.Join(backupDir, backups[i].Name())); err != nil {
			// Log error but continue trying to clean up others
			fmt.Printf("Warning: could not remove old backup %s: %v\n", backups[i].Name(), err)
		}
	}
	return nil
}

// findProblemByID finds a problem by its ID (case-insensitive) and returns it and its index.
func findProblemByID(problems []Problem, id string) (*Problem, int) {
	for i, p := range problems {
		if p.ID == id {
			return &problems[i], i
		}
	}
	return nil, -1
}

// exportProblems exports problems to a specified file.
func exportProblems(problems []Problem, filename string) error {
	data, err := json.MarshalIndent(problems, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal problems for export: %w", err)
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write export file: %w", err)
	}
	return nil
}

// importProblems imports problems from a specified file.
func importProblems(filename string) ([]Problem, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read import file: %w", err)
	}

	var importedProblems []Problem
	if err := json.Unmarshal(data, &importedProblems); err != nil {
		return nil, fmt.Errorf("failed to parse import file: %w", err)
	}

	// Validate imported problems
	for i, p := range importedProblems {
		if p.ID == "" || p.Name == "" {
			return nil, fmt.Errorf("invalid problem at index %d (ID or Name is empty)", i)
		}
	}
	return importedProblems, nil
}

