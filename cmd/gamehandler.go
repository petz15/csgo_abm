package main

import (
	"csgo_abm/internal/engine"
	"csgo_abm/util"
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
	Team1Economics TeamGameEconomics
	Team2Economics TeamGameEconomics
	GameData       *engine.Game // Optional: full game data for advanced analysis
}

// TeamGameEconomics holds economic statistics for a single game
type TeamGameEconomics struct {
	TotalSpent       float64
	TotalEarned      float64
	AverageFunds     float64 // Average funds at round start
	AverageRSEq      float64
	AverageFTEEq     float64
	AverageREEq      float64
	AverageSurvivors float64
	MaxFunds         float64
	MinFunds         float64
	MaxConsecLosses  int
}

// StartGameWithValidatedRules runs a simulation with pre-validated GameRules (optimized for batch processing)
func StartGame_default(team1Name string, team1Strategy string, team2Name string, team2Strategy string,
	gameRules engine.GameRules, simPrefix string, exportJSON bool, exportRounds bool, exportpath string) (*GameResult, error) {

	ID := util.CreateGameID()
	if simPrefix != "" {
		ID = simPrefix + ID
	}

	// Create a new game instance with pre-validated rules
	game := engine.NewGame(ID, team1Name, team1Strategy, team2Name, team2Strategy, gameRules)

	// Start the simulation
	game.Start()

	// Calculate economic statistics for both teams
	team1Econ := calculateTeamEconomics(game.Team1, len(game.Rounds))
	team2Econ := calculateTeamEconomics(game.Team2, len(game.Rounds))

	// Extract results directly from the game object
	result := &GameResult{
		GameID:         ID,
		Team1Won:       game.Is_T1_Winner,
		Team1Score:     game.Score[0],
		Team2Score:     game.Score[1],
		TotalRounds:    len(game.Rounds),
		WentToOvertime: game.OT,
		Team1Economics: team1Econ,
		Team2Economics: team2Econ,
		GameData:       game, // Store game for advanced analysis
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

	// Optionally export round-by-round data
	if exportRounds {
		resultsDir := exportpath
		os.MkdirAll(resultsDir, 0755)

		// Export CSV (most compact, ~10x smaller than JSON)
		roundsPathCSV := filepath.Join(resultsDir, ID+"_rounds.csv")
		fmt.Printf("📊 Exporting round data to CSV: %s\n", roundsPathCSV)
		err := util.ExportRoundsToCSV(game, roundsPathCSV)
		if err != nil {
			fmt.Printf("Warning: Error exporting CSV rounds for %s: %v\n", ID, err)
		} else {
			fmt.Printf("✅ CSV round data exported successfully\n")
		}

		// Export simple JSON version (compact with key metrics only)
		roundsPathSimple := filepath.Join(resultsDir, ID+"_rounds_simple.json")
		fmt.Printf("📊 Exporting simplified JSON round data to: %s\n", roundsPathSimple)
		err = util.ExportRoundsToJSONSimple(game, roundsPathSimple)
		if err != nil {
			fmt.Printf("Warning: Error exporting simple rounds for %s: %v\n", ID, err)
		} else {
			fmt.Printf("✅ Simplified JSON round data exported successfully\n")
		}

		// Also export full version with all details
		roundsPath := filepath.Join(resultsDir, ID+"_rounds_full.json")
		fmt.Printf("📊 Exporting full round data to: %s\n", roundsPath)
		err = util.ExportRoundsToJSON(game, roundsPath)
		if err != nil {
			fmt.Printf("Warning: Error exporting full rounds for %s: %v\n", ID, err)
		} else {
			fmt.Printf("✅ Full round data exported successfully\n")
		}
	}

	// Don't cleanup game yet if it might be used for advanced analysis
	// Caller is responsible for cleanup after using GameData
	if result.GameData == nil {
		// Cleanup game data to free memory
		game.Cleanup()

		// Explicitly clear game data structures to help with garbage collection
		game.Rounds = nil
		game.Team1 = nil
		game.Team2 = nil
		game = nil
	}

	return result, nil
}

// calculateTeamEconomics computes economic statistics for a team from its round data
func calculateTeamEconomics(team *engine.Team, totalRounds int) TeamGameEconomics {
	if len(team.RoundData) == 0 || totalRounds == 0 {
		return TeamGameEconomics{}
	}

	var totalSpent, totalEarned, totalFunds, totalRSEq, totalFTEEq, totalREEq float64
	var totalSurvivors int
	maxFunds := 0.0
	minFunds := team.RoundData[0].Funds
	maxConsecLosses := 0
	currentConsecLosses := 0

	for _, rd := range team.RoundData {
		totalSpent += rd.Spent
		totalEarned += rd.Earned
		totalFunds += rd.Funds
		totalRSEq += rd.RS_Eq_value
		totalFTEEq += rd.FTE_Eq_value
		totalREEq += rd.RE_Eq_value
		totalSurvivors += rd.Survivors

		if rd.Funds > maxFunds {
			maxFunds = rd.Funds
		}
		if rd.Funds < minFunds {
			minFunds = rd.Funds
		}

		// Track consecutive losses
		if rd.Consecutiveloss > currentConsecLosses {
			currentConsecLosses = rd.Consecutiveloss
			if currentConsecLosses > maxConsecLosses {
				maxConsecLosses = currentConsecLosses
			}
		} else if rd.Consecutiveloss == 0 {
			currentConsecLosses = 0
		}
	}

	roundCount := float64(len(team.RoundData))

	return TeamGameEconomics{
		TotalSpent:       totalSpent,
		TotalEarned:      totalEarned,
		AverageFunds:     totalFunds / roundCount,
		AverageRSEq:      totalRSEq / roundCount,
		AverageFTEEq:     totalFTEEq / roundCount,
		AverageREEq:      totalREEq / roundCount,
		AverageSurvivors: float64(totalSurvivors) / roundCount,
		MaxFunds:         maxFunds,
		MinFunds:         minFunds,
		MaxConsecLosses:  maxConsecLosses,
	}
}
