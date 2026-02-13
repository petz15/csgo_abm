# Dynamic Budget Game (DBG) Agent-Based Model (ABM) Simulation

A high-performance agent-based model for simulating Dynamic Budget Game with empirical CS:GO data, testing various investment strategies. This project combines a parallelized Go simulation engine with Python-based Empirical Game-Theoretic Analysis (EGTA) tools to evaluate and compare strategic decision-making algorithms in Dynamic Budget Game scenarios. The Dynamic Budget Game is an abstracted version of Counter-Strike, including only the decision-making during buy phase.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Go Simulation Engine](#go-simulation-engine)
  - [Installation & Building](#installation--building)
  - [Command Line Usage](#command-line-usage)
  - [Available Strategies](#available-strategies)
  - [Game Modes](#game-modes)
  - [Advanced Features](#advanced-features)
  - [Performance Optimization](#performance-optimization)
- [EGTA & Batch Analysis](#egta--batch-analysis)
  - [Batch Tournament Analysis](#batch-tournament-analysis)
  - [EGTA Jupyter Notebook](#egta-jupyter-notebook)
- [Docker EGTA](#docker-egta)
- [Configuration Files](#configuration-files)
- [Project Structure](#project-structure)
- [Performance Benchmarks](#performance-benchmarks)

## Overview

This ABM simulates the economic decision-making aspect of Counter-Strike: Global Offensive matches. Each simulation runs a complete match between two teams, where each team employs a specific strategy to decide how much money to invest in equipment each round based on game state, opponent behavior, and strategic goals.

**Key Features:**
- 🚀 High-performance concurrent simulation engine written in Go
- 🎯 30+ built-in strategies including ML-based approaches
- 📊 Comprehensive tournament system with round-robin support
- 🔬 EGTA analysis tools for game-theoretic evaluation
- 📈 Detailed statistics export (JSON, CSV, round-by-round data)
- ⚡ Handles 1M+ simulations efficiently with memory management
- 🎲 Thread-safe RNG for reproducible results

## Quick Start

### Prerequisites

- **Go 1.18+** (for simulation engine)
- **Python 3.8+** (for EGTA analysis, optional)
- **Git** (for cloning the repository)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/dbg_abm.git
cd dbg_abm

# Build the simulation executable
go build -o dbg_sim.exe ./cmd

# Or use the pre-built versions if no changes have been made
```

### Run Your First Simulation

```bash
# Run a basic 1v1 simulation (default strategies)
./dbg_sim.exe

# Run 1000 simulations with 8 parallel workers
./dbg_sim.exe -n 1000 -c 8

# Compare two specific strategies
./dbg_sim.exe -n 5000 -t1 min_max_v4 -t2 adaptive_eco_v2 -c 8

# Run a tournament with all strategies
./dbg_sim.exe --tournament -s min_max_v4,adaptive_eco_v2,all_in,anti_allin_v3 --games 10000 -o tournament_results

# Export detailed results with CSV format
./dbg_sim.exe -n 10000 -t1 min_max_v4 -t2 xen_model -e --csv 4 -o results
```

---

## Go Simulation Engine

The core of this project is a highly optimized Go-based simulation engine that can run millions of matches with minimal overhead.

### Installation & Building

#### Standard Build (Current Platform)

```bash
# Build for your current platform
go build -o dbg_sim.exe ./cmd
```

#### Cross-Platform Builds

Use the included VS Code tasks or build manually:

```bash
# Linux (amd64)
GOOS=linux GOARCH=amd64 go build -o dbg_sim_linux ./cmd

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o dbg_sim_macos ./cmd

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o dbg_sim_macos_arm ./cmd

# Windows (from Linux/macOS)
GOOS=windows GOARCH=amd64 go build -o dbg_sim.exe ./cmd
```

### Command Line Usage

```bash
dbg_sim.exe [OPTIONS]

Basic Options:
  -n, --num <N>              Number of simulations to run (default: 1)
  -t1, --team1 <STRATEGY>    Team 1 strategy name (default: all_in)
  -t2, --team2 <STRATEGY>    Team 2 strategy name (default: anti_allin_v3)
  -c, --cores <N>            Number of concurrent workers (default: 80% of CPUs)
  -s, --sequential           Run simulations sequentially instead of parallel
  
Output Options:
  -o, --output <PATH>        Output folder for results (default: timestamped folder)
  -e, --export               Export detailed individual game results #not recommended, use CSV export mode
  -r, --rounds               Export round-by-round data #not recommended, use CSV export mode
  --csv <MODE>               CSV export mode (0-4, see below)
  
Tournament Options:
  --tournament               Run tournament mode instead of single matchup
  -s, --strategies <CSV>     Comma-separated list of strategies for tournament
  --format <FORMAT>          Tournament format (roundrobin) (default: roundrobin)
  --games <N>                Games per matchup in tournament (default: 1000)
  
Advanced Options:
  -a, --advanced             Enable advanced analysis (slower, more detailed) #not recommended, use EGTA for all analysis
  -m, --memory <MB>          Memory limit before forced GC (default: 3000)
  -g, --gamerules <PATH>     Custom game rules JSON file
  -dist, --abmmodels <PATH>  Custom ABM distributions JSON file
  -h, --help                 Show help message

CSV Export Modes:
  0 = No CSV export
  1 = Single game files including RNG values
  2 = All games in one file including RNG values #might be very large
  3 = Single game files only most relevant columns
  4 = Single consolidated file only most relevant columns -> use this with the EGTA
```

### Available Strategies

The simulation includes **30+ strategies** organized into categories:

#### **Aggressive Strategies**
- `all_in` - Maximum investment every round
- `all_in_v2` - Variant with slight moderation

#### **Counter Strategies**
- `anti_allin` - Designed to counter aggressive all-in approaches
- `anti_allin_v2`, `anti_allin_v3` - Refined counter-strategies
- `anti_allin_v3_copy` - Additional v3 instance in order to test two variants

#### **Optimized/Min-Max Strategies**
- `min_max` - Classic minimax approach
- `min_max_v2`, `min_max_v3`, `min_max_v4` - Iterative improvements
- `expected_value` - Expected value-based decision making

#### **Adaptive Strategies**
- `adaptive_eco_v1`, `adaptive_eco_v2` - Adaptive economy management
- `smart_v1` - Context-aware smart investment
- `half` - Conservative half-investment approach

#### **Machine Learning Strategies**
With full opponent visibility:
- `ml_dqn` - Deep Q-Network
- `ml_sgd` - Stochastic Gradient Descent
- `ml_tree` - Decision Tree
- `ml_forest` - Random Forest
- `ml_xgboost` - XGBoost
- `ml_logistic` - Logistic Regression

With forbidden opponent information (`*_forbidden` or `*_f` variants):
- `ml_dqn_forbidden` / `ml_dqn_f`
- `ml_sgd_forbidden` / `ml_sgd_f`
- `ml_tree_forbidden` / `ml_tree_f`
- `ml_forest_forbidden` / `ml_forest_f`
- `ml_xgboost_forbidden` / `ml_xgboost_f`
- `ml_logistic_forbidden` / `ml_logistic_f`

#### **Additional ML Models**
- `xen_model` - Advanced XGBoost model with specialized training

#### **Basic/Baseline Strategies**
- `casual` - Balanced casual play
- `random` - Random investment amounts
- `scrooge` - Minimal investment, maximum saving

### Game Modes

#### Single Matchup Mode (Default)

Run repeated matches between two strategies:

```bash
# Basic head-to-head
./dbg_sim.exe -n 10000 -t1 min_max_v4 -t2 adaptive_eco_v2 -c 8

# With detailed exports
./dbg_sim.exe -n 5000 -t1 all_in -t2 anti_allin_v3 -e -r -o matchup_results
```

**Output:**
- Console: Real-time progress, final statistics, winner
- `simulation_summary.json`: Aggregated statistics
- `all_games_minimal.csv`: Game-by-game results (if `--csv 4`)
- Individual game CSVs (if `-e` flag)

#### Tournament Mode

Run round-robin tournaments with multiple strategies:

```bash
# Tournament with 4 strategies
./dbg_sim.exe --tournament \
  -s min_max_v4,adaptive_eco_v2,all_in_v2,anti_allin_v3 \
  --games 10000 \
  -c 8 \
  -o tournament_results

# Large-scale tournament with all ML models
./dbg_sim.exe --tournament \
  -s ml_xgboost,ml_forest,ml_dqn,xen_model,min_max_v4 \
  --games 20000 \
  --csv 4 \
  -o ml_tournament
```

**Output Structure:**
```
tournament_results/
├── matchup_001_min_max_v4_vs_adaptive_eco_v2/
│   ├── simulation_summary.json
│   └── all_games_minimal.csv
├── matchup_002_min_max_v4_vs_all_in_v2/
│   ├── simulation_summary.json
│   └── all_games_minimal.csv
├── ...
└── tournament_standings.txt
```

**Tournament Features:**
- Round-robin format (all vs all)
- Each matchup gets its own folder
- Automatic standings calculation
- Win/loss records for each strategy
- Point-based ranking system

### Advanced Features

#### Custom Game Rules

Modify game parameters using JSON configuration files:

```bash
./dbg_sim.exe -n 1000 -g alt_gamerules_robustness_1.json
```

Example custom game rules structure:
```json
{
  "StartingFunds": 800,
  "HalfLength": 15,
  "OTHalfLength": 3,
  "DefaultEquipment": 200,
  "BombPlantReward": 300,
  ...
}
```

Available custom game rule files in repository:
- `alt_gamerules.json` - Base alternative rules
- `alt_gamerules_robustness_*.json` - Robustness test variants

#### Custom ABM Distributions

Specify custom probability distributions for game outcomes:

```bash
./dbg_sim.exe -n 5000 -dist distributions.json
```

The distributions file defines probabilities for:
- Round outcomes (T win, CT win)
- Equipment effectiveness
- Economic factors

#### Memory Management

For large-scale simulations (100K+ games), control memory usage:

```bash
# Limit memory to 2GB before forcing garbage collection
./dbg_sim.exe -n 1000000 -c 8 -m 2000

# Lower memory limit for resource-constrained systems
./dbg_sim.exe -n 500000 -c 4 -m 1500
```

#### Advanced Analysis Mode

Enable deeper statistical analysis (slower, more comprehensive):

```bash
./dbg_sim.exe -n 10000 -t1 min_max_v4 -t2 xen_model -a
```

Additional analysis includes:
- Round-by-round investment patterns
- Economy state transitions
- Win probability curves
- Strategy effectiveness by game phase

### Performance Optimization

#### Parallel Execution

The engine uses worker pools for maximum throughput:

```bash
# Auto-detect optimal cores (80% of available)
./dbg_sim.exe -n 100000

# Manual core specification
./dbg_sim.exe -n 100000 -c 12

# Sequential (for debugging or single-core machines)
./dbg_sim.exe -n 1000 -s
```

**Performance Tips:**
- Use `-c` with 1-2 fewer cores than available (leaves headroom for OS)
- For 1M+ simulations, increase `-m` to reduce GC overhead
- Use `--csv 4` (single file) instead of individual CSVs for large runs
- Disable `-e` and `-r` flags for maximum speed

#### Output Optimization

Choose the right export mode for your needs:

```bash
# Maximum speed: No CSV export
./dbg_sim.exe -n 1000000 -c 8

# Minimal overhead: Summary only
./dbg_sim.exe -n 500000 -c 8 --csv 1

# Consolidated export: Single CSV file (recommended)
./dbg_sim.exe -n 100000 -c 8 --csv 4

# Full detail: All exports (slowest, for analysis)
./dbg_sim.exe -n 10000 -c 8 -e -r --csv 3
```

#### Benchmarks

Typical performance on modern hardware (AMD Ryzen 9 / Intel i7):

| Simulations | Cores | Avg Time | Throughput |
|------------|-------|----------|------------|
| 1,000      | 8     | ~3s      | ~333 games/s |
| 10,000     | 8     | ~25s     | ~400 games/s |
| 100,000    | 8     | ~4min    | ~420 games/s |
| 1,000,000  | 8     | ~40min   | ~420 games/s |

*Performance varies by strategy complexity and export options*

---

## EGTA & Batch Analysis

The Python components provide game-theoretic analysis and visualization tools for tournament results.

### Batch Tournament Analysis

The `batch_analyze_tournament.py` script automates analysis of all matchups in a tournament folder.

#### Installation

```bash
# Install Python dependencies
pip install papermill jupyter pandas numpy matplotlib seaborn
```

#### Usage

```bash
# Analyze tournament folder
python batch_analyze_tournament.py tournament_results/

# Specify custom notebook template
python batch_analyze_tournament.py tournament_results/ --notebook custom_template.ipynb

# Disable first-round analysis (faster)
python batch_analyze_tournament.py tournament_results/ --no-first-round

# Parallel processing
python batch_analyze_tournament.py tournament_results/ --workers 4
```

#### Features

- **Automatic Matchup Detection**: Finds all `matchup_XXX` folders
- **Parallel Processing**: Analyzes multiple matchups concurrently
- **HTML Report Generation**: Creates interactive visualizations
- **Error Handling**: Continues analysis even if individual matchups fail
- **Progress Tracking**: Real-time progress indicators

#### Output

For each matchup folder:
```
matchup_001_strategy1_vs_strategy2/
├── simulation_summary.json
├── all_games_minimal.csv
├── analysis_report.html           # Generated analysis
└── analysis_report.ipynb           # Executed notebook
```

### EGTA Jupyter Notebook

The `EGTA.ipynb` notebook provides comprehensive empirical game-theoretic analysis.
It automatically produces a overview of all the graphs and calculations in a HTML file.
When running the batch analyzer, the results can be reviewed in the HTML files. 

#### Running the Notebook

```bash
# Start Jupyter
jupyter notebook EGTA.ipynb

# Or use Jupyter Lab
jupyter lab EGTA.ipynb
```

Adjust the FOLDER_PATH = r"(INSERT YOUR FOLDER PATH HERE)"


## Docker EGTA

A containerized environment for running EGTA analysis without local Python setup.

### Quick Start

```bash
cd docker_EGTA

# Build container
docker-compose build

# Run batch analysis
./run_analysis.sh ../tournament_results

# Start Jupyter server
./run_jupyter.sh
```

### Included Tools

- Python 3.9+ with scientific stack
- Jupyter Lab
- All required dependencies (papermill, pandas, matplotlib, seaborn)
- Pre-configured environment for EGTA notebook

See [docker_EGTA/README.md](docker_EGTA/README.md) for detailed Docker instructions.

---

## Configuration Files

### Game Rules

JSON files defining CS:GO match parameters:

- `alt_gamerules.json` - Standard alternative ruleset
- `alt_gamerules_robustness_*.json` - Robustness test variants with modified parameters

Key parameters:
- `StartingFunds`: Initial money per team
- `HalfLength`: Rounds per half (usually 15)
- `OTHalfLength`: Overtime rounds per half (usually 3)
- `DefaultEquipment`: Base equipment cost
- Round rewards and loss bonuses

### Distributions

`distributions.json` defines probability distributions for:
- Round win probabilities by equipment advantage
- Economic factors (save rounds, force buys)
- Bomb plant/defuse probabilities

### ML Models

Pre-trained machine learning models in `ml_models/`:
- `xgboost_model.json` - XGBoost classifier
- `forest_model.json` - Random Forest
- `tree_model.json` - Decision Tree
- `q_network_weights.json` - DQN weights
- `*_forbidden.json` - Models without opponent visibility

---

## Project Structure

```
dbg_abm/
├── cmd/                          # Main application entry point
│   ├── main.go                   # CLI parsing and program entry
│   ├── simulation_concurrent.go  # Parallel simulation engine
│   ├── simulation_sequential.go  # Sequential simulation engine
│   ├── tournament_runner.go      # Tournament management
│   ├── gamehandler.go           # Game initialization and execution
│   └── custom.go                # Custom configuration handling
│
├── internal/                     # Core simulation logic
│   ├── engine/                  # Game engine
│   │   ├── game.go              # Match state machine
│   │   ├── round.go             # Round simulation
│   │   ├── team.go              # Team management
│   │   ├── gamerules.go         # Rule validation and loading
│   │   ├── probabilities.go     # Outcome probability calculations
│   │   ├── distributions_loader.go # ABM distribution loading
│   │   └── strategymanager.go   # Strategy execution interface
│   │
│   ├── strategy/                # Investment strategies
│   │   ├── registry.go          # Strategy registration
│   │   ├── strategycontext.go   # Context passed to strategies
│   │   ├── all_in*.go           # Aggressive strategies
│   │   ├── anti_allin*.go       # Counter strategies
│   │   ├── min_max*.go          # Minimax optimized strategies
│   │   ├── adaptive_eco*.go     # Adaptive strategies
│   │   ├── ml_*.go              # ML-based strategies
│   │   ├── xen_model.go         # Advanced XGBoost model
│   │   └── ...                  # Other strategies
│   │
│   ├── analysis/                # Statistics and export
│   │   ├── stats.go             # Statistics aggregation
│   │   ├── calculator.go        # Metric calculations
│   │   ├── reporter.go          # Console reporting
│   │   ├── exporter.go          # JSON/CSV export
│   │   └── tournament_export.go # Tournament-specific exports
│   │
│   └── tournament/              # Tournament management
│       └── tournament.go        # Tournament logic
│
├── util/                        # Utility functions
│   ├── export_csv.go           # CSV export utilities
│   ├── export_results.go       # Result serialization
│   └── id_generation.go        # Unique ID generation
│
├── ml_models/                   # Pre-trained ML models
├── xen_model/                   # Xen model files
├── docker_EGTA/                 # Docker configuration
│
├── batch_analyze_tournament.py  # Batch analysis script
├── EGTA.ipynb                   # EGTA analysis notebook
├── analysis_notebook_CSF.ipynb  # Additional analysis tools
│
├── *.json                       # Game rules and configurations
├── go.mod                       # Go module definition
└── .vscode/tasks.json          # VS Code tasks (not shown)
```

---


## Citation

If you use this simulation in research, please cite:

```

```

---


## Acknowledgments

- CS:GO game mechanics based on official Valve data
- EGTA methodology from game theory literature

