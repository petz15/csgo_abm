package main

import (
	"CSGO_ABM/internal/analysis"
	"fmt"
	"math"
	"os"
	"runtime"
	"time"
)

// Main entry point for the CS:GO Economy Simulation
//TODO: move gamerules to beginning when everything gets set up, not when each game runs (resourcemanagement and such)
//TODO: solve overtime limit, altough currently it is not a problem
//TODO: the simulation overview file should have a have a unique ID in the name
//TODO: create more models and add better support for them
//TODO: adjust the probabilities in the game engine to make it more realistic/more competitive
//TODO: add support for importing game rules from a file
//TODO: (later) add support for importing models from a file. and AI/machine learning models. no clue how
//TODO: adjust the information the strategy manger gives the models
//TODO: make the resulting terminal output adjustable i.e. simple or pretty

func main() {
	// Default configuration using unified analysis package
	config := analysis.SimulationConfig{
		NumSimulations:        1,                                               // Default to single simulation
		MaxConcurrent:         int(math.Max(1, float64(runtime.NumCPU())*0.8)), // Use 80% of available CPU cores
		MemoryLimit:           3000,                                            // Maximum memory usage before GC (MB)
		Team1Name:             "Team A",
		Team1Strategy:         "all_in",
		Team2Name:             "Team B",
		Team2Strategy:         "default_half",
		Sequential:            false, // Default to parallel simulations
		ExportDetailedResults: false, // Default to not export individual results
		Exportpath:            "",    // Will be set after parsing args
	}

	// Parse command line arguments first to check for custom output path
	customOutputPath := ""
	customGameRulesPath := ""
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
			config.ExportDetailedResults = true
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
		case "-o", "--output":
			if i+1 < len(args) {
				customOutputPath = args[i+1]
				i++
			}
		case "-g", "--gamerules":
			if i+1 < len(args) {
				customGameRulesPath = args[i+1]
				i++
			}
		case "-h", "--help":
			printUsage()
			return
		}
	}

	// Set the results directory - use custom path if specified, otherwise create timestamped directory
	if customOutputPath != "" {
		config.Exportpath = customOutputPath
	} else {
		config.Exportpath = fmt.Sprintf("results_%s", time.Now().Format("20060102_150405"))
	}

	// Create the results directory
	if err := os.MkdirAll(config.Exportpath, 0755); err != nil {
		fmt.Printf("Error creating results directory: %v\n", err)
		os.Exit(1)
	}

	// Validate and prepare all customizations before starting simulations
	customConfig, err := ValidateAndPrepareCustomizations(customGameRulesPath, config.Exportpath)
	if err != nil {
		fmt.Printf("❌ Configuration validation failed: %v\n", err)
		os.Exit(1)
	}

	config.GameRules = customConfig.GameRules

	// Validate team strategies
	if err := ValidateStrategies(config.Team1Strategy, config.Team2Strategy); err != nil {
		fmt.Printf("❌ Strategy validation failed: %v\n", err)
		os.Exit(1)
	}

	// Run simulation(s)
	if config.NumSimulations == 1 {
		// Single simulation mode
		fmt.Println("Running single simulation...")
		result, err := StartGame_default(
			config.Team1Name,
			config.Team1Strategy,
			config.Team2Name,
			config.Team2Strategy,
			customConfig.GameRules,
			"",
			config.ExportDetailedResults,
			config.Exportpath,
		)
		if err != nil {
			fmt.Printf("Error running simulation: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Simulation completed. Game ID: %s\n", result.GameID)
		if result.Team1Won {
			fmt.Printf("Winner: %s (%d-%d)\n", config.Team1Name, result.Team1Score, result.Team2Score)
		} else {
			fmt.Printf("Winner: %s (%d-%d)\n", config.Team2Name, result.Team2Score, result.Team1Score)
		}
	} else if config.Sequential {
		// Sequential simulations mode
		err := sequentialsimulation(config, customConfig.GameRules)
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
	fmt.Println("  -o, --output <path>    Results output directory (default: results_YYYYMMDD_HHMMSS)")
	fmt.Println("  -g, --gamerules <file> Path to JSON file with custom game rules (default: built-in defaults)")
	fmt.Println("  -t1, --team1 <strategy> Team 1 strategy (default: all_in)")
	fmt.Println("  -t2, --team2 <strategy> Team 2 strategy (default: default_half)")
	fmt.Println("  -h, --help             Print this help message")
	fmt.Println("\nAvailable strategies:")
	fmt.Println("  all_in                 Always invest all available funds")
	fmt.Println("  default_half           Default strategy that invests half of available funds")
	fmt.Println("  adaptive_eco_v1        Advanced adaptive economic strategy")
	fmt.Println("  yolo                   Random investment strategy (high risk)")
	fmt.Println("  scrooge                Minimal investment strategy (ultra-conservative)")
	fmt.Println("\nGame Rules Configuration:")
	fmt.Println("  You can customize game parameters using a JSON file. Example:")
	fmt.Println("  go run ./cmd -g example_gamerules.json")
	fmt.Println("  ")
	fmt.Println("  The JSON file can contain any or all of these fields:")
	fmt.Println("  {")
	fmt.Println("    \"defaultEquipment\": 250.0,    // Equipment cost per player")
	fmt.Println("    \"otFunds\": 10000.0,           // Overtime starting funds")
	fmt.Println("    \"startingFunds\": 1000.0,      // Match starting funds")
	fmt.Println("    \"halfLength\": 12,             // Rounds per half")
	fmt.Println("    \"csfR\": 0.7,                  // Contest success function parameter")
	fmt.Println("    \"otHalfLength\": 3             // Overtime rounds per half")
	fmt.Println("  }")
	fmt.Println("  ")
	fmt.Println("  Missing fields will use default values automatically.")
	fmt.Println("\nExamples:")
	fmt.Println("  # Run a single simulation")
	fmt.Println("  go run ./cmd")
	fmt.Println("  # Run 100 simulations using 8 cores")
	fmt.Println("  go run ./cmd -n 100 -c 8")
	fmt.Println("  # Run 1000 simulations with individual result export")
	fmt.Println("  go run ./cmd -n 1000 -c 4 -e")
	fmt.Println("  # Run 500 simulations sequentially with result export to custom folder")
	fmt.Println("  go run ./cmd -n 500 -s -e -o my_results")
	fmt.Println("  # Run 10,000 simulations with optimized memory settings")
	fmt.Println("  go run ./cmd -n 10000 -c 4 -m 500")
	fmt.Println("  # Run 100,000 simulations with memory-efficient settings")
	fmt.Println("  go run ./cmd -n 100000 -c 4 -m 1000")
	fmt.Println("  # Run 1,000,000 simulations (use lower concurrent workers for stability)")
	fmt.Println("  go run ./cmd -n 1000000 -c 2 -m 2000")
	fmt.Println("\nResult Export Options:")
	fmt.Println("  # Export individual results in parallel mode")
	fmt.Println("  go run ./cmd -n 100 -e")
	fmt.Println("  # Export individual results to custom directory")
	fmt.Println("  go run ./cmd -n 100 -s -e -o custom_output")
	fmt.Println("  # Summary-only mode (default, recommended for large simulations)")
	fmt.Println("  go run ./cmd -n 10000")
	fmt.Println("\nLarge-scale simulation recommendations:")
	fmt.Println("  - For 100K+ simulations: Use 4 or fewer cores to avoid memory pressure")
	fmt.Println("  - For 1M+ simulations: Use 2-4 cores with 2GB+ memory limit")
	fmt.Println("  - Avoid -e flag for very large runs (>10K) to prevent filesystem issues")
	fmt.Println("  - Monitor system resources during long-running simulations")
	fmt.Println("  - Results are processed directly in memory for optimal performance")
}
