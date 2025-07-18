package main

import (
	"CSGO_ABM/internal/engine"
	"CSGO_ABM/output"
	"crypto/rand"
	"io"
	"os"
	"time"
)

func main() {
	// Initialize the game engine

	//for test cases
	Team1Name := "Team A"
	Team1Strategy := "all_in"
	Team2Name := "Team B"
	Team2Strategy := "default_half"
	gamerules := "default"

	ID := generateID()

	// Create a new game instance
	game := engine.NewGame(ID, Team1Name, Team1Strategy, Team2Name, Team2Strategy, gamerules)

	// Start the simulation
	game.Start()

	// Export the results to JSON
	err := output.ExportResultsToJSON(game, ID+".json")
	if err != nil {
		println("Error exporting results:", err)
	}
}

func generateID() string {
	// Get the hostname of the server/desktop
	hostname := "host" // Replace with actual hostname retrieval logic if needed

	name, err := os.Hostname()
	if err != nil {
		hostname = name
	}

	//uniqueID
	uniqueID := EncodeToString(4)

	// Generate a unique ID using the current time and hostname and random number
	// This ensures that the ID is unique and can be used to identify the game session
	// The format is YYYYMMDD_HHMMSS_hostname

	ID := time.Now().Format("20060102_150405") + "_" + hostname + "_" + uniqueID

	return ID
}

func EncodeToString(max int) string {
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

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
