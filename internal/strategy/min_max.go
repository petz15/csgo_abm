package strategy

func InvestDecisionMaking_min_max(ctx StrategyContext_simple) float64 {
	// min_max, spends all, in the first 3 rounds of a half, then saves the next 3, then spends 3 etc.
	//always spending the last two rounds of a half or when a team is about to win.

	//in overtime it spends all if OT half is <=3 else it spends 2 rounds all, 2 rounds saves etc
	//spending all if a team is about to win.
	if !ctx.IsOvertime {
		if ctx.CurrentRound <= 3 || (ctx.CurrentRound > ctx.GameRules_strategy.HalfLength && ctx.CurrentRound-ctx.GameRules_strategy.HalfLength <= 3) {
			return ctx.Funds
		} else if ctx.OpponentScore == ctx.GameRules_strategy.HalfLength || ctx.OwnScore == ctx.GameRules_strategy.HalfLength {
			return ctx.Funds
		} else if ctx.CurrentRound == ctx.GameRules_strategy.HalfLength-1 || ctx.CurrentRound == ctx.GameRules_strategy.HalfLength-2 {
			return ctx.Funds
		} else if ctx.CurrentRound == ctx.GameRules_strategy.HalfLength*2-1 || ctx.CurrentRound == ctx.GameRules_strategy.HalfLength*2-2 {
			return ctx.Funds
		}

		if (ctx.CurrentRound)%6 <= 3 {
			// spend round
			return ctx.Funds
		} else {
			// save round
			return 0
		}
	} else {
		if ctx.GameRules_strategy.OTHalfLength <= 3 {
			return ctx.Funds
		} else if ctx.OpponentScore == ctx.GameRules_strategy.HalfLength+(ctx.GameRules_strategy.OTHalfLength*ctx.CurrentRound/(ctx.GameRules_strategy.OTHalfLength*2)) || ctx.OwnScore == ctx.GameRules_strategy.HalfLength+(ctx.GameRules_strategy.OTHalfLength*ctx.CurrentRound/(ctx.GameRules_strategy.OTHalfLength*2)) {
			return ctx.Funds
		} else if ctx.CurrentRound%4 <= 2 {
			return ctx.Funds
		} else {
			return 0
		}
	}

	return 0
}
