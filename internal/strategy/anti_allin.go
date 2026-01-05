package strategy

import "math"

func InvestDecisionMaking_anti_allin(ctx StrategyContext_simple) float64 {

	if ctx.IsLastRoundHalf {
		return ctx.Funds
	} else if ctx.IsPistolRound {
		return ctx.Funds
	} else if ctx.IsAfterPistol && ctx.ConsecutiveLosses > 0 {
		return 0
	} else if ctx.GameRules_strategy.HalfLength-ctx.OpponentScore == 1 && !ctx.IsOvertime {
		return ctx.Funds * 0.8
	} else if ctx.GameRules_strategy.HalfLength-ctx.OpponentScore == 0 && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.GameRules_strategy.HalfLength-ctx.OwnScore == 1 && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.GameRules_strategy.HalfLength-ctx.OwnScore == 0 && !ctx.IsOvertime {
		return ctx.Funds
	}

	if ctx.EnemySurvivors < 1 && ctx.ConsecutiveWins <= 2 {
		//max money the enemy can spend to win the round is 1900 * 5 + default equipment *5
		max_enemy_investment := 1900*5 + ctx.GameRules_strategy.DefaultEquipment*5
		investment := (max_enemy_investment * 2) - ctx.Equipment
		return math.Min(investment, ctx.Funds)
	} else if ctx.ConsecutiveWins == 1 {
		return ctx.Funds
	} else if ctx.ConsecutiveWins >= 2 {
		ratio := 1 - (float64(ctx.ConsecutiveWins) * 0.1)
		ratio = math.Max(ratio, 0.4) //set a minimum ratio of 0.4
		return ctx.Funds * ratio
	}

	return ctx.Funds * 0.9
}
