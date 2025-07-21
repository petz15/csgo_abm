package main

import (
	"CSGO_ABM/internal/engine"
	"CSGO_ABM/util"
	"fmt"
	"os"
	"path/filepath"
)

// GameResult holds the essential results from a game simulation
type GameResult struct {
	GameID         string
	Team1Won       bool
	Team1Score     int
	Team2Score     int
	TotalRounds    int
	WentToOvertime bool
}

// StartGame initializes and runs a single game simulation (legacy function for single simulations)
func StartGame_simple(team1Name string, team1Strategy string, team2Name string, team2Strategy string,
	gameRules engine.GameRules, simPrefix string) string {

	ID := util.CreateGameID()
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
	err := util.ExportResultsToJSON(game, resultsPath)
	if err != nil {
		fmt.Printf("Error exporting results: %v\n", err)
	}

	// Explicitly clear game to help with garbage collection
	game = nil

	return ID
}

// StartGameWithValidatedRules runs a simulation with pre-validated GameRules (optimized for batch processing)
func StartGame_default(team1Name string, team1Strategy string, team2Name string, team2Strategy string,
	gameRules engine.GameRules, simPrefix string, exportJSON bool, exportpath string) (*GameResult, error) {

	ID := util.CreateGameID()
	if simPrefix != "" {
		ID = simPrefix + ID
	}

	// Create a new game instance with pre-validated rules
	game := engine.NewGame(ID, team1Name, team1Strategy, team2Name, team2Strategy, gameRules)

	// Start the simulation
	game.Start()

	// Extract results directly from the game object
	result := &GameResult{
		GameID:         ID,
		Team1Won:       !game.WinnerTeam, // WinnerTeam is false if Team1 wins
		Team1Score:     game.Score[0],
		Team2Score:     game.Score[1],
		TotalRounds:    len(game.Rounds),
		WentToOvertime: game.OT,
	}

	// Optionally export to JSON for debugging/analysis
	if exportJSON {
		resultsDir := exportpath
		os.MkdirAll(resultsDir, 0755)
		resultsPath := filepath.Join(resultsDir, ID+".json")
		err := util.ExportResultsToJSON(game, resultsPath)
		if err != nil {
			fmt.Printf("Warning: Error exporting detailed results for %s: %v\n", ID, err)
		}
	}

	// Explicitly clear game to help with garbage collection
	game = nil

	return result, nil
}
