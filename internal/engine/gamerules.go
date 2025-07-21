package engine

import (
	"encoding/json"
	"fmt"
	"os"
)

type GameRules struct {
	DefaultEquipment float64 `json:"defaultEquipment"` // Default equipment cost
	OTFunds          float64 `json:"otFunds"`          // Overtime funds
	StartingFunds    float64 `json:"startingFunds"`    // Starting funds for teams
	HalfLength       int     `json:"halfLength"`       // Length of a half in rounds
	CSF_r            float64 `json:"csfR"`             // Contest Success Function parameter
	OTHalfLength     int     `json:"otHalfLength"`     // Length of overtime half in rounds
	MaxFunds         float64 `json:"maxFunds"`         // Maximum funds allowed for a team
}

// getDefaultRules returns the default game rules configuration
func getDefaultRules() GameRules {
	return GameRules{
		DefaultEquipment: 200,
		OTFunds:          10000,
		StartingFunds:    800,
		HalfLength:       15,
		CSF_r:            0.8,    // Default value for CSF_r
		OTHalfLength:     3,      // Default value for Overtime half length
		MaxFunds:         999999, // Default value for Maximum funds
	}
}

// validateGameRulesStrict performs strict validation - returns false if any value is invalid
func validateGameRulesStrict(rules GameRules) bool {
	// Validate economic values are positive
	if rules.DefaultEquipment < 0 {
		fmt.Println("Default equipment must be non-negative")
		return false
	}
	if rules.StartingFunds < 0 {
		fmt.Println("Starting funds must be non-negative")
		return false
	}
	if rules.OTFunds < 0 {
		fmt.Println("Overtime funds must be non-negative")
		return false
	}

	if rules.MaxFunds < 0 {
		fmt.Println("Maximum funds must be non-negative")
		return false
	}

	// Validate round counts are reasonable
	if rules.HalfLength <= 0 {
		fmt.Println("Half length must be positive")
		return false
	}
	if rules.OTHalfLength <= 0 {
		fmt.Println("Overtime half length must be positive")
		return false
	}

	// Validate Contest Success Function parameter
	if rules.CSF_r < 0 {
		fmt.Println("Contest Success Function (r) must be non-negative")
		return false
	}

	return true
}

func NewGameRules(pathtoFile string) (GameRules, bool) {
	// Start with default rules
	rules := getDefaultRules()

	if pathtoFile != "" && pathtoFile != "default" {
		// Attempt to read game rules from a JSON file
		file, err := os.Open(pathtoFile)
		if err != nil {
			fmt.Printf("Warning: Could not open game rules file '%s': %v. Using defaults.\n", pathtoFile, err)
			return rules, false
		}
		defer file.Close()

		// Create a temporary rules struct for JSON parsing
		var jsonRules GameRules
		decoder := json.NewDecoder(file)

		if err := decoder.Decode(&jsonRules); err != nil {
			fmt.Printf("Warning: Could not parse game rules file '%s': %v. Using defaults.\n", pathtoFile, err)
			return rules, false
		}

		// Merge JSON values with defaults (only overwrite non-zero values)
		candidateRules := rules // Start with defaults

		if jsonRules.DefaultEquipment > 0 {
			candidateRules.DefaultEquipment = jsonRules.DefaultEquipment
		}
		if jsonRules.OTFunds > 0 {
			candidateRules.OTFunds = jsonRules.OTFunds
		}
		if jsonRules.StartingFunds > 0 {
			candidateRules.StartingFunds = jsonRules.StartingFunds
		}
		if jsonRules.HalfLength > 0 {
			candidateRules.HalfLength = jsonRules.HalfLength
		}
		if jsonRules.CSF_r > 0 {
			candidateRules.CSF_r = jsonRules.CSF_r
		}
		if jsonRules.OTHalfLength > 0 {
			candidateRules.OTHalfLength = jsonRules.OTHalfLength
		}
		if jsonRules.MaxFunds > 0 {
			candidateRules.MaxFunds = jsonRules.MaxFunds
		}

		// Strict validation - if ANY value fails, use all defaults
		if !validateGameRulesStrict(candidateRules) {
			fmt.Printf("Warning: Game rules in '%s' failed validation. Using all default values.\n", pathtoFile)
			return rules, false // Return original defaults
		}

		// All validations passed, use the candidate rules
		rules = candidateRules
		return rules, true
	}

	return rules, false
}
