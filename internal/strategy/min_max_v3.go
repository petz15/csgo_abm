package strategy

func InvestDecisionMaking_min_max_v3(ctx StrategyContext_simple) float64 {

	Score_to_Win := ctx.GameRules_strategy.HalfLength + (ctx.GameRules_strategy.OTHalfLength * (ctx.OvertimeAmount)) + 1

	// min_max_v3, invests all in the beginning and end of halves/overtime. In between it tries to build up wealth
	//when losing until max loss bonus has been reached or funds are larger than average winner funds

	if Score_to_Win-ctx.OpponentScore == 1 || Score_to_Win-ctx.OwnScore == 1 {
		return ctx.Funds
	}

	if ctx.IsFirstRoundHalf {
		return ctx.Funds
	} else if ctx.IsLastRoundHalf {
		return ctx.Funds
	}

	if ctx.ConsecutiveLosses < len(ctx.GameRules_strategy.LossBonus) {
		return 0.0
	}

	return ctx.Funds

}
