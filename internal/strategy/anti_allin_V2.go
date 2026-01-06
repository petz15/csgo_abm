package strategy

import "math"

func InvestDecisionMaking_anti_allin_v2(ctx StrategyContext_simple) float64 {
	//Round context decision making
	if ctx.IsLastRoundHalf {
		return ctx.Funds
	} else if ctx.IsFirstRoundHalf {
		return ctx.Funds
	} else if ctx.IsSecondRoundHalf && ctx.ConsecutiveLosses > 1 {
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

	//Winning context decision making
	if ctx.EnemySurvivors < 1 && ctx.ConsecutiveWins > 0 {
		//max money the enemy can spend to win the round is Lossbonus * 5 + default equipment *5 + elimination rewards
		current_loss_bonus_enemy := math.Min(float64(ctx.ConsecutiveWins), float64(len(ctx.GameRules_strategy.LossBonus)-1))
		max_enemy_investment := ctx.GameRules_strategy.LossBonus[int(current_loss_bonus_enemy)]*5 + ctx.GameRules_strategy.DefaultEquipment*5
		if ctx.Side {
			//CT side elimination rewards
			max_enemy_investment += (ctx.GameRules_strategy.AdditionalReward_T_Elimination + ctx.GameRules_strategy.EliminationReward) * (float64(5 - ctx.OwnSurvivors))
		} else {
			//T side elimination rewards
			max_enemy_investment += (ctx.GameRules_strategy.AdditionalReward_CT_Elimination + ctx.GameRules_strategy.EliminationReward) * (float64(5 - ctx.OwnSurvivors))
		}
		investment := (max_enemy_investment * 1.5) - ctx.Equipment //aim to have double the investment of the enemy
		return math.Min(investment, ctx.Funds)
	} else if ctx.ConsecutiveWins == 1 {
		return ctx.Funds
	} else if ctx.ConsecutiveWins >= 2 {
		ratio := 1 - (float64(ctx.ConsecutiveWins) * 0.05)
		ratio = math.Max(ratio, 0.4) //set a minimum ratio of 0.4
		return ctx.Funds * ratio
	}

	//Losing context decision making
	if ctx.ConsecutiveLosses == 1 {
		if ctx.OwnSurvivors >= 3 { //if we had at least 3 survivors, we can afford to invest more
			return ctx.Funds * 1.0
		}
		return 0
	} else if ctx.ConsecutiveLosses == 2 {
		if ctx.OwnSurvivors >= 2 { //if we had at least 2 survivors, we can afford to invest more
			return ctx.Funds
		}
		return 0
	} else if ctx.ConsecutiveLosses == 3 {
		return ctx.Funds * 0.8
	} else if ctx.ConsecutiveLosses >= 4 {
		return ctx.Funds * 0.9
	}

	return ctx.Funds * 0.85
}
