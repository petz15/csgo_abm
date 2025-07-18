package engine

import (
	"encoding/json"
	"os"
)

type GameRules struct {
	DefaultEquipment float64 // Default equipment cost
	OTFunds          float64 // Overtime funds
	StartingFunds    float64 // Starting funds for teams
	HalfLength       int     // Length of a half in rounds
	CSF_r            float64 // Contest Success Function parameter
}

func NewGameRules(pathtoFile string) GameRules {
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
		DefaultEquipment: 200,
		OTFunds:          10000,
		StartingFunds:    800,
		HalfLength:       15,
		CSF_r:            0.5, // Default value for CSF_r
	}
}
