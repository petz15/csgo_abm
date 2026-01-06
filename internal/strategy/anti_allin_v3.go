package strategy

import (
	"encoding/json"
	"math"
	"os"
	"sync"
)

// Configuration struct for anti_allin_v3 strategy parameters
type AntiAllinV3Config struct {
	PressingRatio     float64 `json:"pressing_ratio"`
	OverturnThreshold float64 `json:"overturn_threshold"`
	OverturnRatio     float64 `json:"overturn_ratio"`
}

var (
	antiAllinV3Config     AntiAllinV3Config
	antiAllinV3ConfigOnce sync.Once
)

func InvestDecisionMaking_anti_allin_v3(ctx StrategyContext_simple) float64 {

	// anti_allin_v3, invests all in the beginning and end of halves/overtime. In between it tries to build up wealth
	Score_to_Win := ctx.GameRules_strategy.HalfLength + (ctx.GameRules_strategy.OTHalfLength * (ctx.OvertimeAmount)) + 1

	// Load configuration once
	config := getAntiAllinV3Config()
	pressing_ratio := config.PressingRatio
	overturn_threshold := config.OverturnThreshold
	overturn_ratio := config.OverturnRatio

	//always go all in if pistol round (in not OT rounds) or last round of half
	if ctx.IsLastRoundHalf {
		return ctx.Funds
	} else if ctx.IsFirstRoundHalf && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.IsLastRoundHalf && !ctx.IsOvertime && ctx.ConsecutiveLosses < 1 {
		//also go all in if round after pistol and it is not overtime and we won pistol
		return ctx.Funds
	} else if Score_to_Win-ctx.OpponentScore == 2 || Score_to_Win-ctx.OwnScore == 2 {
		return ctx.Funds
	} else if Score_to_Win-ctx.OpponentScore == 1 || Score_to_Win-ctx.OwnScore == 1 {
		return ctx.Funds
	}

	if !ctx.IsOvertime {
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

	} else {

		//if somebody is about to win, go all in
		if Score_to_Win-ctx.OpponentScore == 1 || Score_to_Win-ctx.OwnScore == 1 {
			return ctx.Funds
		} else {
			//try to divide the funds evenly over the OT rounds
			return ctx.Funds / float64(ctx.GameRules_strategy.OTHalfLength)
		}
	}

	return 0.0
}

func getAntiAllinV3Config() AntiAllinV3Config {
	antiAllinV3ConfigOnce.Do(func() {
		antiAllinV3Config = loadAntiAllinV3Config()
	})
	return antiAllinV3Config
}

func loadAntiAllinV3Config() AntiAllinV3Config {
	// Default fallback values
	defaultConfig := AntiAllinV3Config{
		PressingRatio:     1.15,
		OverturnThreshold: 0.3,
		OverturnRatio:     0.8,
	}

	// Try to load from JSON file
	configPaths := []string{
		"configs/anti_allin_v3.json",
		"config/anti_allin_v3.json",
		"anti_allin_v3.json",
	}

	for _, path := range configPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue // Try next path
		}

		var config AntiAllinV3Config
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
