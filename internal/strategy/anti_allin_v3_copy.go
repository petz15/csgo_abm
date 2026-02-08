package strategy

import (
	"encoding/json"
	"math"
	"os"
	"sync"
)

// Configuration struct for anti_allin_v3 strategy parameters
type AntiAllinV3Config_copy struct {
	PressingRatio     float64 `json:"pressing_ratio"`
	OverturnThreshold float64 `json:"overturn_threshold"`
	OverturnRatio     float64 `json:"overturn_ratio"`
}

var (
	antiAllinV3Config_copy     AntiAllinV3Config_copy
	antiAllinV3ConfigOnce_copy sync.Once
)

func InvestDecisionMaking_anti_allin_v3_copy(ctx StrategyContext_simple) float64 {

	// anti_allin_v3, invests all in the beginning and end of halves/overtime. In between it tries to build up wealth
	Score_to_Win := ctx.GameRules_strategy.HalfLength + (ctx.GameRules_strategy.OTHalfLength * (ctx.OvertimeAmount)) + 1

	// Load configuration once
	config := getAntiAllinV3Config()
	pressing_ratio := config.PressingRatio
	overturn_threshold := config.OverturnThreshold
	overturn_ratio := config.OverturnRatio

	//always go all in if last round of half or close scores
	if ctx.IsLastRoundHalf {
		return ctx.Funds
	} else if Score_to_Win-ctx.OpponentScore == 2 || Score_to_Win-ctx.OwnScore == 2 {
		return ctx.Funds
	} else if Score_to_Win-ctx.OpponentScore == 1 || Score_to_Win-ctx.OwnScore == 1 {
		return ctx.Funds
	}

	if !ctx.IsOvertime {
		if ctx.IsFirstRoundHalf {
			return ctx.Funds
		} else if ctx.IsSecondRoundHalf && ctx.ConsecutiveLosses < 1 {
			//also go all in if round after pistol and it is not overtime and we won pistol
			return ctx.Funds
		}
	} else if ctx.IsFirstRoundHalf && ctx.GameRules_strategy.OTFunds > avgArray(ctx.GameRules_strategy.RoundOutcomeReward[:])*2.5 {
		// try to bait all in strat to spend everything in the first round of overtime
		//hoping in the consecutive rounds they will have considerably less funds
		// but only if OT funds are high enough (approx. 2x avg round reward)
		return 0.0
	} else if ctx.IsSecondRoundHalf && ctx.ConsecutiveLosses == 1 {
		//if we lost the first OT round and it is the second round
		//assume that about 80% of the opponents OT funds remained per survivor + avg round reward
		approx_funds := (ctx.GameRules_strategy.RoundOutcomeReward[ctx.RoundEndReason-1] + ctx.GameRules_strategy.OTEquipment) * 5
		approx_saved_eq := float64(ctx.EnemySurvivors) * ((ctx.GameRules_strategy.OTEquipment + ctx.GameRules_strategy.OTEquipment) / 5.0) * 0.9
		additional_rewards := (5 - float64(ctx.OwnSurvivors)) * ctx.GameRules_strategy.EliminationReward

		if !ctx.GameRules_strategy.WithSaves {
			approx_saved_eq = 0.0
		}

		approx_allin_investment := approx_funds + approx_saved_eq + additional_rewards

		return math.Min((ctx.Funds + ctx.Equipment), approx_allin_investment*pressing_ratio)
	}

	//Regular time round decision making
	if ctx.ConsecutiveLosses < 1 {
		//try to invest accoring to the pressing factor
		approx_funds := (ctx.GameRules_strategy.LossBonus[ctx.LossBonusLevel_opponent] + ctx.GameRules_strategy.DefaultEquipment) * 5
		approx_saved_eq := 0.0
		if ctx.ConsecutiveWins == 1 {
			approx_saved_eq += float64(ctx.EnemySurvivors) * avgArray(ctx.GameRules_strategy.RoundOutcomeReward[:])
		} else {
			approx_saved_eq += float64(ctx.EnemySurvivors) * ctx.GameRules_strategy.LossBonus[int(math.Max(0, float64(ctx.LossBonusLevel_opponent-1)))]
		}
		additional_rewards := (5 - float64(ctx.OwnSurvivors)) * ctx.GameRules_strategy.EliminationReward

		if !ctx.GameRules_strategy.WithSaves {
			approx_saved_eq = 0.0
		}

		approx_allin_investment := approx_funds + approx_saved_eq + additional_rewards

		return math.Min((ctx.Funds + ctx.Equipment), approx_allin_investment*pressing_ratio)
	}

	approx_funds := (ctx.GameRules_strategy.RoundOutcomeReward[ctx.RoundEndReason-1] + ctx.GameRules_strategy.DefaultEquipment) * 5
	approx_saved_eq := float64(ctx.EnemySurvivors) * avgArray(ctx.GameRules_strategy.RoundOutcomeReward[:])
	additional_rewards := (5 - float64(ctx.OwnSurvivors)) * ctx.GameRules_strategy.EliminationReward

	if !ctx.GameRules_strategy.WithSaves {
		approx_saved_eq = 0.0
	}

	approx_allin_investment := approx_funds + approx_saved_eq + additional_rewards

	//if we have enough funds to challenge all in, do it
	if ctx.Funds+ctx.Equipment >= approx_allin_investment*overturn_threshold {
		return math.Min((ctx.Funds + ctx.Equipment), approx_allin_investment*overturn_ratio)
	}

	//otherwise save until enough funds are built up
	return 0.0
}

func getAntiAllinV3Config_copy() AntiAllinV3Config_copy {
	antiAllinV3ConfigOnce_copy.Do(func() {
		antiAllinV3Config_copy = loadAntiAllinV3Config_copy()
	})
	return antiAllinV3Config_copy
}

func loadAntiAllinV3Config_copy() AntiAllinV3Config_copy {
	// Default fallback values
	defaultConfig := AntiAllinV3Config_copy{
		PressingRatio:     1.15,
		OverturnThreshold: 0.3,
		OverturnRatio:     0.8,
	}

	// Try to load from JSON file
	configPaths := []string{
		"configs/anti_allin_v3_copy.json",
		"config/anti_allin_v3_copy.json",
		"anti_allin_v3_copy.json",
	}

	for _, path := range configPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue // Try next path
		}

		var config AntiAllinV3Config_copy
		if err := json.Unmarshal(data, &config); err != nil {
			continue // Try next path or use defaults
		}

		// Validate loaded values (optional)
		if config.PressingRatio > 0 && config.OverturnThreshold > 0 && config.OverturnRatio > 0 {
			return config
		}
	}

	// Return defaults if no valid config found
	return defaultConfig
}
