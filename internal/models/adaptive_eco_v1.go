package models

import "math"

// InvestDecisionMaking_adaptive implements an advanced economic strategy that adapts based on:
// - Current economic state and momentum
// - Round importance (pistol, anti-eco, buy rounds)
// - Score differential and match situation
// - Risk assessment based on remaining funds
func InvestDecisionMaking_adaptive_v1(ctx StrategyContext_simple) float64 {
	// Basic economic thresholds (team-based: 5 players)

	funds := ctx.Funds
	curround := ctx.CurrentRound
	curscoreopo := ctx.OpponentScore

	const (
		fullBuyThreshold  = 3500.0 * 5 // Minimum for full buy
		forceBuyThreshold = 2500.0 * 5 // Force buy threshold
		savingThreshold   = 1000.0 * 5 // Below this, save for next round
		maxInvestment     = 4000.0 * 5 // Cap investment to prevent over-spending
	)

	// Calculate round type and importance
	roundType := determineRoundType(curround)
	scoreImportance := calculateScoreImportance(curround, curscoreopo)
	economicState := assessEconomicState(funds)

	// Base investment calculation
	baseInvestment := calculateBaseInvestment(funds, economicState)

	// Apply modifiers based on situation
	investment := baseInvestment

	// Round type modifiers
	switch roundType {
	case "pistol":
		// Pistol rounds are crucial - invest more aggressively
		investment = math.Min(funds, fullBuyThreshold*0.8)
	case "anti-eco":
		// Against likely eco, invest moderately to maintain economy
		investment = math.Min(funds*0.6, forceBuyThreshold)
	case "force-buy":
		// Critical buy round - be aggressive
		investment = math.Min(funds*0.85, maxInvestment)
	case "eco":
		// Save for next round, minimal investment
		investment = math.Min(funds*0.2, 800.0*5)
	case "buy":
		// Standard buy round
		investment = baseInvestment
	}

	// Score-based adjustments
	investment = applyScoreModifiers(investment, funds, curround, curscoreopo, scoreImportance)

	// Economic momentum adjustments
	investment = applyEconomicMomentum(investment, funds, economicState)

	// Final bounds checking
	investment = math.Max(investment, 300.0*5)       // Minimum investment (team-based)
	investment = math.Min(investment, funds)         // Can't spend more than we have
	investment = math.Min(investment, maxInvestment) // Investment cap

	return investment
}

// determineRoundType analyzes the current round to determine the likely round type
func determineRoundType(curround int) string {
	// Pistol rounds (1st and 16th round of each half)
	if curround == 1 || curround == 16 {
		return "pistol"
	}

	// Likely anti-eco rounds (after pistol wins)
	if curround == 2 || curround == 17 {
		return "anti-eco"
	}

	// Late round situations often force critical buys
	if curround >= 14 && curround <= 15 || curround >= 29 && curround <= 30 {
		return "force-buy"
	}

	// Check for likely eco situations (every 3rd round pattern)
	if (curround-1)%3 == 0 && curround > 4 {
		return "eco"
	}

	return "buy"
}

// calculateScoreImportance determines how critical the current round is based on score
func calculateScoreImportance(curround int, curscoreopo int) float64 {
	// Calculate own score (assuming total rounds - opponent score - current round)
	ownScore := (curround - 1) - curscoreopo
	scoreDiff := ownScore - curscoreopo

	// High importance in close games
	if math.Abs(float64(scoreDiff)) <= 2 {
		return 1.2 // 20% importance boost
	}

	// Very high importance when behind significantly
	if scoreDiff <= -3 {
		return 1.4 // 40% importance boost
	}

	// Moderate importance when ahead
	if scoreDiff >= 3 {
		return 0.9 // 10% importance reduction
	}

	return 1.0 // Normal importance
}

// assessEconomicState categorizes the team's economic situation (team-based thresholds)
func assessEconomicState(funds float64) string {
	if funds >= 10000*5 {
		return "rich"
	} else if funds >= 6000*5 {
		return "comfortable"
	} else if funds >= 3500*5 {
		return "moderate"
	} else if funds >= 2000*5 {
		return "tight"
	} else {
		return "poor"
	}
}

// calculateBaseInvestment determines base investment based on economic state (team-based values)
func calculateBaseInvestment(funds float64, economicState string) float64 {
	switch economicState {
	case "rich":
		// When rich, invest aggressively but maintain buffer
		return math.Min(funds*0.7, 4000.0*5)
	case "comfortable":
		// Balanced approach
		return math.Min(funds*0.65, 3500.0*5)
	case "moderate":
		// Careful investment
		return math.Min(funds*0.6, 3000.0*5)
	case "tight":
		// Conservative, prepare for potential eco
		return math.Min(funds*0.5, 2000.0*5)
	case "poor":
		// Save mode, minimal investment
		return math.Min(funds*0.3, 1000.0*5)
	default:
		return funds * 0.5
	}
}

// applyScoreModifiers adjusts investment based on match situation
func applyScoreModifiers(investment float64, funds float64, curround int, curscoreopo int, scoreImportance float64) float64 {
	ownScore := (curround - 1) - curscoreopo
	scoreDiff := ownScore - curscoreopo

	// If losing badly and have funds, be more aggressive
	if scoreDiff <= -4 && funds >= 4000*5 {
		investment *= 1.3
	}

	// If winning comfortably, be more conservative
	if scoreDiff >= 5 {
		investment *= 0.8
	}

	// Apply score importance multiplier
	investment *= scoreImportance

	// Match point situations (assuming 16 rounds to win)
	if ownScore >= 15 || curscoreopo >= 15 {
		// Critical rounds, invest more aggressively
		investment *= 1.25
	}

	return investment
}

// applyEconomicMomentum simulates economic momentum based on recent performance
func applyEconomicMomentum(investment float64, funds float64, economicState string) float64 {
	// Simulate momentum based on economic state
	// In a real implementation, this would track recent round outcomes

	if economicState == "rich" {
		// Assume good momentum, maintain aggressive spending
		return investment * 1.1
	}

	if economicState == "poor" {
		// Assume bad momentum, be more conservative
		return investment * 0.9
	}

	// Moderate momentum adjustment for middle states
	return investment
}
