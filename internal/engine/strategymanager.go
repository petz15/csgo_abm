package engine

import (
	"csgo-economy-sim/internal/models"
)

type StrategyManager struct {
	Name string
}

func CallStrategy(team *Team, curround int, curscoreopo int) int {
	// Implement the logic to call the appropriate strategy for the team
	switch team.Strategy {
	case "all_in":
		return models.InvestDecisionMaking_allin(team.Funds)
	case "simple":
		return models.InvestDecisionMaking_half(team.Funds, curround, curscoreopo)
	default:
		return models.InvestDecisionMaking_allin(team.Funds)
	}
}
