package models

func InvestDecisionMaking_half(ctx StrategyContext_simple) float64 {

	if ctx.CurrentRound == 15 || ctx.CurrentRound == 30 {
		return ctx.Funds // Invest half of the funds at round 15 and 30
	} else {
		return ctx.Funds / 2 // Invest half of the funds in other rounds
	}

}
