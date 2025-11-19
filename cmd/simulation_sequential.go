package main

import (
	"csgo_abm/internal/analysis"
	"csgo_abm/internal/engine"
	"csgo_abm/util"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func updateglobalstats(stats *analysis.SimulationStats, result *GameResult) {
	// Use the unified analysis package method
	stats.UpdateGameResult(
		result.Team1Won,
		result.Team1Score,
		result.Team2Score,
		result.TotalRounds,
		result.WentToOvertime,
		0, // responseTime - not tracked in sequential mode
	)

	// Convert and update economic statistics
	team1Econ := analysis.TeamGameEconomics{
		TotalSpent:       result.Team1Economics.TotalSpent,
		TotalEarned:      result.Team1Economics.TotalEarned,
		AverageFunds:     result.Team1Economics.AverageFunds,
		AverageRSEq:      result.Team1Economics.AverageRSEq,
		AverageFTEEq:     result.Team1Economics.AverageFTEEq,
		AverageREEq:      result.Team1Economics.AverageREEq,
		AverageSurvivors: result.Team1Economics.AverageSurvivors,
		MaxFunds:         result.Team1Economics.MaxFunds,
		MinFunds:         result.Team1Economics.MinFunds,
		MaxConsecLosses:  result.Team1Economics.MaxConsecLosses,
	}
	team2Econ := analysis.TeamGameEconomics{
		TotalSpent:       result.Team2Economics.TotalSpent,
		TotalEarned:      result.Team2Economics.TotalEarned,
		AverageFunds:     result.Team2Economics.AverageFunds,
		AverageRSEq:      result.Team2Economics.AverageRSEq,
		AverageFTEEq:     result.Team2Economics.AverageFTEEq,
		AverageREEq:      result.Team2Economics.AverageREEq,
		AverageSurvivors: result.Team2Economics.AverageSurvivors,
		MaxFunds:         result.Team2Economics.MaxFunds,
		MinFunds:         result.Team2Economics.MinFunds,
		MaxConsecLosses:  result.Team2Economics.MaxConsecLosses,
	}
	stats.UpdateEconomicStats(team1Econ, team2Econ, result.TotalRounds)
}

func showstats(stats *analysis.SimulationStats) {
	// Use the unified analysis package for enhanced reporting
	analysis.PrintEnhancedStats(stats)
}

func exportSummary_v2(stats *analysis.SimulationStats, pathOrID string) error {
	var filename string

	// Check if pathOrID contains a path separator (indicates it's a full path)
	if strings.Contains(pathOrID, "/") || strings.Contains(pathOrID, "\\") {
		filename = pathOrID
	} else {
		// Original behavior - create filename from ID
		filename = fmt.Sprintf("%s_results.json", pathOrID)
	}

	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// sequentialsimulationWithRules is an optimized version that uses pre-validated GameRules
func sequentialsimulation(config SimulationConfig, gameRules engine.GameRules) error {
	starttime := time.Now()
	stats := analysis.NewStats(config.NumSimulations, "sequential")

	// Display simulation mode and export information
	if config.ExportDetailedResults {
		fmt.Printf("Sequential simulation with individual result export enabled\n")
		fmt.Printf("Results will be saved to: %s/\n", config.Exportpath)
		if config.NumSimulations > 10000 {
			fmt.Printf("WARNING: Exporting %d individual results may create filesystem pressure\n", config.NumSimulations)
		}
	} else {
		fmt.Println("Sequential simulation with summary-only mode")
	}

	fmt.Printf("Starting %d sequential simulations...\n", config.NumSimulations)

	for i := 0; i < config.NumSimulations; i++ {
		// Generate a simulation prefix for this run
		simPrefix := fmt.Sprintf("seq_sim_%d_", i+1)

		// Simulate a single game with pre-validated rules
		result, err := StartGame_default(config.Team1Name, config.Team1Strategy, config.Team2Name,
			config.Team2Strategy, gameRules, simPrefix, config.ExportDetailedResults, false, config.Exportpath)
		if err != nil {
			fmt.Printf("Simulation %d failed: %v\n", i+1, err)
			continue
		}

		// Update statistics with the result
		updateglobalstats(stats, result)

		// Periodic garbage collection every 10% of simulations if over 1000, otherwise every 100
		var gcInterval int
		if config.NumSimulations > 1000 {
			gcInterval = config.NumSimulations / 10
		} else {
			gcInterval = 100
		}
		if (i+1)%gcInterval == 0 {
			runtime.GC()
		}

		// Progress reporting every 10% of simulations if over 1000, otherwise every 100
		var progressInterval int
		if config.NumSimulations > 1000 {
			progressInterval = config.NumSimulations / 10
		} else {
			progressInterval = 100
		}
		if config.NumSimulations > 100 && (i+1)%progressInterval == 0 {
			fmt.Printf("Progress: %d/%d simulations completed\n", i+1, config.NumSimulations)
		}
	}

	// Calculate final execution time
	endtime := time.Now()
	stats.ExecutionTime = endtime.Sub(starttime)

	// Generate comprehensive final summary
	showstats(stats)

	// Export summary statistics
	if config.ExportDetailedResults {
		summaryPath := fmt.Sprintf("%s/simulation_summary.json", config.Exportpath)
		exportSummary_v2(stats, summaryPath)
		fmt.Printf("\nResults exported to: %s/\n", config.Exportpath)
		fmt.Printf("- Individual game results: %d JSON files\n", stats.CompletedSims)
		fmt.Println("- Summary statistics: simulation_summary.json")
	} else {
		// Use original export method for single file
		simID := util.CreateGameID()
		exportSummary_v2(stats, simID)
	}

	return nil
}
