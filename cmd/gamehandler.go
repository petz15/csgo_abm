package main

import (
	"CSGO_ABM/internal/engine"
	"CSGO_ABM/output"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// StartGame initializes and runs a single game simulation
func StartGame(team1Name string, team1Strategy string, team2Name string, team2Strategy string,
	gameRules string, simPrefix string) string {

	ID := createGameID()
	if simPrefix != "" {
		ID = simPrefix + ID
	}

	// Create a new game instance
	game := engine.NewGame(ID, team1Name, team1Strategy, team2Name, team2Strategy, gameRules)

	// Start the simulation
	game.Start()

	// Create results directory if it doesn't exist
	resultsDir := "results"
	os.MkdirAll(resultsDir, 0755)

	// Export the results to JSON
	resultsPath := filepath.Join(resultsDir, ID+".json")
	err := output.ExportResultsToJSON(game, resultsPath)
	if err != nil {
		fmt.Printf("Error exporting results: %v\n", err)
	}

	return ID
}

// createGameID creates a unique identifier for a game session
func createGameID() string {
	// Get the hostname
	hostname := "host"
	name, err := os.Hostname()
	if err == nil {
		hostname = name
	}

	// Generate random component
	uniqueID := generateRandomString(4)

	// Format: YYYYMMDD_HHMMSS_hostname_randomID
	ID := time.Now().Format("20060102_150405") + "_" + hostname + "_" + uniqueID

	return ID
}

// generateRandomString generates a random string of specified length
func generateRandomString(max int) string {
	table := [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
