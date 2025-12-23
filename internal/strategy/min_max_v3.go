package strategy

func InvestDecisionMaking_min_max_v3(ctx StrategyContext_simple) float64 {
	// min_max_v2, spends all, in the first and last round. otherwise it's a 50-50 chance to spend all or nothing

	if !ctx.IsOvertime {
		if ctx.CurrentRound == 1 || ctx.CurrentRound == ctx.GameRules_strategy.HalfLength+1 {
			return ctx.Funds
		} else if ctx.CurrentRound == ctx.GameRules_strategy.HalfLength || ctx.CurrentRound == ctx.GameRules_strategy.HalfLength*2 {
			return ctx.Funds
		}
	} else {
		curR := ctx.CurrentRound - ctx.GameRules_strategy.HalfLength*2
		if curR%ctx.GameRules_strategy.OTHalfLength == 1 || curR%ctx.GameRules_strategy.OTHalfLength == 0 {
			return ctx.Funds
		}
	}

	if ctx.RNG.Float64() < 0.5 {
		return ctx.Funds
	} else {
		return 0
	}
}
