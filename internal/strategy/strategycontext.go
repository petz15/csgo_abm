package strategy

import "math/rand"

// StrategyContext holds all relevant information for economic decision making
type StrategyContext_simple struct {
	Funds              float64
	Equipment          float64
	PlayersAlive       int     // Number of players alive in the team
	RoundImportance    float64 // Importance of the current round
	CurrentRound       int
	OpponentScore      int
	OwnScore           int
	ConsecutiveLosses  int
	ConsecutiveWins    int
	LossBonusLevel     int  // Current loss bonus level
	Side               bool // true = CT, false = T
	IsPistolRound      bool
	IsLastRoundHalf    bool
	IsOvertime         bool
	OvertimeAmount     int // Number of overtime periods played
	IsAfterPistol      bool
	OwnSurvivors       int  // Number of survivors in the previous RoundEndReason
	EnemySurvivors     int  // Number of enemy survivors in the previous RoundEndReason
	RoundEndReason     int  // Reason for the end of the last round
	Is_BombPlanted     bool // Whether the bomb was planted in the last round
	RNG                *rand.Rand
	GameRules_strategy GameRules_strategymanager
}

type GameRules_strategymanager struct {
	DefaultEquipment                float64    `json:"defaultEquipment"`              // Default equipment cost
	OTFunds                         float64    `json:"otFunds"`                       // Overtime funds
	OTEquipment                     float64    `json:"otEquipment"`                   // Overtime equipment cost
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
}

func empty() {
	// This function is intentionally left empty to ensure the package compiles correctly.
	// It serves as a placeholder for future enhancements or additional logic.

}

// ReLU activation function
func relu(x float64) float64 {
	if x > 0 {
		return x
	}
	return 0
}

func avgArray(arr []float64) float64 {
	if len(arr) == 0 {
		return 0
	} else {
		sum := 0.0
		for _, v := range arr {
			sum += v
		}
		return sum / float64(len(arr))
	}
}
