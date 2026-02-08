package strategy

import (
	"math"
	"math/rand"
)

func InvestDecisionMaking_random(ctx StrategyContext_simple) float64 {
	// Randomly invest between 0% and 100% of available funds
	// Because why plan when you can YOLO?

	// Generate random ratio between 0.0 and 1.0
	randomRatio := rand.Float64()
	investment := math.Round(ctx.Funds * randomRatio) // in order to avoid too many decimals

	return investment
}
