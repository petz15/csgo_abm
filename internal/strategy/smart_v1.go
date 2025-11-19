package strategy

import "math"

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

	if ctx.WithSaves {
		if ctx.ConsecutiveLosses < 1 && ctx.EnemySurvivors < 1 {
			return math.Min(ctx.DefaultEquipment*2*5, ctx.Funds)
		}
	}

	if ctx.ConsecutiveLosses < 1 {
		return ctx.Funds * 0.9
	} else if ctx.ConsecutiveLosses == 1 {
		return ctx.Funds * 0.1
	} else if ctx.ConsecutiveLosses >= 4 {
		return ctx.Funds
	}

	return ctx.Funds * 0.8

}
