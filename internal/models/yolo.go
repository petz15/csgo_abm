package models

import "math/rand"

// InvestDecisionMaking_yolo implements a completely random investment strategy
// This strategy throws caution to the wind and makes random investment decisions
// regardless of game state, economy, or logic - pure chaos!
func InvestDecisionMaking_yolo(funds float64, curround int, curscoreopo int) float64 {
	// Random investment between 10% and 95% of available funds
	// Because why plan when you can YOLO?
	minInvestment := funds * 0.1
	maxInvestment := funds * 0.95

	// Generate random investment amount
	randomFactor := rand.Float64() // 0.0 to 1.0
	investment := minInvestment + (randomFactor * (maxInvestment - minInvestment))

	// Sometimes go completely crazy and invest everything (5% chance)
	if rand.Float64() < 0.05 {
		investment = funds
	}

	// Sometimes be super conservative (5% chance)
	if rand.Float64() < 0.05 {
		investment = funds * 0.05
	}

	return investment
}
