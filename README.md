# CS:GO Agent-Based Economy Model (CSGO_ABM)

This project implements an agent-based model simulating the economy aspect of Counter-Strike: Global Offensive (CS:GO). The simulation focuses on how teams manage their finances, make investment decisions, and the impact of these decisions on game outcomes through probabilistic models.

## Project Structure

```
CSGO_ABM
├── cmd
│   └── main.go               # Entry point of the application
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

The simulation creates a game with two teams using specified strategies:

```go
// From cmd/main.go
Team1Name := "Team A"
Team1Strategy := "all_in"
Team2Name := "Team B"
Team2Strategy := "default_half"
gamerules := "default"

game := engine.NewGame(ID, Team1Name, Team1Strategy, Team2Name, Team2Strategy, gamerules)
game.Start()
```

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

## Future Developments

- Improved strategy manager with more sophisticated decision-making models
- Enhanced probability models for more realistic outcome simulation
- Support for batch simulations and comparative analysis
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