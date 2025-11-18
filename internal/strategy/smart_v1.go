package strategy

func InvestDecisionMaking_smart_v1(ctx StrategyContext_simple) float64 {

	if ctx.IsLastRoundHalf {
		return ctx.Funds
	} else if ctx.IsPistolRound {
		return ctx.Funds
	} else if ctx.HalfLength-ctx.OpponentScore == 1 && !ctx.IsOvertime {
		return ctx.Funds * 0.8
	} else if ctx.HalfLength-ctx.OpponentScore == 0 && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.HalfLength-ctx.OwnScore == 1 && !ctx.IsOvertime {
		return ctx.Funds * 0.8
	} else if ctx.HalfLength-ctx.OwnScore == 0 && !ctx.IsOvertime {
		return ctx.Funds
	}

	if ctx.ConsecutiveLosses < 1 {
		return ctx.Funds * 0.8
	} else if ctx.ConsecutiveLosses == 1 {
		return ctx.Funds * 0.2
	} else if ctx.ConsecutiveLosses == 2 {
		return ctx.Funds * 0.3
	} else if ctx.ConsecutiveLosses == 3 {
		return ctx.Funds * 0.4
	} else {
		return ctx.Funds * 0.9
	}

}
