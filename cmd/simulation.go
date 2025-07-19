package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SimulationConfig holds the configuration for a batch of simulations
type SimulationConfig struct {
	NumSimulations int
	MaxConcurrent  int
	Team1Name      string
	Team1Strategy  string
	Team2Name      string
	Team2Strategy  string
	GameRules      string
}

// RunParallelSimulations runs multiple game simulations in parallel
func RunParallelSimulations(config SimulationConfig) error {
	fmt.Printf("Running %d simulations using %d cores...\n", config.NumSimulations, config.MaxConcurrent)
	fmt.Printf("Team 1: %s with '%s' strategy\n", config.Team1Name, config.Team1Strategy)
	fmt.Printf("Team 2: %s with '%s' strategy\n", config.Team2Name, config.Team2Strategy)

	// Create a wait group to wait for all simulations to complete
	var wg sync.WaitGroup

	// Create a semaphore channel to limit concurrent executions
	sem := make(chan bool, config.MaxConcurrent)

	// Start time for performance measurement
	startTime := time.Now()

	// Create results directory if it doesn't exist
	resultsDir := "results"
	os.MkdirAll(resultsDir, 0755)

	// Run simulations
	for i := 0; i < config.NumSimulations; i++ {
		wg.Add(1)

		// Run each simulation in a goroutine
		go func(simIndex int) {
			defer wg.Done()
			sem <- true              // Acquire semaphore
			defer func() { <-sem }() // Release semaphore when done

			fmt.Printf("Starting simulation %d...\n", simIndex+1)

			simPrefix := fmt.Sprintf("sim_%d_", simIndex+1)
			StartGame(
				config.Team1Name,
				config.Team1Strategy,
				config.Team2Name,
				config.Team2Strategy,
				config.GameRules,
				simPrefix,
			)

			fmt.Printf("Completed simulation %d\n", simIndex+1)
		}(i)
	}

	// Wait for all simulations to complete
	wg.Wait()

	// Calculate and display execution time
	duration := time.Since(startTime)
	fmt.Printf("All simulations completed in %s\n", duration)

	// Organize results
	return organizeAndSummarizeResults()
}

// organizeAndSummarizeResults moves simulation results to a timestamped directory
// and generates summary statistics
func organizeAndSummarizeResults() error {
	// Create timestamp-based directory name for this batch of simulations
	batchID := time.Now().Format("20060102_150405")
	resultsDir := filepath.Join("results", batchID)

	// Move individual result files to the batch directory
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("error creating results directory: %w", err)
	}

	// Move files
	files, err := os.ReadDir("results")
	if err != nil {
		return fmt.Errorf("error reading results directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		oldPath := filepath.Join("results", file.Name())
		newPath := filepath.Join(resultsDir, file.Name())
		if err := os.Rename(oldPath, newPath); err != nil {
			fmt.Printf("Warning: Could not move file %s: %v\n", file.Name(), err)
		}
	}
	/*
		// Aggregate and save summary statistics
		summaryFile := filepath.Join(resultsDir, "simulation_summary.json")
		if err := aggregateResults(resultsDir, summaryFile); err != nil {
			return fmt.Errorf("error aggregating results: %w", err)
		}

		fmt.Printf("Summary statistics saved to %s\n", summaryFile)

	*/

	fmt.Printf("All result files moved to %s\n", resultsDir)

	return nil
}

// gameStats represents summary data for a single game
type gameStats struct {
	ID            string `json:"id"`
	Team1Name     string `json:"team1_name"`
	Team1Strategy string `json:"team1_strategy"`
	Team2Name     string `json:"team2_name"`
	Team2Strategy string `json:"team2_strategy"`
	Team1Score    int    `json:"team1_score"`
	Team2Score    int    `json:"team2_score"`
	Team1Win      bool   `json:"team1_win"`
	RoundsPlayed  int    `json:"rounds_played"`
	OvertimeFlag  bool   `json:"overtime_flag"`
}

// simulationStats represents aggregate statistics across multiple games
type simulationStats struct {
	TotalGames        int            `json:"total_games"`
	Team1WinRate      float64        `json:"team1_win_rate"`
	Team2WinRate      float64        `json:"team2_win_rate"`
	AverageRounds     float64        `json:"average_rounds"`
	OvertimeRate      float64        `json:"overtime_rate"`
	ScoreDistribution map[string]int `json:"score_distribution"`
	Games             []gameStats    `json:"games"`
	SimulationTime    string         `json:"simulation_time"`
}

//currently not used, not sure if I want to keep this or do it in python. And also the way
// this is structured is not very efficient for large numbers of games

// aggregateResults generates summary statistics from individual game results
func aggregateResults(resultsDir, outputFile string) error {
	// List all JSON files in the results directory
	files, err := os.ReadDir(resultsDir)
	if err != nil {
		return fmt.Errorf("error reading results directory: %w", err)
	}

	// Initialize summary data
	summary := simulationStats{
		ScoreDistribution: make(map[string]int),
	}

	team1Wins := 0
	totalRounds := 0
	overtimeGames := 0

	// Process each JSON file
	for _, file := range files {
		if file.IsDir() || file.Name() == "simulation_summary.json" {
			continue
		}

		// This is a simplified placeholder - in a real implementation
		// you would parse each JSON file and extract the relevant data
		// Then accumulate statistics in the summary structure

		// For now, just add a placeholder game summary
		gameSummary := gameStats{
			ID:            file.Name(),
			Team1Name:     "Team A", // These would be parsed from the file
			Team1Strategy: "all_in",
			Team2Name:     "Team B",
			Team2Strategy: "default_half",
			Team1Score:    16, // Example score
			Team2Score:    14,
			Team1Win:      true,
			RoundsPlayed:  30,
			OvertimeFlag:  false,
		}

		// Update statistics
		if gameSummary.Team1Win {
			team1Wins++
		}
		totalRounds += gameSummary.RoundsPlayed
		if gameSummary.OvertimeFlag {
			overtimeGames++
		}

		// Update score distribution
		scoreKey := fmt.Sprintf("%d-%d", gameSummary.Team1Score, gameSummary.Team2Score)
		summary.ScoreDistribution[scoreKey]++

		summary.Games = append(summary.Games, gameSummary)
	}

	// Calculate summary statistics
	summary.TotalGames = len(summary.Games)
	if summary.TotalGames > 0 {
		summary.Team1WinRate = float64(team1Wins) / float64(summary.TotalGames)
		summary.Team2WinRate = 1.0 - summary.Team1WinRate
		summary.AverageRounds = float64(totalRounds) / float64(summary.TotalGames)
		summary.OvertimeRate = float64(overtimeGames) / float64(summary.TotalGames)
	}

	summary.SimulationTime = time.Now().Format(time.RFC3339)

	// Write summary to file
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling summary data: %w", err)
	}

	return os.WriteFile(outputFile, data, 0644)
}
