package main

import (
	"csgo-economy-sim/internal/engine"
	"encoding/json"
	"fmt"
)

func main() {
	// Initialize the game engine

	//for test cases
	Team1Name := "Team A"
	Team1Strategy := "all_in"
	Team2Name := "Team B"
	Team2Strategy := "simple"
	gamerules := "default"

	game := engine.NewGame(Team1Name, Team1Strategy, Team2Name, Team2Strategy, gamerules)

	// Start the simulation
	game.Start()

	// Print all variables in a readable JSON format
	gameState := map[string]interface{}{
		"Team1Name":     Team1Name,
		"Team1Strategy": Team1Strategy,
		"Team2Name":     Team2Name,
		"Team2Strategy": Team2Strategy,
		"gamerules":     gamerules,
		"Game":          game,
	}
	jsonBytes, err := json.MarshalIndent(gameState, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
	} else {
		fmt.Println(string(jsonBytes))
	}
}
