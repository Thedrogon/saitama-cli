// problem.go
package main

import (
	"encoding/json"
	"os"
)

// Problem defines the structure for a coding problem
type Problem struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

const dbFile = "problems.json"

// loadProblems reads the problems from the JSON file
func loadProblems() ([]Problem, error) {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Problem{}, nil // Return empty list if file doesn't exist
		}
		return nil, err
	}

	var problems []Problem
	err = json.Unmarshal(data, &problems)
	return problems, err
}

// saveProblems writes the current list of problems to the JSON file
func saveProblems(problems []Problem) error {
	data, err := json.MarshalIndent(problems, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dbFile, data, 0644)
}