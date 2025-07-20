package models

// InvestDecisionMaking_scrooge implements the most conservative strategy possible
// This strategy hoards money like Scrooge McDuck and refuses to spend on anything
// Always invests the absolute minimum to survive
func InvestDecisionMaking_scrooge(funds float64, curround int, curscoreopo int) float64 {
	// Scrooge only spends the bare minimum (team-based: 5 players)
	const minimumSpending = 300.0 * 5 // Absolute minimum for 5 players

	// If we have less than minimum, spend what we have
	if funds < minimumSpending {
		return funds * 0.8 // Spend 80% if we're really poor
	}

	// Otherwise, spend the absolute minimum
	// Scrooge believes in saving every penny!
	return minimumSpending
}
