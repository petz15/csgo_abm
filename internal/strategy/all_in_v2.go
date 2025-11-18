package strategy

func InvestDecisionMaking_allin_v2(ctx StrategyContext_simple) float64 {
	// Implement your decision-making logic here
	// This is a placeholder implementation

	if ctx.IsLastRoundHalf {
		return ctx.Funds
	} else if ctx.IsPistolRound {
		return ctx.Funds
	} else if ctx.HalfLength-ctx.OpponentScore == 1 && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.HalfLength-ctx.OpponentScore == 0 && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.HalfLength-ctx.OwnScore == 1 && !ctx.IsOvertime {
		return ctx.Funds
	} else if ctx.HalfLength-ctx.OwnScore == 0 && !ctx.IsOvertime {
		return ctx.Funds
	}

	if ctx.ConsecutiveLosses < 1 {
		return ctx.Funds * 0.9
	} else if ctx.ConsecutiveLosses == 1 {
		return ctx.Funds * 0.2
	} else {
		return ctx.Funds * 0.9
	}
}
