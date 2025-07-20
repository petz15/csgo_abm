# CS:GO Agent-Based Economy Model (CSGO_ABM)

This project implements an agent-based model simulating the economy aspect of Counter-Strike: Global Offensive (CS:GO). The simulation focuses on how teams manage their finances, make investment decisions, and the impact of these decisions on game outcomes through probabilistic models.

## Project Structure

```
CSGO_ABM
├── cmd
│   ├── main.go                    # Entry point with CLI argument handling
│   ├── gamehandler.go            # Game handling and orchestration
│   ├── simulation.go             # High-performance parallel simulation management
│   └── simulation_sequential.go  # Sequential simulation for comparison/debugging
├── internal
│   ├── engine
│   │   ├── game.go               # Core game logic and state management
│   │   ├── gamerules.go          # Game rules and configuration
│   │   ├── probabilities.go      # Probability distributions and CSF functions
│   │   ├── round.go              # Individual round simulation
│   │   ├── strategymanager.go    # Strategy selection and implementation
│   │   └── team.go               # Team state and behavior
│   └── models
│       ├── all_in.go             # "All-in" spending strategy
│       └── default_half.go       # "Default-half" spending strategy
├── output
│   └── results.go                # Functions for exporting simulation results
├── results_*/                    # Timestamped directories with simulation outputs
├── go.mod                        # Module definition
├── go.sum                        # Dependency checksums
├── PARALLEL_SIMULATION_GUIDE.md  # Detailed guide for parallel simulations
└── README.md                     # Project documentation
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
- Contest Success Function (CSF) for probabilistic outcome determination
- Equipment value tracking and survival calculations
- Complete match simulation including:
  - Half-time side swapping after 15 rounds
  - Overtime mechanics when tied at 15-15
  - Loss bonus calculation based on consecutive losses
  - Bomb plant mechanics and related bonuses

### 📊 **Advanced Analytics**
- Real-time progress monitoring with performance metrics
- Comprehensive statistical analysis including win rates, score distributions, and overtime rates
- Memory usage tracking and optimization
- JSON export with detailed simulation metadata
- Aggregated results with summary statistics

### 🔧 **Concurrency Improvements (Latest Version)**
- **Fixed Memory Leaks**: Resolved unbounded map growth and goroutine leaks
- **Prevented Hanging**: Proper shutdown sequences and timeout handling
- **Race Condition Fixes**: Atomic operations for all shared state
- **Context Propagation**: Proper cancellation signal handling throughout the system
- **Resource Cleanup**: Guaranteed cleanup of monitoring goroutines and resources

## Probability Models

The simulation uses several probability models to determine outcomes:

1. **Contest Success Function (CSF)**: Calculates the probability of winning based on team expenditures
   - Uses Tullock contest model with adjustable parameter r
   - Higher r values make outcomes more deterministic

2. **CSF Normal Distribution**: Generates values from a normal distribution with:
   - Mean positioned according to the CSF probability
   - Standard deviation set to 1/4 of the output range
   - Configurable output range for different game aspects

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
   ```
   go run cmd/main.go
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

#### Simulation Output

Each parallel run creates a timestamped directory containing:
- `simulation_summary.json`: Aggregate statistics and performance metrics
- Individual simulation results (if enabled)
- Memory usage and performance data

Example summary output:
```json
{
  "total_simulations": 10000,
  "completed_simulations": 10000,
  "team1_wins": 5234,
  "team2_wins": 4766,
  "team1_win_rate": 52.34,
  "team2_win_rate": 47.66,
  "average_rounds": 28.7,
  "overtime_rate": 23.4,
  "execution_time": "45.2s",
  "simulations_per_second": 221.2,
  "peak_memory_usage_mb": 1847,
  "total_gc_runs": 12
}
```

#### Command-Line Arguments

- `-n, --num <number>`: Number of simulations to run (default: 1)
- `-c, --cores <number>`: Number of concurrent workers (default: number of CPU cores)
- `-m, --memory <MB>`: Memory limit in MB before forcing garbage collection (default: 1000)
- `-t1, --team1 <strategy>`: Team 1 strategy (default: all_in)
- `-t2, --team2 <strategy>`: Team 2 strategy (default: default_half)
- `-r, --rules <ruleset>`: Game rules to use (default: standard)
- `-h, --help`: Print help message

### Performance Benchmarks

Based on testing with the improved concurrency implementation:

| Simulations | Workers | Memory Limit | Avg. Time | Rate (sims/sec) | Peak Memory |
|-------------|---------|--------------|-----------|-----------------|-------------|
| 1,000       | 4       | 1000MB       | 4.2s      | 238/sec         | 145MB       |
| 10,000      | 8       | 1000MB       | 42s       | 238/sec         | 520MB       |
| 100,000     | 16      | 2000MB       | 7.1min    | 234/sec         | 1.8GB       |

*Results may vary based on hardware specifications and strategy complexity.*

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

### Concurrency Architecture

The simulation uses a robust worker pool pattern with the following components:

1. **Worker Pool**: Manages a fixed number of goroutines for simulation processing
2. **Job Queue**: Buffered channel for distributing simulation tasks
3. **Result Collection**: Dedicated goroutine for aggregating simulation results
4. **Memory Monitoring**: Background goroutine tracking memory usage and triggering GC
5. **Progress Reporting**: Real-time status updates during long-running simulations

### Memory Management

- **Bounded Data Structures**: Score distribution limited to 1000 unique entries
- **Automatic GC**: Triggered when memory usage exceeds configured limits
- **Atomic Operations**: Race-condition-free statistics updates
- **Resource Cleanup**: Guaranteed cleanup of all goroutines and channels

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

## Future Developments

### Planned Enhancements
- **GPU Acceleration**: CUDA support for massive parallel simulations
- **Distributed Computing**: Multi-node simulation processing
- **Machine Learning Integration**: AI-driven strategy optimization
- **Real-time Visualization**: Live simulation monitoring dashboard
- **Database Integration**: Persistent storage for large-scale result analysis

### Strategy Development
- **Advanced Economic Models**: More sophisticated spending algorithms
- **Adaptive Strategies**: Dynamic strategy adjustment based on game state
- **Meta-game Analysis**: Strategy effectiveness across different scenarios
- **Player Behavior Modeling**: Individual player decision-making simulation

### Performance Improvements
- **Memory Pool Reuse**: Reduce garbage collection overhead
- **Vectorized Operations**: SIMD optimizations for probability calculations
- **Streaming Results**: Real-time result processing for very large runs
- **Smart Caching**: Precomputed probability tables and strategy decisions

## Troubleshooting

### Common Issues

**Memory Usage Growth**: 
- Increase memory limit with `-m` flag
- Reduce number of concurrent workers
- Monitor with built-in memory tracking

**Hanging Simulations**:
- Check for infinite loops in custom strategies
- Verify timeout settings (default: 5 minutes per simulation)
- Use sequential mode for debugging: `simulation_sequential.go`

**Performance Issues**:
- Optimize worker count (typically 1-2x CPU cores)
- Ensure adequate system memory
- Consider reducing simulation complexity for large runs

### Debug Mode

For debugging individual simulations:
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
- **Performance Optimization**: Improve simulation speed and memory efficiency  
- **Testing**: Add comprehensive test suites for parallel simulation reliability
- **Documentation**: Improve code documentation and usage examples
- **Analysis Tools**: Create visualization and statistical analysis utilities
- **Game Mechanics**: Add new CS:GO features (weapon mechanics, map-specific strategies)

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