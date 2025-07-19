# CS:GO Agent-Based Economy Model (CSGO_ABM)

This project implements an agent-based model simulating the economy aspect of Counter-Strike: Global Offensive (CS:GO). The simulation focuses on how teams manage their finances, make investment decisions, and the impact of these decisions on game outcomes through probabilistic models.

## Project Structure

```
CSGO_ABM
├── cmd
│   ├── main.go               # Entry point with CLI argument handling
│   ├── game.go               # Individual game simulation logic
│   ├── simulation.go         # Parallel simulation management
│   └── results               # Directory for simulation results
├── internal
│   ├── engine
│   │   ├── game.go           # Core game logic and state management
│   │   ├── gamerules.go      # Game rules and configuration
│   │   ├── probabilities.go  # Probability distributions and CSF functions
│   │   ├── round.go          # Individual round simulation
│   │   ├── strategymanager.go# Strategy selection and implementation
│   │   └── team.go           # Team state and behavior
│   └── models
│       ├── all_in.go         # "All-in" spending strategy
│       └── default_half.go   # "Default-half" spending strategy
├── output
│   └── results.go            # Functions for exporting simulation results
├── go.mod                    # Module definition
├── go.sum                    # Dependency checksums
└── README.md                 # Project documentation
```

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
- JSON export of simulation results with unique game identifiers

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

### Running Simulations

You can run simulations with various command-line arguments:

```bash
# Run a single simulation
go run ./cmd

# Run 100 simulations using 8 cores 
go run ./cmd -n 100 -c 8

# Use specific team strategies
go run ./cmd -t1 all_in -t2 default_half
```

#### Command-Line Arguments

- `-n, --num <number>`: Number of simulations to run (default: 1)
- `-c, --cores <number>`: Number of concurrent simulations (default: number of CPU cores)
- `-t1, --team1 <strategy>`: Team 1 strategy (default: all_in)
- `-t2, --team2 <strategy>`: Team 2 strategy (default: default_half)
- `-h, --help`: Print help message

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

### Analyzing Results

Results are exported as JSON files with unique identifiers:
- Format: `YYYYMMDD_HHMMSS_hostname_randomID.json`
- Contains complete game state including rounds, team performance, and outcome

When running multiple simulations in parallel, the system:
1. Creates individual result files for each simulation
2. Moves all results to a timestamped directory
3. Generates a summary file (`simulation_summary.json`) with aggregate statistics including:
   - Win rates for each team
   - Score distribution
   - Average number of rounds played
   - Overtime rate

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

- Improved strategy manager with more sophisticated decision-making models
- Enhanced probability models for more realistic outcome simulation
- Advanced statistical analysis tools for simulation results
- Visualization tools for simulation results
- Additional game mechanics such as eco rounds and force buys

## Contributing

Contributions are welcome! Areas that could benefit from improvements:
- Additional strategies
- Refined probability distributions
- Enhanced simulation parameters
- Output analysis tools

## License

This project is available as open source under the terms of the MIT License.