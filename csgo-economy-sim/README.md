# CS:GO Economy Simulation

This project simulates the economy aspect of the popular game Counter-Strike: Global Offensive (CS:GO). The simulation focuses on how teams manage their finances, make investment decisions, and the impact of these decisions on the game's outcome.

## Project Structure

```
csgo-economy-sim
├── cmd
│   └── main.go               # Entry point of the application
├── internal
│   ├── engine
│   │   ├── economy.go        # Manages financial aspects of the game
│   │   ├── game.go           # Manages overall game state
│   │   ├── round.go          # Handles individual rounds
│   │   ├── simulation.go      # Runs the game simulation
│   │   └── team.go           # Represents teams and their strategies
│   ├── models
│   │   ├── player.go         # Represents individual players
│   │   ├── equipment.go      # Represents equipment available for purchase
│   │   └── strategy.go       # Represents spending strategies
│   └── util
│       ├── csf.go           # Contest Success Function logic
│       └── random.go        # Utility functions for randomness
├── strategies
│   ├── aggressive.json       # Aggressive spending strategy
│   ├── conservative.json     # Conservative spending strategy
│   └── balanced.json         # Balanced spending strategy
├── output
│   └── results.go            # Functions for exporting simulation results
├── go.mod                    # Module definition
├── go.sum                    # Dependency checksums
└── README.md                 # Project documentation
```

## Features

- Random assignment of teams to Counter-Terrorist (CT) and Terrorist (T) sides.
- Starting money allocation for both teams.
- Teams can decide how much to invest in each buy phase.
- Round outcomes determined by a Contest Success Function (CSF) based on team expenditures.
- Equipment, surviving players, and bomb planting effects are influenced by the CSF.
- Half-time after 15 rounds, with overtime rules if the score reaches 15:15.
- Import spending strategies from JSON files.
- Export results of the simulation for analysis.

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/csgo-economy-sim.git
   cd csgo-economy-sim
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

- Modify or create new spending strategies in the `strategies` directory.
- Analyze the results generated in the `output` directory after running the simulation.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements.

## License

This project is licensed under the MIT License. See the LICENSE file for details.