package strategy

func InvestDecisionMaking_casual(ctx StrategyContext_simple) float64 {
	//based off a casual player's (called Tim) logic

	min_treshold := 5000 * 5                    //minimal threshold when all funds are spent
	safety_threshold := 2000 * 5                //safety threshold, which are not to be spent in order to reach the minimum threshold
	max_threshold := 10000*5 - safety_threshold //maximum threshold, above which we only invest up to this amount in order to save up money

	if ctx.IsLastRoundHalf {
		return ctx.Funds
	} else if ctx.IsFirstRoundHalf {
		return ctx.Funds
	} else if ctx.GameRules_strategy.HalfLength-ctx.OpponentScore == 1 && !ctx.IsOvertime {
		return ctx.Funds * 0.8
	} else if ctx.GameRules_strategy.HalfLength-ctx.OpponentScore == 0 && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.GameRules_strategy.HalfLength-ctx.OwnScore == 1 && !ctx.IsOvertime {
		return ctx.Funds * 0.8
	} else if ctx.GameRules_strategy.HalfLength-ctx.OwnScore == 0 && !ctx.IsOvertime {
		return ctx.Funds
	}

	if ctx.Funds > float64(max_threshold) {
		return float64(max_threshold)
	} else if ctx.Funds > float64(min_treshold) {
		return ctx.Funds
	} else if ctx.Funds > float64(safety_threshold) {
		return ctx.Funds - float64(safety_threshold)
	} else {
		return 0
	}

}
