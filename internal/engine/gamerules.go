package engine

import (
	"encoding/json"
	"os"
)

type GameRules struct {
	defaultEquipment int     // Default equipment cost
	otEquipment      int     // Overtime equipment cost
	otFunds          int     // Overtime funds
	startingFunds    int     // Starting funds for teams
	CSF_r            float32 // Contest Success Function parameter
}

func initGameRules(pathtoFile string) GameRules {
	if pathtoFile != "" {
		// Attempt to read game rules from a JSON file
		file, err := os.Open(pathtoFile)
		if err == nil {
			defer file.Close()
			decoder := json.NewDecoder(file)
			var rules GameRules
			if err := decoder.Decode(&rules); err == nil {
				return rules
			}
		}
		// If file not found or parse error, fallback to defaults
	}
	return GameRules{
		defaultEquipment: 200,
		otEquipment:      200,
		otFunds:          10000,
		startingFunds:    800,
		CSF_r:            0.5, // Default value for CSF_r
	}
}
