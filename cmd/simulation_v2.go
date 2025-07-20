package main

import (
	"CSGO_ABM/util"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// SimulationConfig holds configuration for parallel simulations
type SimulationConfig_v2 struct {
	NumSimulations int
	Team1Name      string
	Team1Strategy  string
	Team2Name      string
	Team2Strategy  string
	GameRules      string
}

type SimulationStats_v2 struct {
	TotalSimulations  int64            `json:"total_simulations"`
	CompletedSims     int64            `json:"completed_simulations"`
	Team1Wins         int64            `json:"team1_wins"`
	Team2Wins         int64            `json:"team2_wins"`
	TotalRounds       int64            `json:"total_rounds"`
	OvertimeGames     int64            `json:"overtime_games"`
	AverageRounds     float64          `json:"average_rounds"`
	Team1WinRate      float64          `json:"team1_win_rate"`
	Team2WinRate      float64          `json:"team2_win_rate"`
	OvertimeRate      float64          `json:"overtime_rate"`
	ExecutionTime     time.Duration    `json:"execution_time"`
	ScoreDistribution map[string]int64 `json:"score_distribution"`
}

func sequentialsimulation(config SimulationConfig_v2) error {
	// This function is a placeholder for sequential simulation logic
	// It can be used to run simulations one after another without parallelism
	// For now, it does nothing but can be expanded later

	starttime := time.Now()
	stats := SimulationStats_v2{
		TotalSimulations:  int64(config.NumSimulations),
		ScoreDistribution: make(map[string]int64),
	}
	emptyString := ""

	for i := 0; i < config.NumSimulations; i++ {
		// Simulate a single game or task
		// This could involve calling StartGame or similar functions
		result, err := StartGameWithResults(config.Team1Name, config.Team1Strategy, config.Team2Name, config.Team2Strategy, config.GameRules, emptyString)
		if err != nil {
			return fmt.Errorf("failed to start game: %v", err)
		}
		updateglobalstats(&stats, result)
	}

	endtime := time.Now()

	stats.ExecutionTime = endtime.Sub(starttime)

	// Show stats
	simID := util.CreateGameID()
	showstats(&stats)
	exportSummary_v2(&stats, simID+"_results.json")

	return nil
}

func updateglobalstats(stats *SimulationStats_v2, result *GameResult) {
	stats.CompletedSims++
	if result.Team1Won {
		stats.Team1Wins++
	} else {
		stats.Team2Wins++
	}

	stats.TotalRounds += int64(result.TotalRounds)
	if result.WentToOvertime {
		stats.OvertimeGames++
	}

	// Update score distribution
	scoreKey := fmt.Sprintf("%d-%d", result.Team1Score, result.Team2Score)
	stats.ScoreDistribution[scoreKey]++

	// Calculate averages and rates
	stats.AverageRounds = float64(stats.TotalRounds) / float64(stats.CompletedSims)
	stats.Team1WinRate = float64(stats.Team1Wins) / float64(stats.CompletedSims) * 100
	stats.Team2WinRate = float64(stats.Team2Wins) / float64(stats.CompletedSims) * 100
	stats.OvertimeRate = float64(stats.OvertimeGames) / float64(stats.CompletedSims) * 100
}

func showstats(stats *SimulationStats_v2) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("SIMULATION RESULTS SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total Simulations: %d\n", stats.CompletedSims)
	fmt.Printf("Execution Time: %s\n", stats.ExecutionTime.Round(time.Second))
	fmt.Println()
	fmt.Printf("Team 1 Wins: %d (%.1f%%)\n", stats.Team1Wins, stats.Team1WinRate)
	fmt.Printf("Team 2 Wins: %d (%.1f%%)\n", stats.Team2Wins, stats.Team2WinRate)
	fmt.Printf("Average Rounds per Game: %.1f\n", stats.AverageRounds)
	fmt.Printf("Overtime Rate: %.1f%%\n", stats.OvertimeRate)

	if len(stats.ScoreDistribution) > 0 {
		fmt.Println("\nTop Score Distributions:")
		// Print top 5 most common scores
		// This is a simplified display - in production you might want to sort properly
		count := 0
		for score, freq := range stats.ScoreDistribution {
			if count >= 5 {
				break
			}
			percentage := float64(freq) / float64(stats.CompletedSims) * 100
			fmt.Printf("  %s: %d games (%.1f%%)\n", score, freq, percentage)
			count++
		}
	}
	fmt.Println(strings.Repeat("=", 60))
}

func exportSummary_v2(stats *SimulationStats_v2, path string) error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
