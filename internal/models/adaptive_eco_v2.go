package models

import (
	"math"
)

// InvestDecisionMaking_adaptive_v2 implements an advanced economic strategy with enhanced context awareness
func InvestDecisionMaking_adaptive_v2(ctx StrategyContext_simple) float64 {
	// Economic thresholds (5-player team basis)

	// Determine base economic state
	economicState := assessEconomicStateV2(ctx.Funds)

	// Calculate base investment based on economic state
	baseInvestment := calculateBaseInvestmentV2(ctx.Funds, economicState)

	// Apply round-specific modifiers
	investment := applyRoundModifiersV2(baseInvestment, ctx, economicState)

	// Apply score-based pressure adjustments
	investment = applyScorePressureV2(investment, ctx)

	// Apply side-specific adjustments (CT vs T)
	investment = applySideModifiersV2(investment, ctx)

	// Apply loss bonus considerations
	investment = applyLossBonusLogicV2(investment, ctx)

	// Final safety checks and caps
	investment = applySafetyChecksV2(investment, ctx.Funds)

	return investment
}

// assessEconomicStateV2 determines the team's economic situation
func assessEconomicStateV2(funds float64) string {
	const (
		fullBuyThreshold  = 3500.0 * 5
		forceBuyThreshold = 2500.0 * 5
		ecoThreshold      = 1000.0 * 5
	)

	if funds >= fullBuyThreshold {
		return "healthy"
	} else if funds >= forceBuyThreshold {
		return "moderate"
	} else if funds >= ecoThreshold {
		return "poor"
	}
	return "critical"
}

// calculateBaseInvestmentV2 calculates investment based on economic state
func calculateBaseInvestmentV2(funds float64, economicState string) float64 {
	switch economicState {
	case "healthy":
		return math.Min(funds*0.75, 4000.0*5) // Invest heavily but keep reserves
	case "moderate":
		return math.Min(funds*0.65, 3000.0*5) // Balanced investment
	case "poor":
		return math.Min(funds*0.45, 1800.0*5) // Conservative investment
	case "critical":
		return math.Min(funds*0.25, 800.0*5) // Minimal investment, save for next round
	default:
		return funds * 0.5
	}
}

// applyRoundModifiersV2 adjusts investment based on round type and importance
func applyRoundModifiersV2(baseInvestment float64, ctx StrategyContext_simple, economicState string) float64 {
	investment := baseInvestment

	// Pistol round logic
	if ctx.IsPistolRound {
		// Pistol rounds are crucial - always invest appropriately
		pistolInvestment := math.Min(ctx.Funds*0.8, 800.0*5)
		return math.Max(investment, pistolInvestment)
	}

	// Anti-eco after pistol win (round 2/17)
	if ctx.IsEcoAfterPistol && ctx.OwnScore > ctx.OpponentScore {
		// Light buy against likely eco
		return math.Min(investment*0.7, 2000.0*5)
	}

	// Last round of half - be more aggressive
	if ctx.IsLastRoundHalf {
		multiplier := 1.0 + (ctx.RoundImportance * 0.3) // Up to 30% increase
		investment *= multiplier
	}

	// Overtime rounds - higher stakes
	if ctx.IsOvertime {
		// More aggressive in overtime due to reset after each round pair
		investment = math.Min(investment*1.2, ctx.Funds*0.9)
	}

	return investment
}

// applyScorePressureV2 adjusts investment based on match score situation
func applyScorePressureV2(investment float64, ctx StrategyContext_simple) float64 {
	scoreDiff := ctx.OwnScore - ctx.OpponentScore
	totalScore := ctx.OwnScore + ctx.OpponentScore

	// Calculate match pressure (higher when close to match point)
	matchPressure := float64(totalScore) / 30.0 // Normalize to 0-1 scale

	// When behind, be more aggressive
	if scoreDiff < -3 {
		// Significantly behind - take more risks
		aggressionBonus := 1.0 + (math.Abs(float64(scoreDiff))/15.0)*0.4
		investment *= aggressionBonus
	} else if scoreDiff < -1 {
		// Slightly behind - moderate aggression
		investment *= (1.0 + matchPressure*0.2)
	} else if scoreDiff > 3 {
		// Significantly ahead - be more conservative
		investment *= (0.9 - matchPressure*0.1)
	}

	// Match point situations (when either team is close to winning)
	if ctx.OwnScore >= 14 || ctx.OpponentScore >= 14 {
		// Critical rounds - invest more heavily
		investment *= (1.0 + ctx.RoundImportance*0.3)
	}

	return investment
}

// applySideModifiersV2 applies CT vs T side economic adjustments
func applySideModifiersV2(investment float64, ctx StrategyContext_simple) float64 {
	if ctx.Side { // CT side
		// CTs generally need more expensive equipment (armor + utility)
		// But also have more defensive options for eco rounds
		if ctx.Funds < 2500.0*5 {
			// CT eco rounds can be more effective, invest less
			investment *= 0.85
		} else {
			// CT buy rounds need more investment for utility
			investment *= 1.1
		}
	} else { // T side
		// Ts can be more effective with lighter buys
		// But need coordination for site takes
		if ctx.Funds >= 3000.0*5 {
			// T side full buys don't need as much per player
			investment *= 0.95
		} else if ctx.Funds >= 1500.0*5 {
			// T force buys can be very effective
			investment *= 1.15
		}
	}

	return investment
}

// applyLossBonusLogicV2 factors in loss bonus economics
func applyLossBonusLogicV2(investment float64, ctx StrategyContext_simple) float64 {
	// Loss bonus increases investment capability next round
	if ctx.ConsecutiveLosses >= 2 {
		// Multiple losses mean better economy next round if we lose
		// Can afford to be more aggressive this round
		lossBonusMultiplier := 1.0 + float64(ctx.ConsecutiveLosses)*0.1
		investment *= math.Min(lossBonusMultiplier, 1.4) // Cap at 40% increase
	} else if ctx.ConsecutiveLosses == 0 {
		// Just won - be more conservative to maintain economy
		investment *= 0.9
	}

	return investment
}

// applySafetyChecksV2 ensures investment doesn't exceed safe limits
func applySafetyChecksV2(investment float64, funds float64) float64 {
	// Never invest more than available funds
	investment = math.Min(investment, funds)

	// Always keep minimum reserve for next round (except in critical situations)
	minReserve := 500.0 * 5 // Minimum pistol + armor per player
	if funds > minReserve*2 {
		investment = math.Min(investment, funds-minReserve)
	}

	// Cap maximum investment to prevent over-spending
	maxInvestment := 4500.0 * 5
	investment = math.Min(investment, maxInvestment)

	// Ensure minimum investment (don't invest nothing unless truly broke)
	minInvestment := 200.0 * 5
	if funds >= minInvestment {
		investment = math.Max(investment, minInvestment)
	}

	return investment
}
