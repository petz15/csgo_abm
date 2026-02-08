package strategy

func InvestDecisionMaking_half(ctx StrategyContext_simple) float64 {

	return ctx.Funds / 2 // Invest half of the funds in all rounds

}
