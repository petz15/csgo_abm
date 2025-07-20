package main

import (
	"fmt"
	"math"
	"os"
	"runtime"

	"CSGO_ABM/internal/analysis"
)

func main() {
	// Default configuration using unified analysis package
	config := analysis.SimulationConfig{
		NumSimulations: 1,                                               // Default to single simulation
		MaxConcurrent:  int(math.Max(1, float64(runtime.NumCPU())*0.8)), // Use 80% of available CPU cores
		MemoryLimit:    3000,                                            // Maximum memory usage before GC (MB)
		Team1Name:      "Team A",
		Team1Strategy:  "all_in",
		Team2Name:      "Team B",
		Team2Strategy:  "default_half",
		GameRules:      "default",
		Sequential:     false, // Default to parallel simulations
		ExportResults:  false, // Default to not export individual results
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
		case "-s", "--sequential":
			config.Sequential = true
		case "-e", "--export":
			config.ExportResults = true
		case "-c", "--cores":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &config.MaxConcurrent)
				i++
			}
		case "-m", "--memory":
			if i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &config.MemoryLimit)
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
	} else if config.Sequential {
		// Sequential simulations mode
		err := sequentialsimulation(config)
		if err != nil {
			fmt.Printf("Error running sequential simulations: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Multiple simulations mode
		err := RunParallelSimulations(config)
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
	fmt.Println("  -m, --memory <number>  Memory limit in MB before forcing GC (default: 3000)")
	fmt.Println("  -s, --sequential       Run simulations sequentially instead of in parallel")
	fmt.Println("  -e, --export           Export individual game results as JSON files (works with both modes)")
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
	fmt.Println("  # Run 1000 simulations with individual result export")
	fmt.Println("  go run ./cmd -n 1000 -c 4 -e")
	fmt.Println("  # Run 500 simulations sequentially with result export")
	fmt.Println("  go run ./cmd -n 500 -s -e")
	fmt.Println("  # Run 10,000 simulations with optimized memory settings")
	fmt.Println("  go run ./cmd -n 10000 -c 4 -m 500")
	fmt.Println("  # Run 100,000 simulations with memory-efficient settings")
	fmt.Println("  go run ./cmd -n 100000 -c 4 -m 1000")
	fmt.Println("  # Run 1,000,000 simulations (use lower concurrent workers for stability)")
	fmt.Println("  go run ./cmd -n 1000000 -c 2 -m 2000")
	fmt.Println("\nResult Export Options:")
	fmt.Println("  # Export individual results in parallel mode")
	fmt.Println("  go run ./cmd -n 100 -e")
	fmt.Println("  # Export individual results in sequential mode")
	fmt.Println("  go run ./cmd -n 100 -s -e")
	fmt.Println("  # Summary-only mode (default, recommended for large simulations)")
	fmt.Println("  go run ./cmd -n 10000")
	fmt.Println("\nLarge-scale simulation recommendations:")
	fmt.Println("  - For 100K+ simulations: Use 4 or fewer cores to avoid memory pressure")
	fmt.Println("  - For 1M+ simulations: Use 2-4 cores with 2GB+ memory limit")
	fmt.Println("  - Avoid -e flag for very large runs (>10K) to prevent filesystem issues")
	fmt.Println("  - Monitor system resources during long-running simulations")
	fmt.Println("  - Results are processed directly in memory for optimal performance")
}
