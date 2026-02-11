package engine

import (
	"encoding/json"
	"fmt"
	"os"
)

type GameRules struct {
	DefaultEquipment                float64    `json:"defaultEquipment"`              // Default equipment
	OTFunds                         float64    `json:"otFunds"`                       // Overtime funds
	OTEquipment                     float64    `json:"otEquipment"`                   // Overtime equipment
	StartingFunds                   float64    `json:"startingFunds"`                 // Starting funds for teams
	HalfLength                      int        `json:"halfLength"`                    // Length of a half in rounds
	OTHalfLength                    int        `json:"otHalfLength"`                  // Length of overtime half in rounds
	MaxFunds                        float64    `json:"maxFunds"`                      // Maximum funds allowed for a team
	LossBonusCalc                   bool       `json:"lossBonusCalc"`                 // true: loss bonus reduced 1 after each win, false: resets after each win
	WithSaves                       bool       `json:"withSaves"`                     // true: teams can save weapons between rounds
	LossBonus                       []float64  `json:"lossBonus"`                     // Custom loss bonus per round (if empty, use default logic)
	RoundOutcomeReward              [4]float64 `json:"roundOutcomeReward"`            // Custom rewards for round outcomes
	EliminationReward               float64    `json:"eliminationReward"`             // Reward for eliminating a opponent
	BombplantRewardall              float64    `json:"bombplantRewardall"`            // Reward for planting the bomb for all players
	BombplantReward                 float64    `json:"bombplantReward"`               // Reward for planting the bomb
	BombdefuseReward                float64    `json:"bombdefuseReward"`              // Reward for defusing the bomb
	AdditionalReward_CT_Elimination float64    `json:"additionalCTEliminationReward"` // Additional reward for CT team for eliminations
	AdditionalReward_T_Elimination  float64    `json:"additionalTEliminationReward"`  // Additional reward for T team for eliminations
	Custom_CSF_r_value              float64    `json:"customRValue"`                  // Custom r value the CSF default is to use from the probabilities.json
}

// getDefaultRules returns the default game rules configuration
func getDefaultRules() GameRules {
	return GameRules{
		DefaultEquipment:                200,
		OTFunds:                         10000,
		OTEquipment:                     200,
		StartingFunds:                   800,
		HalfLength:                      15,
		OTHalfLength:                    3,           // Default value for Overtime half length
		MaxFunds:                        (16000 * 5), // Default value for Maximum funds
		LossBonusCalc:                   true,
		WithSaves:                       true,
		LossBonus:                       []float64{1400, 1900, 2400, 2900, 3400},
		RoundOutcomeReward:              [4]float64{3500, 3250, 3500, 3250}, // Default rewards for Roundoutcome 1-4
		EliminationReward:               300,                                //technically this is dependent on the weapon used, but we use a flat value for simplicity as there are no weapons
		BombplantRewardall:              800,                                //in cs2 this is 600 but in csgo its 800
		BombplantReward:                 300,
		BombdefuseReward:                300,
		AdditionalReward_CT_Elimination: 0, //introduced in cs2, 50 per elimination for all CT players
		AdditionalReward_T_Elimination:  0,
		Custom_CSF_r_value:              -1, // Default value of 0 means "use default from probabilities.json"
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

	if rules.OTEquipment < 0 {
		fmt.Println("Overtime equipment must be non-negative")
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

	if rules.LossBonusCalc != true && rules.LossBonusCalc != false {
		fmt.Println("Loss bonus calculation flag must be true or false")
		return false
	}

	if rules.WithSaves != true && rules.WithSaves != false {
		fmt.Println("With saves flag must be true or false")
		return false
	}

	// Validate reward values are non-negative
	if rules.EliminationReward < 0 {
		fmt.Println("Elimination reward must be non-negative")
		return false
	}
	if rules.BombplantRewardall < 0 {
		fmt.Println("Bomb plant reward (all) must be non-negative")
		return false
	}
	if rules.BombplantReward < 0 {
		fmt.Println("Bomb plant reward must be non-negative")
		return false
	}
	if rules.BombdefuseReward < 0 {
		fmt.Println("Bomb defuse reward must be non-negative")
		return false
	}
	if rules.AdditionalReward_CT_Elimination < 0 {
		fmt.Println("Additional CT elimination reward must be non-negative")
		return false
	}
	if rules.AdditionalReward_T_Elimination < 0 {
		fmt.Println("Additional T elimination reward must be non-negative")
		return false
	}

	// Validate RoundOutcomeReward values
	for i, reward := range rules.RoundOutcomeReward {
		if reward < 0 {
			fmt.Printf("Round outcome reward [%d] must be non-negative\n", i)
			return false
		}
	}

	// Validate LossBonus values if provided
	for i, bonus := range rules.LossBonus {
		if bonus < 0 {
			fmt.Printf("Loss bonus [%d] must be non-negative\n", i)
			return false
		}
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
		if jsonRules.OTEquipment > 0 {
			candidateRules.OTEquipment = jsonRules.OTEquipment
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
		if jsonRules.OTHalfLength > 0 {
			candidateRules.OTHalfLength = jsonRules.OTHalfLength
		}
		if jsonRules.MaxFunds > 0 {
			candidateRules.MaxFunds = jsonRules.MaxFunds
		}
		if jsonRules.LossBonusCalc == true || jsonRules.LossBonusCalc == false {
			candidateRules.LossBonusCalc = jsonRules.LossBonusCalc
		}
		if jsonRules.WithSaves == true || jsonRules.WithSaves == false {
			candidateRules.WithSaves = jsonRules.WithSaves
		}

		// Load reward values (can be 0, so check if they were explicitly set)
		// For rewards, we allow 0 values, so we only skip if not present in JSON
		if jsonRules.EliminationReward >= 0 {
			candidateRules.EliminationReward = jsonRules.EliminationReward
		}
		if jsonRules.BombplantRewardall >= 0 {
			candidateRules.BombplantRewardall = jsonRules.BombplantRewardall
		}
		if jsonRules.BombplantReward >= 0 {
			candidateRules.BombplantReward = jsonRules.BombplantReward
		}
		if jsonRules.BombdefuseReward >= 0 {
			candidateRules.BombdefuseReward = jsonRules.BombdefuseReward
		}
		if jsonRules.AdditionalReward_CT_Elimination >= 0 {
			candidateRules.AdditionalReward_CT_Elimination = jsonRules.AdditionalReward_CT_Elimination
		}
		if jsonRules.AdditionalReward_T_Elimination >= 0 {
			candidateRules.AdditionalReward_T_Elimination = jsonRules.AdditionalReward_T_Elimination
		}

		if jsonRules.Custom_CSF_r_value >= 0 {
			candidateRules.Custom_CSF_r_value = jsonRules.Custom_CSF_r_value
		}

		// Load RoundOutcomeReward array (check if any values are set)
		hasRoundRewards := false
		for _, val := range jsonRules.RoundOutcomeReward {
			if val != 0 {
				hasRoundRewards = true
				break
			}
		}
		if hasRoundRewards {
			candidateRules.RoundOutcomeReward = jsonRules.RoundOutcomeReward
		}

		// Load LossBonus array (only if not empty)
		if len(jsonRules.LossBonus) > 0 {
			candidateRules.LossBonus = jsonRules.LossBonus
		}

		// Load Custom_CSF_r_value (can be 0 to use default from probabilities.json)
		// Negative values are invalid, but 0 means "use default"
		if jsonRules.Custom_CSF_r_value >= 0 {
			candidateRules.Custom_CSF_r_value = jsonRules.Custom_CSF_r_value
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
