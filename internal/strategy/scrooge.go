package strategy

func InvestDecisionMaking_scrooge(ctx StrategyContext_simple) float64 {

	// Scrooge believes in saving every penny!
	return ctx.Funds * 0.1
}
