# CS:GO Agent-Based Econom   └── models
       ├── all_in.go                 # "All-in" spending strategy
       ├── default_half.go           # "Default-half" spending strategy
       ├── adaptive_eco_v1.go        # Advanced adaptive economic strategy
       ├── adaptive_eco_v2.go        # Enhanced contextual adaptive strategy
       ├── yolo.go                   # High-risk random investment strategy
       └── scrooge.go                # Ultra-conservative minimal investment strategyodel (CSGO_ABM)

This project implements an agent-based model simulating the economy aspect of Counter-Strike: Global Offensive (CS:GO). The simulation focuses on how teams manage their finances, make investment decisions, and the impact of these decisions on game outcomes through probabilistic models.

## Project Structure

```
CSGO_ABM
├── cmd
│   ├── main.go                        # Entry point with unified CLI
│   ├── gamehandler.go                # Game handling and orchestration
│   ├── simulation_concurrent.go      # High-performance parallel simulations
│   └── simulation_sequential.go      # Sequential simulation for comparison
├── internal
│   ├── analysis                      # Unified analysis package
│   │   ├── stats.go                  # Comprehensive statistics structure
│   │   ├── calculator.go             # Statistical computation engine
│   │   ├── reporter.go               # Enhanced reporting with insights
│   │   └── exporter.go               # Multi-format export system
│   ├── engine
│   │   ├── game.go                   # Core game logic and state management
│   │   ├── gamerules.go              # Game rules and configuration
│   │   ├── probabilities.go          # Probability distributions and CSF functions
│   │   ├── round.go                  # Individual round simulation
│   │   ├── strategymanager.go        # Strategy selection and implementation
│   │   └── team.go                   # Team state and behavior
│   └── models
│       ├── all_in.go                 # "All-in" spending strategy
│       └── default_half.go           # "Default-half" spending strategy
├── output
│   └── results.go                    # Functions for exporting simulation results
├── results_*/                        # Timestamped directories with simulation outputs
├── go.mod                            # Module definition
├── go.sum                            # Dependency checksums
├── SIMULATION_GUIDE.md               # Comprehensive simulation guide
├── internal/analysis/README.md       # Analysis package documentation
└── README.md                         # Project documentation
```

## Key Features

### 🚀 **High-Performance Parallel Processing**
- **Worker Pool Architecture**: Efficient concurrent simulation processing with configurable worker counts
- **Memory Management**: Intelligent memory monitoring with automatic garbage collection
- **Scalable Design**: Handles 100,000+ simulations with optimized resource usage
- **Fault Tolerance**: Robust error handling with timeout mechanisms and graceful degradation

### 🎮 **Realistic Game Simulation**
- Random assignment of teams to Counter-Terrorist (CT) and Terrorist (T) sides
- Dynamic economy system with realistic fund allocation and spending strategies
- Contest Success Function (CSF) for probabilistic outcome determination with customizable skewness
- Equipment value tracking and survival calculations
- Complete match simulation including:
  - Half-time side swapping after 15 rounds
  - Overtime mechanics when tied at 15-15
  - Loss bonus calculation based on consecutive losses
  - Bomb plant mechanics and related bonuses
- **Advanced Strategy System**: Multiple AI strategies with contextual decision-making
- **Enhanced Probability Models**: Skewed normal distributions for realistic outcome variance

### 📊 **Unified Analysis System**
- **Advanced Statistics**: Win streak analysis, statistical significance testing, and balance scoring
- **Enhanced Reporting**: Rich visual output with insights and strategy recommendations
- **Smart Insights**: Strategy imbalance detection, competitiveness analysis, and optimization suggestions
- **Multiple Export Formats**: JSON, CSV exports for comprehensive data analysis
- **Real-time Monitoring**: Progress tracking with performance metrics and memory optimization
- **Thread-Safe Operations**: Atomic updates for concurrent simulation processing

### 🔧 **Latest Improvements (v2.5)**
- **Enhanced Strategy System**: New `adaptive_eco_v2` with contextual awareness and psychological modeling
- **Advanced Probability Models**: Skewed normal distributions with momentum effects
- **Economic Intelligence**: Strategies now consider opponent state, loss bonuses, and round importance
- **Multiple Strategy Types**: 6 distinct strategies from ultra-conservative to high-risk gambling
- **Unified Analysis Package**: Consolidated statistical processing across all simulation modes
- **Enhanced Reporting**: Rich visual output with strategic insights and recommendations
- **Code Simplification**: Eliminated duplicate code, unified configuration system
- **Advanced Metrics**: Statistical significance testing, balance scoring, competitiveness analysis
- **Multi-format Export**: JSON, CSV exports with comprehensive data analysis options
- **Thread-Safe Design**: Atomic operations and proper concurrency patterns throughout

## Probability Models

The simulation uses several advanced probability models to determine outcomes:

1. **Contest Success Function (CSF)**: Calculates the probability of winning based on team expenditures
   - Uses Tullock contest model with adjustable parameter r
   - Higher r values make outcomes more deterministic

2. **CSF Normal Distribution with Skewness**: Generates values from a normal distribution with:
   - Mean positioned according to the CSF probability
   - Configurable standard deviation (controlled by stdDevFactor)
   - **Skewness support**: Azzalini's skew-normal transformation for realistic asymmetric outcomes
   - Configurable output range for different game aspects
   - Perfect for modeling momentum effects and psychological factors

3. **Enhanced Strategy System**: AI strategies with contextual awareness:
   - **Economic state assessment**: Healthy, moderate, poor, critical funding levels
   - **Round-specific logic**: Pistol rounds, anti-eco situations, overtime pressure
   - **Side-specific adjustments**: CT vs T economic differences
   - **Score pressure dynamics**: Aggressive play when behind, conservative when ahead
   - **Loss bonus optimization**: Strategic use of consecutive loss economics

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/petz15/csgo_abm.git
   cd csgo_abm
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Run the simulation:
   ```bash
   # Single simulation
   go run ./cmd
   
   # Quick test with enhanced reporting
   go run ./cmd -n 10 -s
   
   # Large-scale parallel simulation
   go run ./cmd -n 1000 -c 4
   ```

## Usage

### Running Large-Scale Simulations

The parallel simulation system is optimized for large-scale runs:

```bash
# Run 100,000 simulations with 16 workers and 2GB memory limit
go run ./cmd -n 100000 -c 16 -m 2000

# Quick performance test with monitoring
go run ./cmd -n 10000 -c 8
```

#### Performance Optimizations

- **Automatic scaling**: Memory limits and worker counts adjust for very large simulations (100k+)
- **Progress monitoring**: Real-time updates every 10-30 seconds depending on simulation size
- **Memory management**: Automatic garbage collection when memory usage exceeds limits
- **Bounded growth**: Score distribution tracking limited to prevent memory issues
- **Graceful shutdown**: Proper cleanup even if interrupted

#### Enhanced Simulation Output

Each simulation run provides comprehensive analysis with:
- **Rich Statistical Summary**: Enhanced visual reporting with insights
- **Multiple Export Formats**: JSON (complete data), CSV (key metrics), distribution files
- **Strategic Insights**: Balance analysis, competitiveness scoring, optimization recommendations
- **Performance Metrics**: Memory usage, processing rates, and execution timing

Example enhanced output:
```
============================================================
🎮 SIMULATION RESULTS SUMMARY
============================================================
Simulation Mode: Concurrent
Total Simulations: 1,000 (Completed: 1,000, Failed: 0)
Execution Time: 4s

📊 GAME ANALYSIS
----------------------------------------
Team 1 Wins: 523 (52.3%)
Team 2 Wins: 477 (47.7%)
Average Rounds per Game: 28.7
Overtime Rate: 23.4%
Strategy Balance Score: 89.2% (Excellent)
Statistical Significance: High confidence (p<0.01)

Game Length: Standard (avg 28.7 rounds)
Close Games (≤3 round diff): 234 (23.4%)
Blowout Games (>10 round diff): 45 (4.5%)
Competitiveness: High

📈 SCORE DISTRIBUTIONS
----------------------------------------
Most Common Scores:
  16-14: 89 games (8.9%)
  16-13: 76 games (7.6%)
  16-12: 67 games (6.7%)

🔍 INSIGHTS & RECOMMENDATIONS
----------------------------------------
✅ Well-Balanced Strategies: Both teams show competitive performance
⚡ Moderate Overtime Rate: Games often decided in regulation
📊 Statistically Significant: Results are reliable with this sample size
============================================================
```

#### Command-Line Arguments

- `-n, --num <number>`: Number of simulations to run (default: 1)
- `-c, --cores <number>`: Number of concurrent workers (default: CPU cores × 0.8)
- `-m, --memory <MB>`: Memory limit in MB before forcing GC (default: 3000)
- `-s, --sequential`: Run simulations sequentially instead of parallel
- `-e, --export`: Export individual game results as JSON files
- `-t1, --team1 <strategy>`: Team 1 strategy (default: all_in)
- `-t2, --team2 <strategy>`: Team 2 strategy (default: default_half)
- `-h, --help`: Print comprehensive help message

#### Export Options

```bash
# Summary-only mode (recommended for large runs)
go run ./cmd -n 10000

# Export individual results and comprehensive analysis
go run ./cmd -n 100 -e

# Sequential mode with exports
go run ./cmd -n 500 -s -e

# Multiple export formats automatically generated:
# - simulation_summary.json (complete data)
# - summary.csv (key metrics)
# - score_distribution.csv (score frequency analysis)
# - round_distribution.csv (game length analysis)
```

### Performance Benchmarks

Based on testing with the unified analysis system:

| Simulations | Workers | Memory Limit | Avg. Time | Rate (sims/sec) | Analysis Features |
|-------------|---------|--------------|-----------|-----------------|-------------------|
| 100         | 2       | 3000MB       | 0.5s      | 200/sec         | Full insights     |
| 1,000       | 4       | 3000MB       | 4.2s      | 238/sec         | Statistical sig.  |
| 10,000      | 8       | 3000MB       | 42s       | 238/sec         | Advanced metrics  |
| 100,000     | 4       | 3000MB       | 7.1min    | 234/sec         | All features      |

*Enhanced reporting and analysis add minimal overhead while providing comprehensive insights.*

### Creating Custom Strategies

To implement a custom strategy:

1. Create a new file in the `internal/models` directory
2. Implement a function with the signature:
   ```go
   func InvestDecisionMaking_yourStrategy(funds float64, curround int, curscoreopo int) float64
   ```
3. Add your strategy to the strategy manager:
   ```go
   // In internal/engine/strategymanager.go
   case "your_strategy":
       return models.InvestDecisionMaking_yourStrategy(team.Funds, curround, curscoreopo)
   ```

## Technical Implementation

### Unified Analysis Architecture

The simulation features a comprehensive analysis system with:

1. **SimulationStats**: Unified statistics structure for both sequential and concurrent modes
2. **Statistical Calculator**: Thread-safe computation engine with atomic operations
3. **Enhanced Reporting**: Rich visual output with strategic insights and recommendations
4. **Multi-format Export**: JSON, CSV, and distribution analysis files
5. **Advanced Metrics**: Balance scoring, statistical significance testing, competitiveness analysis

### Concurrency Architecture

The simulation uses a robust worker pool pattern with:

1. **Worker Pool**: Fixed number of goroutines for simulation processing
2. **Thread-Safe Statistics**: Atomic operations for all shared state updates
3. **Result Collection**: Unified analysis pipeline for both modes
4. **Memory Monitoring**: Background monitoring with automatic GC triggers
5. **Progress Reporting**: Real-time status with performance metrics

### Analysis Features

- **Win Streak Analysis**: Maximum and average win streak calculations
- **Statistical Significance**: Chi-square testing for result reliability
- **Balance Scoring**: Strategy competitiveness measurement (0-100%)
- **Game Categorization**: Close games vs blowouts analysis
- **Distribution Analysis**: Score and round length frequency analysis
- **Strategic Insights**: Automated recommendations and warnings

### Error Handling

- **Timeout Protection**: Individual simulations timeout after 5 minutes
- **Panic Recovery**: Graceful handling of simulation panics
- **Graceful Degradation**: System continues operation even with partial failures
- **Comprehensive Logging**: Detailed error reporting for debugging

### Analyzing Results

Results are exported as JSON files with unique identifiers:
- **Individual Results**: Format `YYYYMMDD_HHMMSS_hostname_randomID.json`
- **Summary Statistics**: Aggregate data in `simulation_summary.json`
- **Performance Metrics**: Execution time, processing rate, memory usage

The summary includes:
- Win rates and score distributions for both teams
- Average rounds per game and overtime statistics
- Performance metrics (simulations/second, peak memory usage)
- System resource utilization data

## Features

- Random assignment of teams to Counter-Terrorist (CT) and Terrorist (T) sides
- Dynamic economy system with realistic fund allocation and spending strategies
- Contest Success Function (CSF) for probabilistic outcome determination
- Normal distribution sampling based on team expenditures with customizable parameters
- Equipment value tracking and survival calculations
- Complete match simulation including:
  - Half-time side swapping after 15 rounds
  - Overtime mechanics when tied at 15-15
  - Loss bonus calculation based on consecutive losses
  - Bomb plant mechanics and related bonuses
- Strategy implementation through pluggable strategy modules
- Parallel simulation processing with configurable concurrency
- Simulation result aggregation and statistical analysis
- JSON export of simulation results with unique game identifiers

## Quick Reference

### Essential Commands
```bash
# Basic simulation
go run ./cmd

# Test run with enhanced reporting  
go run ./cmd -n 10 -s

# Large-scale simulation
go run ./cmd -n 1000 -c 4

# Full analysis with exports
go run ./cmd -n 100 -e
```

### Key Flags
- `-n <number>`: Number of simulations
- `-s`: Sequential mode (vs parallel)
- `-c <workers>`: Concurrent worker count
- `-e`: Export individual results
- `-m <MB>`: Memory limit
- `-h`: Help

### Documentation
- **[SIMULATION_GUIDE.md](SIMULATION_GUIDE.md)**: Comprehensive usage guide
- **[internal/analysis/README.md](internal/analysis/README.md)**: Analysis package documentation
- **[Go Docs](https://pkg.go.dev)**: API documentation (run `go doc ./...`)

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Commit changes: `git commit -am 'Add feature'`
4. Push to branch: `git push origin feature-name`
5. Submit a Pull Request

### Development Setup
```bash
git clone https://github.com/petz15/csgo_abm.git
cd csgo_abm
go mod tidy
go test ./...
go run ./cmd -n 10 -s  # Test run
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
```bash
# Run single simulation with detailed output
go run ./cmd -n 1 -c 1

# Use sequential implementation for step-by-step debugging
# (modify main.go to call sequentialsimulation instead)
```

## Contributing

Contributions are welcome! The project has been significantly improved with robust concurrency handling and is ready for community contributions.

### Areas for Contribution
- **Strategy Development**: Implement new economic strategies in `internal/models/`
- **Machine Learning Integration**: Add ML-based adaptive strategies and opponent modeling
- **Performance Optimization**: Improve simulation speed and memory efficiency  
- **Testing**: Add comprehensive test suites for parallel simulation reliability
- **Documentation**: Improve code documentation and usage examples
- **Analysis Tools**: Create visualization and statistical analysis utilities
- **Game Mechanics**: Add new CS:GO features (weapon mechanics, map-specific strategies)
- **Web Dashboard**: Real-time monitoring and interactive strategy comparison
- **Tournament Simulation**: Multi-game series and bracket simulation
- **Economic Intelligence**: Opponent fund estimation and counter-strategy development

### High-Priority Improvements
1. **Enhanced Strategy Migration**: Move all strategies to contextual decision-making system
2. **Economic Intelligence**: Opponent fund estimation and behavior prediction
3. **Real-time Dashboard**: Web-based simulation monitoring and analysis
4. **Machine Learning Framework**: Foundation for adaptive and learning strategies
5. **Advanced Game Mechanics**: Map control, weapon-specific modeling, tactical rounds

### Future Vision
- **Professional Integration**: Real CS:GO team data and strategy modeling
- **Academic Research Platform**: Tools for economic and game theory research  
- **Tournament Simulation**: Complete esports ecosystem modeling
- **Multi-Agent Reinforcement Learning**: Self-improving strategies through AI

For detailed roadmap and implementation suggestions, see `PROJECT_ROADMAP.md`.

### Development Guidelines
- Follow Go best practices for concurrency
- Ensure thread-safety for any shared state
- Add benchmarks for performance-critical code
- Include error handling and timeout mechanisms
- Document any new command-line flags or configuration options

### Testing Large Simulations
Before submitting changes that affect parallel processing:
```bash
# Test with various scales
go run ./cmd -n 1000 -c 4    # Small scale
go run ./cmd -n 10000 -c 8   # Medium scale  
go run ./cmd -n 50000 -c 16  # Large scale (if hardware permits)
```

## License

This project is available as open source under the terms of the MIT License.

---

## Quick Start Example

```bash
# Clone and setup
git clone https://github.com/petz15/csgo_abm.git
cd csgo_abm
go mod tidy

# Run a quick performance test
go run ./cmd -n 1000 -c 4 -t1 all_in -t2 default_half

# Expected output:
# Starting 1000 simulations with 4 concurrent workers...
# Memory limit: 1000 MB
# Progress: 1000/1000 (100.0%) - Rate: 245.2 sims/sec - Elapsed: 4s
# 
# ============================================================
# SIMULATION RESULTS SUMMARY
# ============================================================
# Total Simulations: 1000
# Execution Time: 4s
# Processing Rate: 245.2 simulations/second
# Peak Memory Usage: 156 MB
# ...
```

For detailed parallel simulation information, see `PARALLEL_SIMULATION_GUIDE.md`.