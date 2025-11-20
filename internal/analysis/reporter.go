package analysis

import (
	"fmt"
	"strings"
	"time"
)

// PrintEnhancedStats prints comprehensive statistics with enhanced analysis
func PrintEnhancedStats(stats *SimulationStats) {
	printBasicStats(stats)
	printGameAnalysis(stats)
}

// printBasicStats prints core simulation statistics
func printBasicStats(stats *SimulationStats) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎮 SIMULATION RESULTS SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Simulation Mode: %s\n", strings.Title(stats.SimulationMode))
	fmt.Printf("Total Simulations: %d (Completed: %d, Failed: %d)\n",
		stats.TotalSimulations, stats.CompletedSims, stats.FailedSims)
	fmt.Printf("Execution Time: %s\n", stats.ExecutionTime.Round(time.Second))
	if stats.ProcessingRate > 0 {
		fmt.Printf("Processing Rate: %.1f simulations/second\n", stats.ProcessingRate)
	}
	fmt.Println()
}

// printGameAnalysis prints game-specific analysis
func printGameAnalysis(stats *SimulationStats) {
	if stats.CompletedSims == 0 {
		return
	}

	fmt.Println("📊 GAME ANALYSIS")
	fmt.Println(strings.Repeat("-", 40))

	// Team performance
	fmt.Printf("Team 1 Wins: %d (%.1f%%)\n", stats.Team1Wins, stats.Team1WinRate)
	fmt.Printf("Team 2 Wins: %d (%.1f%%)\n", stats.Team2Wins, stats.Team2WinRate)

	// Game characteristics
	fmt.Printf("Average Rounds per Game: %.1f\n", stats.AverageRounds)
	fmt.Printf("Overtime Rate: %.1f%%\n", stats.OvertimeRate)

	fmt.Println()
}

// Advanced analysis removed

// PrintProgressSummary prints a compact progress update (for live monitoring)
func PrintProgressSummary(stats *SimulationStats) {
	if stats.CompletedSims == 0 {
		return
	}

	team1Rate := float64(stats.Team1Wins) / float64(stats.CompletedSims) * 100
	progressPercent := float64(stats.CompletedSims) / float64(stats.TotalSimulations) * 100

	fmt.Printf("\r🎮 Progress: %d/%d (%.1f%%) | Team1: %.1f%% | Team2: %.1f%% | Rate: %.0f/sec",
		stats.CompletedSims, stats.TotalSimulations, progressPercent,
		team1Rate, 100-team1Rate, stats.ProcessingRate)
}

// ClearProgressLine clears the progress line
func ClearProgressLine() {
	fmt.Print("\r" + strings.Repeat(" ", 80) + "\r")
}
