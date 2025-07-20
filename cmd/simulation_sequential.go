package main

import (
	"CSGO_ABM/internal/analysis"
	"CSGO_ABM/util"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func sequentialsimulation(config analysis.SimulationConfig) error {
	// Sequential simulation logic - runs simulations one after another without parallelism

	starttime := time.Now()
	stats := analysis.NewStats(config.NumSimulations, "sequential")

	// Create results directory if exporting individual results
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

		// Simulate a single game
		result, err := StartGameWithResultsAndExport(config.Team1Name, config.Team1Strategy, config.Team2Name,
			config.Team2Strategy, config.GameRules, simPrefix, config.ExportDetailedResults, config.Exportpath)
		if err != nil {
			fmt.Printf("Simulation %d failed: %v\n", i+1, err)
			continue
		}

		// Export individual result if requested
		if config.ExportDetailedResults && result.GameID != "" {
			filename := fmt.Sprintf("%s/%s.json", config.Exportpath, result.GameID)
			if resultData, err := json.MarshalIndent(result, "", "  "); err == nil {
				os.WriteFile(filename, resultData, 0644)
			}
		}

		updateglobalstats(stats, result)

		// Progress update for long runs
		if config.NumSimulations > 100 && (i+1)%100 == 0 {
			fmt.Printf("Progress: %d/%d simulations completed\n", i+1, config.NumSimulations)
		}
	}

	endtime := time.Now()
	stats.ExecutionTime = endtime.Sub(starttime)

	// Show stats
	showstats(stats)

	// Export summary
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
