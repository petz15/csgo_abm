package main

import (
	"csgo-economy-sim/internal/engine"
	"csgo-economy-sim/internal/util"
	"fmt"
)

func main() {
	// Initialize the game engine
	game := engine.NewGame()

	// Set up teams
	game.SetupTeams()

	// Start the simulation
	results := game.RunSimulation()

	// Export results
	err := util.ExportResults(results)
	if err != nil {
		fmt.Println("Error exporting results:", err)
	} else {
		fmt.Println("Simulation results exported successfully.")
	}
}
