package main

import (
	"dbg_abm/internal/analysis"
	"dbg_abm/internal/engine"
	"dbg_abm/util"
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

	// Advanced analysis removed

	// Storage for games if we need combined CSV export (modes 2 or 4)
	var allGames []*engine.Game
	if config.CSVExportMode == 2 || config.CSVExportMode == 4 {
		allGames = make([]*engine.Game, 0, config.NumSimulations)
	}

	// Display simulation mode and export information
	if !config.SuppressOutput {
		if config.ExportDetailedResults {
			fmt.Printf("Sequential simulation with individual result export enabled\n")
			fmt.Printf("Results will be saved to: %s/\n", config.Exportpath)
			if config.NumSimulations > 10000 {
				fmt.Printf("WARNING: Exporting %d individual results may create filesystem pressure\n", config.NumSimulations)
			}
		} else {
			fmt.Println("Sequential simulation with summary-only mode")
		}

		if config.CSVExportMode > 0 {
			fmt.Printf("CSV export mode: %d\n", config.CSVExportMode)
		}

		fmt.Printf("Starting %d sequential simulations...\n", config.NumSimulations)
	}

	for i := 0; i < config.NumSimulations; i++ {
		// Generate a simulation prefix for this run
		simPrefix := fmt.Sprintf("seq_sim_%d_", i+1)

		// Simulate a single game with pre-validated rules
		result, err := StartGame_default(config.Team1Name, config.Team1Strategy, config.Team2Name,
			config.Team2Strategy, gameRules, simPrefix, config.ExportDetailedResults, false, config.CSVExportMode, config.Exportpath)
		if err != nil {
			if !config.SuppressOutput {
				fmt.Printf("Simulation %d failed: %v\n", i+1, err)
			}
			continue
		}

		// Store game data for combined CSV export
		if config.CSVExportMode == 2 || config.CSVExportMode == 4 {
			if result.GameData != nil {
				allGames = append(allGames, result.GameData)
			}
		}

		// Update statistics with the result
		updateglobalstats(stats, result)

		// Advanced analysis removed

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
			if !config.SuppressOutput {
				fmt.Printf("Progress: %d/%d simulations completed\n", i+1, config.NumSimulations)
			}
		}
	}

	// Calculate final execution time
	endtime := time.Now()
	stats.ExecutionTime = endtime.Sub(starttime)

	// Export combined CSV if mode 2 or 4
	if config.CSVExportMode == 2 && len(allGames) > 0 {
		if !config.SuppressOutput {
			fmt.Println("\nExporting combined full CSV...")
		}
		csvPath := fmt.Sprintf("%s/all_games_full.csv", config.Exportpath)
		err := util.ExportAllGamesAllDataCSV(allGames, csvPath)
		if err != nil {
			if !config.SuppressOutput {
				fmt.Printf("Warning: Error exporting combined full CSV: %v\n", err)
			}
		} else {
			if !config.SuppressOutput {
				fmt.Printf("✅ Combined full CSV exported: %s\n", csvPath)
			}
		}
	} else if config.CSVExportMode == 4 && len(allGames) > 0 {
		if !config.SuppressOutput {
			fmt.Println("\nExporting combined minimal CSV...")
		}
		csvPath := fmt.Sprintf("%s/all_games_minimal.csv", config.Exportpath)
		err := util.ExportAllGamesMinimalCSV(allGames, csvPath)
		if err != nil {
			if !config.SuppressOutput {
				fmt.Printf("Warning: Error exporting combined minimal CSV: %v\n", err)
			}
		} else {
			if !config.SuppressOutput {
				fmt.Printf("✅ Combined minimal CSV exported: %s\n", csvPath)
			}
		}
	}

	// Generate comprehensive final summary
	if !config.SuppressOutput {
		showstats(stats)
	}

	// Advanced analysis removed

	// Export summary statistics (always export to simulation_summary.json)
	summaryPath := fmt.Sprintf("%s/simulation_summary.json", config.Exportpath)
	if err := exportSummary_v2(stats, summaryPath); err != nil {
		if !config.SuppressOutput {
			fmt.Printf("Warning: Failed to export summary: %v\n", err)
		}
	}

	if !config.SuppressOutput {
		fmt.Printf("\nResults exported to: %s/\n", config.Exportpath)
		if config.ExportDetailedResults {
			fmt.Printf("- Individual game results: %d JSON files\n", stats.CompletedSims)
		}
		fmt.Println("- Summary statistics: simulation_summary.json")
	}

	return nil
}
