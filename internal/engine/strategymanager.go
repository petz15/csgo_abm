package engine

import (
	"csgo_abm/internal/models"
	"math"
)

type StrategyManager struct {
	Name string
}

func CallStrategy(team *Team, opponent *Team, curround int, isOvertime bool, gameR GameRules) float64 {
	// Implement the logic to call the appropriate strategy for the team
	ctx := models.StrategyContext_simple{
		Funds:             team.GetCurrentFunds(),
		CurrentRound:      curround,
		OpponentScore:     opponent.GetScore(),
		OwnScore:          team.GetScore(),
		ConsecutiveLosses: team.GetConsecutiveloss(),
		Side:              team.GetSide(),
		Equipment:         team.GetRSEquipment(),
		IsOvertime:        isOvertime,
		IsPistolRound:     isPistolRound(curround, gameR.HalfLength),
		IsEcoAfterPistol:  isEcoAfterPistol(curround, gameR.HalfLength, team.GetScore(), opponent.GetScore()),
		IsLastRoundHalf:   isLastRoundHalf(curround, gameR.HalfLength),
		RoundImportance:   calculateRoundImportance(team.GetScore(), opponent.GetScore(), curround, gameR.HalfLength),
		EconomicAdvantage: calculateEconomicAdvantage(team.GetCurrentFunds(), opponent.GetCurrentFunds()),
		WinProbability:    calculateWinProbability(team.GetCurrentFunds(), opponent.GetCurrentFunds(), team.GetScore(), opponent.GetScore()),
		HalfLength:        gameR.HalfLength,
		OTHalfLength:      gameR.OTHalfLength,
	}

	switch team.Strategy {
	case "all_in":
		return models.InvestDecisionMaking_allin(ctx)
	case "default_half":
		return models.InvestDecisionMaking_half(ctx)
	case "adaptive_eco_v1":
		return models.InvestDecisionMaking_adaptive_v1(ctx)
	case "adaptive_eco_v2":
		return models.InvestDecisionMaking_adaptive_v2(ctx)
	case "yolo":
		return models.InvestDecisionMaking_yolo(ctx)
	case "scrooge":
		return models.InvestDecisionMaking_scrooge(ctx)
	default:
		return models.InvestDecisionMaking_allin(ctx)
	}
}

//most of the following functions are used to enhance the context for decision making
//which should be done by each model individually. For now, for testing purposes, they are implemented here.

//TODO: Clean up these functions and move them to a more appropriate place in the models package i.e. into the individual models

// Helper functions for enhanced context
func isPistolRound(round int, halfLength int) bool {
	return round == 1 || round == halfLength+1
}

func isEcoAfterPistol(round int, halfLength int, ownScore, opponentScore int) bool {
	// Round 2 or 17, and we won the pistol round
	return (round == 2 && ownScore > opponentScore) || (round == halfLength+1 && ownScore > opponentScore)
}

func isLastRoundHalf(round int, halfLength int) bool {
	return round == halfLength || round == halfLength*2
}

func calculateRoundImportance(ownScore, opponentScore, round int, halfLength int) float64 {
	// Calculate round importance based on score situation and round number
	scoreDiff := math.Abs(float64(ownScore - opponentScore))
	totalScore := ownScore + opponentScore

	importance := 0.5 // Base importance

	// Pistol rounds are always important
	if isPistolRound(round, halfLength) {
		importance = 0.8
	}

	// Match point situations
	if ownScore >= halfLength || opponentScore >= halfLength {
		importance = 0.9
	}

	// Close games are more important
	if scoreDiff <= 2 && totalScore >= 20 {
		importance += 0.2
	}

	// Last round of half
	if isLastRoundHalf(round, halfLength) {
		importance += 0.1
	}

	// Overtime rounds
	if round > halfLength*2 {
		importance = 0.9
	}

	// Ensure importance is between 0 and 1
	if importance > 1.0 {
		importance = 1.0
	}

	return importance
}

func calculateEconomicAdvantage(ownFunds, opponentFunds float64) float64 {
	// Calculate relative economic advantage (-1 to 1 scale)
	totalFunds := ownFunds + opponentFunds
	if totalFunds == 0 {
		return 0
	}

	// Normalize to -1 to 1 scale
	advantage := (ownFunds - opponentFunds) / totalFunds
	return advantage
}

func calculateWinProbability(ownFunds, opponentFunds float64, ownScore, opponentScore int) float64 {
	// Simple win probability calculation based on economic and score factors
	// This is a placeholder - could be enhanced with CSF calculations

	// Economic factor (0.4 to 0.6 based on funds)
	economicFactor := 0.5
	if ownFunds+opponentFunds > 0 {
		economicFactor = ownFunds / (ownFunds + opponentFunds)
	}

	// Score factor (weighted by how close the game is)
	totalRounds := ownScore + opponentScore
	if totalRounds > 0 {
		scoreFactor := float64(ownScore) / float64(totalRounds)
		// Weight economic vs score factors (economic matters more early, score matters more late)
		gameProgress := float64(totalRounds) / 30.0
		if gameProgress > 1.0 {
			gameProgress = 1.0
		}

		// Blend economic and score factors
		return economicFactor*(1.0-gameProgress) + scoreFactor*gameProgress
	}

	return economicFactor
}
