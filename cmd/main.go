package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	// Default configuration
	config := struct {
		NumSimulations int
		MaxConcurrent  int
		Team1Name      string
		Team1Strategy  string
		Team2Name      string
		Team2Strategy  string
		GameRules      string
	}{
		NumSimulations: 1, // Default to single simulation
		MaxConcurrent:  runtime.NumCPU(),
		Team1Name:      "Team A",
		Team1Strategy:  "all_in",
		Team2Name:      "Team B",
		Team2Strategy:  "default_half",
		GameRules:      "default",
	}

	// Parse command line arguments
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-n", "--num":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &config.NumSimulations)
				i++
			}
		case "-c", "--cores":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &config.MaxConcurrent)
				i++
			}
		case "-t1", "--team1":
			if i+1 < len(args) {
				config.Team1Strategy = args[i+1]
				i++
			}
		case "-t2", "--team2":
			if i+1 < len(args) {
				config.Team2Strategy = args[i+1]
				i++
			}
		case "-h", "--help":
			printUsage()
			return
		}
	}

	// Run simulation(s)
	if config.NumSimulations == 1 {
		// Single simulation mode
		fmt.Println("Running single simulation...")
		gameID := StartGame(
			config.Team1Name,
			config.Team1Strategy,
			config.Team2Name,
			config.Team2Strategy,
			config.GameRules,
			"",
		)
		fmt.Printf("Simulation completed. Game ID: %s\n", gameID)
	} else {
		// Multiple simulations mode
		err := RunParallelSimulations(SimulationConfig{
			NumSimulations: config.NumSimulations,
			MaxConcurrent:  config.MaxConcurrent,
			Team1Name:      config.Team1Name,
			Team1Strategy:  config.Team1Strategy,
			Team2Name:      config.Team2Name,
			Team2Strategy:  config.Team2Strategy,
			GameRules:      config.GameRules,
		})
		if err != nil {
			fmt.Printf("Error running parallel simulations: %v\n", err)
			os.Exit(1)
		}
	}
}

// Print usage information for command-line arguments
func printUsage() {
	fmt.Println("CS:GO Economy Simulation Usage:")
	fmt.Println("  -n, --num <number>     Number of simulations to run (default: 1)")
	fmt.Println("  -c, --cores <number>   Number of concurrent simulations (default: number of CPU cores)")
	fmt.Println("  -t1, --team1 <strategy> Team 1 strategy (default: all_in)")
	fmt.Println("  -t2, --team2 <strategy> Team 2 strategy (default: default_half)")
	fmt.Println("  -h, --help             Print this help message")
	fmt.Println("\nAvailable strategies:")
	fmt.Println("  all_in                 Always invest all available funds")
	fmt.Println("  default_half           Default strategy that invests half of available funds")
	fmt.Println("\nExamples:")
	fmt.Println("  # Run a single simulation")
	fmt.Println("  go run ./cmd")
	fmt.Println("  # Run 100 simulations using 8 cores")
	fmt.Println("  go run ./cmd -n 100 -c 8")
}
