package engine

import (
	"CSGO_ABM/internal/models"
)

type StrategyManager struct {
	Name string
}

func CallStrategy(team *Team, curround int, curscoreopo int) float64 {
	// Implement the logic to call the appropriate strategy for the team
	switch team.Strategy {
	case "all_in":
		return models.InvestDecisionMaking_allin(team.Funds)
	case "default_half":
		return models.InvestDecisionMaking_half(team.Funds, curround, curscoreopo)
	case "adaptive_eco_v1":
		return models.InvestDecisionMaking_adaptive_v1(team.Funds, curround, curscoreopo)
	default:
		return models.InvestDecisionMaking_allin(team.Funds)
	}
}
