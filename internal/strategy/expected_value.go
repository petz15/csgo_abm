package strategy

import "math"

func InvestDecisionMaking_expected_value(ctx StrategyContext_simple) float64 {
	if ctx.IsLastRoundHalf {
		return ctx.Funds
	}

	// spends as much as the expected value of the earnings of this round i.e.
	// 0.5 * (funds won if win) + 0.5 * (funds won if lose) approximately
	approx_win_funds := 0.0
	if ctx.Side {
		//CT Side win option are reasoncode 3 & 4
		approx_win_funds += 0.5 * (ctx.GameRules_strategy.RoundOutcomeReward[2]*5 + ctx.GameRules_strategy.BombdefuseReward) //defuse win
		approx_win_funds += 0.5 * (ctx.GameRules_strategy.RoundOutcomeReward[3] * 5)                                         //elimination win
	} else {
		//T Side win option are reasoncode 1 & 2
		approx_win_funds += 0.5 * (ctx.GameRules_strategy.RoundOutcomeReward[0]*5 + ctx.GameRules_strategy.BombplantReward) //bomb explosion win
		approx_win_funds += 0.5 * (ctx.GameRules_strategy.RoundOutcomeReward[1] * 5)                                        //elimination win
	}

	approx_loss_funds := ctx.GameRules_strategy.LossBonus[int(math.Min(float64(ctx.LossBonusLevel+1), float64(len(ctx.GameRules_strategy.LossBonus)-1)))] * 5

	investment := 0.5*approx_win_funds + 0.5*approx_loss_funds
	return math.Min(investment, ctx.Funds)
}
