package strategy

func InvestDecisionMaking_allin_v2(ctx StrategyContext_simple) float64 {

	if ctx.IsLastRoundHalf {
		return ctx.Funds
	} else if ctx.IsPistolRound {
		return ctx.Funds
	} else if ctx.GameRules_strategy.HalfLength-ctx.OpponentScore == 1 && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.GameRules_strategy.HalfLength-ctx.OpponentScore == 0 && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.GameRules_strategy.HalfLength-ctx.OwnScore == 1 && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.GameRules_strategy.HalfLength-ctx.OwnScore == 0 && !ctx.IsOvertime {
		return ctx.Funds
	}

	if ctx.ConsecutiveLosses < 1 {
		return ctx.Funds * 0.9
	} else if ctx.ConsecutiveLosses == 1 {
		return ctx.Funds
	} else if ctx.ConsecutiveLosses == 2 {
		return ctx.Funds * 0.1
	} else if ctx.ConsecutiveLosses == 3 {
		return ctx.Funds * 0.3
	} else if ctx.ConsecutiveLosses >= 4 {
		return ctx.Funds
	}

	return ctx.Funds * 0.8
}
