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
	printPerformanceMetrics(stats)
	printDistributions(stats)
	printAdvancedAnalysis(stats)
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

	if stats.AdvancedStats != nil {
		// Balance analysis
		fmt.Printf("Strategy Balance Score: %.1f%%", stats.AdvancedStats.BalanceScore)
		if stats.AdvancedStats.BalanceScore > 90 {
			fmt.Print(" (Excellent)")
		} else if stats.AdvancedStats.BalanceScore > 80 {
			fmt.Print(" (Good)")
		} else if stats.AdvancedStats.BalanceScore > 70 {
			fmt.Print(" (Fair)")
		} else {
			fmt.Print(" (Poor)")
		}
		fmt.Println()

		// Statistical significance
		fmt.Printf("Statistical Significance: %s", stats.AdvancedStats.StatisticalSignificance)
		if stats.AdvancedStats.StatisticalSignificance == "significant" {
			fmt.Printf(" (χ² = %.2f)", stats.AdvancedStats.ChiSquareValue)
		}
		fmt.Println()

		// Game length analysis
		if stats.AverageRounds > 26 {
			fmt.Printf("Game Length: Long (avg %.1f rounds) - Strategies well-matched\n", stats.AverageRounds)
		} else if stats.AverageRounds < 22 {
			fmt.Printf("Game Length: Short (avg %.1f rounds) - One strategy dominates\n", stats.AverageRounds)
		} else {
			fmt.Printf("Game Length: Normal (avg %.1f rounds)\n", stats.AverageRounds)
		}

		// Competitiveness analysis
		closeGameRate := float64(stats.AdvancedStats.CloseGames) / float64(stats.CompletedSims) * 100
		blowoutRate := float64(stats.AdvancedStats.BlowoutGames) / float64(stats.CompletedSims) * 100

		fmt.Printf("Close Games (≤3 round diff): %d (%.1f%%)\n", stats.AdvancedStats.CloseGames, closeGameRate)
		fmt.Printf("Blowout Games (>10 round diff): %d (%.1f%%)\n", stats.AdvancedStats.BlowoutGames, blowoutRate)

		if closeGameRate > 30 {
			fmt.Println("Competitiveness: High (many close games)")
		} else if closeGameRate < 15 {
			fmt.Println("Competitiveness: Low (few close games)")
		} else {
			fmt.Println("Competitiveness: Moderate")
		}
	}
	fmt.Println()
}

// printPerformanceMetrics prints performance-related statistics
func printPerformanceMetrics(stats *SimulationStats) {
	if stats.SimulationMode != "concurrent" {
		return
	}

	fmt.Println("⚡ PERFORMANCE METRICS")
	fmt.Println(strings.Repeat("-", 40))
	if stats.PeakMemoryUsage > 0 {
		fmt.Printf("Peak Memory Usage: %d MB\n", stats.PeakMemoryUsage)
		if stats.CompletedSims > 0 {
			memEfficiency := float64(stats.CompletedSims) / float64(stats.PeakMemoryUsage)
			fmt.Printf("Memory Efficiency: %.1f simulations/MB\n", memEfficiency)
		}
	}
	if stats.TotalGCRuns > 0 {
		fmt.Printf("Total GC Runs: %d\n", stats.TotalGCRuns)
	}

	if stats.AdvancedStats != nil && len(stats.AdvancedStats.ResponseTimes) > 0 {
		fmt.Printf("Response Time P50: %v\n", stats.AdvancedStats.P50ResponseTime.Round(time.Millisecond))
		fmt.Printf("Response Time P95: %v\n", stats.AdvancedStats.P95ResponseTime.Round(time.Millisecond))
		fmt.Printf("Response Time P99: %v\n", stats.AdvancedStats.P99ResponseTime.Round(time.Millisecond))
	}
	fmt.Println()
}

// printDistributions prints score and round distributions
func printDistributions(stats *SimulationStats) {
	if stats.CompletedSims == 0 {
		return
	}

	// Score distribution
	if len(stats.ScoreDistribution) > 0 && stats.AdvancedStats != nil {
		fmt.Println("📈 SCORE DISTRIBUTIONS")
		fmt.Println(strings.Repeat("-", 40))

		count := 0
		maxDisplay := 8
		for _, scoreLine := range stats.AdvancedStats.TopScoreLines {
			if count >= maxDisplay {
				break
			}
			fmt.Printf("  %s: %d games (%.1f%%)\n",
				scoreLine.Score, scoreLine.Count, scoreLine.Frequency)
			count++
		}
		if len(stats.AdvancedStats.TopScoreLines) > maxDisplay {
			fmt.Printf("  ... and %d more score combinations\n",
				len(stats.ScoreDistribution)-maxDisplay)
		}
		fmt.Println()
	}

	// Round distribution highlights
	if len(stats.RoundDistribution) > 0 && stats.AdvancedStats != nil {
		fmt.Println("🎯 ROUND ANALYSIS")
		fmt.Println(strings.Repeat("-", 40))

		if stats.AdvancedStats.MedianRounds > 0 {
			fmt.Printf("Median Rounds: %.1f\n", stats.AdvancedStats.MedianRounds)
			fmt.Printf("Standard Deviation: %.1f\n", stats.AdvancedStats.StdDevRounds)
		}

		// Find most common round counts
		var mostCommonRounds []struct {
			rounds int
			count  int64
		}
		for rounds, count := range stats.RoundDistribution {
			mostCommonRounds = append(mostCommonRounds, struct {
				rounds int
				count  int64
			}{rounds, count})
		}

		// Sort by count (descending)
		for i := 0; i < len(mostCommonRounds)-1; i++ {
			for j := i + 1; j < len(mostCommonRounds); j++ {
				if mostCommonRounds[i].count < mostCommonRounds[j].count {
					mostCommonRounds[i], mostCommonRounds[j] = mostCommonRounds[j], mostCommonRounds[i]
				}
			}
		}

		fmt.Println("Most Common Game Lengths:")
		displayCount := 5
		if len(mostCommonRounds) < displayCount {
			displayCount = len(mostCommonRounds)
		}
		for i := 0; i < displayCount; i++ {
			rounds := mostCommonRounds[i].rounds
			count := mostCommonRounds[i].count
			percentage := float64(count) / float64(stats.CompletedSims) * 100
			fmt.Printf("  %d rounds: %d games (%.1f%%)\n", rounds, count, percentage)
		}
		fmt.Println()
	}
}

// printAdvancedAnalysis prints additional insights
func printAdvancedAnalysis(stats *SimulationStats) {
	if stats.AdvancedStats == nil || stats.CompletedSims == 0 {
		return
	}

	fmt.Println("🔍 INSIGHTS & RECOMMENDATIONS")
	fmt.Println(strings.Repeat("-", 40))

	// Strategy recommendations based on results
	if stats.AdvancedStats.BalanceScore < 70 {
		winnerStrategy := "Team 1"
		if stats.Team2WinRate > stats.Team1WinRate {
			winnerStrategy = "Team 2"
		}
		fmt.Printf("⚠️  Strategy Imbalance Detected: %s strategy appears stronger\n", winnerStrategy)
		fmt.Println("   Consider adjusting parameters or testing different strategies")
	} else if stats.AdvancedStats.BalanceScore > 95 {
		fmt.Println("✅ Excellent Strategy Balance: Strategies are well-matched")
		fmt.Println("   This suggests good game balance for these strategies")
	}

	// Game length insights
	if stats.OvertimeRate > 25 {
		fmt.Println("⏱️  High Overtime Rate: Games are very competitive")
		fmt.Println("   Strategies are closely matched in effectiveness")
	} else if stats.OvertimeRate < 5 {
		fmt.Println("⚡ Low Overtime Rate: Games often decided in regulation")
		fmt.Println("   One strategy may gain consistent early advantages")
	}

	// Performance insights (for concurrent mode)
	if stats.SimulationMode == "concurrent" && stats.ProcessingRate > 0 {
		if stats.ProcessingRate > 200 {
			fmt.Println("🚀 Excellent Performance: High simulation throughput")
		} else if stats.ProcessingRate < 50 {
			fmt.Println("🐌 Consider Performance Tuning: Low simulation throughput")
			fmt.Println("   Try reducing concurrent workers or increasing memory limit")
		}
	}

	// Sample size adequacy
	if stats.CompletedSims < 100 {
		fmt.Println("📊 Small Sample Size: Consider running more simulations for statistical confidence")
	} else if stats.CompletedSims > 10000 {
		fmt.Println("📊 Large Sample Size: Results are statistically robust")
	}

	fmt.Println(strings.Repeat("=", 60))
}

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
