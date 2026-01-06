package engine

import (
	"csgo_abm/internal/strategy"
	"fmt"
	"math"
)

type StrategyManager struct {
	Name string
}

func CallStrategy(team *Team, opponent *Team, curround int, isOvertime bool, gameR GameRules, g *Game) float64 {
	// Build context once
	ctx := strategy.StrategyContext_simple{
		Funds:                   team.GetCurrentFunds(),
		CurrentRound:            curround,
		OpponentScore:           opponent.GetScore(),
		OwnScore:                team.GetScore(),
		ConsecutiveLosses:       team.GetConsecutiveloss(),
		ConsecutiveWins:         team.GetConsecutivewins(),
		LossBonusLevel:          team.GetCurrentLossBonusLevel(),
		LossBonusLevel_opponent: opponent.GetCurrentLossBonusLevel(),
		Side:                    team.GetSide(),
		Equipment:               team.GetRSEquipment(),
		IsOvertime:              isOvertime,
		OvertimeAmount:          g.OTcounter,
		IsFirstRoundHalf:        IsFirstRoundHalf(curround, gameR.HalfLength, isOvertime, gameR.OTHalfLength),
		IsSecondRoundHalf:       IsSecondRoundHalf(curround, gameR.HalfLength, team.GetScore(), opponent.GetScore(), isOvertime, gameR.OTHalfLength),
		IsLastRoundHalf:         isLastRoundHalf(curround, gameR.HalfLength, isOvertime, gameR.OTHalfLength),
		OwnSurvivors:            team.GetpreviousSurvivors(),
		EnemySurvivors:          opponent.GetpreviousSurvivors(),
		RoundEndReason:          g.GetPreviousRoundEndReason(),
		Is_BombPlanted:          g.GetPreviousBombPlant(),
		RNG:                     g.rng,
		GameRules_strategy: strategy.GameRules_strategymanager{
			DefaultEquipment:                gameR.DefaultEquipment,
			OTFunds:                         gameR.OTFunds,
			OTEquipment:                     gameR.OTEquipment,
			StartingFunds:                   gameR.StartingFunds,
			HalfLength:                      gameR.HalfLength,
			OTHalfLength:                    gameR.OTHalfLength,
			MaxFunds:                        gameR.MaxFunds,
			LossBonusCalc:                   gameR.LossBonusCalc,
			WithSaves:                       gameR.WithSaves,
			LossBonus:                       gameR.LossBonus,
			RoundOutcomeReward:              gameR.RoundOutcomeReward,
			EliminationReward:               gameR.EliminationReward,
			BombplantRewardall:              gameR.BombplantRewardall,
			BombplantReward:                 gameR.BombplantReward,
			BombdefuseReward:                gameR.BombdefuseReward,
			AdditionalReward_CT_Elimination: gameR.AdditionalReward_CT_Elimination,
			AdditionalReward_T_Elimination:  gameR.AdditionalReward_T_Elimination,
		},
	}

	// Get strategy function from registry
	strategyFunc, err := strategy.GetStrategy(team.Strategy)
	if err != nil {
		// This should never happen if validation is done upfront
		// But provide a safe fallback just in case
		panic(fmt.Sprintf("FATAL: Invalid strategy '%s' for team - this should have been caught during validation!", team.Strategy))
	}

	invest := strategyFunc(ctx)
	if err != nil {
		panic(fmt.Sprintf("FATAL: Error executing strategy '%s': %v", team.Strategy, err))
	}
	return invest
}

//most of the following functions are used to enhance the context for decision making
//which should be done by each strategy individually. For now, for testing purposes, they are implemented here.

//TODO: Clean up these functions and move them to a more appropriate place in the strategy package i.e. into the individual strategies

// Helper functions for enhanced context
func IsFirstRoundHalf(round int, halfLength int, isOT bool, OTHalfLength int) bool {
	if !isOT {
		return round == 1 || round == halfLength+1
	} else {
		return (round-halfLength*2)%OTHalfLength == 1
	}

}

func IsSecondRoundHalf(round int, halfLength int, ownScore, opponentScore int, isOT bool, OTHalfLength int) bool {
	if !isOT {
		return round == 2 || round == halfLength+2
	} else {
		return (round-halfLength*2)%OTHalfLength == 2
	}
}

func isLastRoundHalf(round int, halfLength int, isOT bool, OTHalfLength int) bool {
	if !isOT {
		return round == halfLength || round == halfLength*2
	} else {
		return (round-halfLength*2)%OTHalfLength == 0
	}
}

func calculateRoundImportance(ownScore, opponentScore, round int, halfLength int, isOT bool, OTHalfLength int) float64 {
	// Calculate round importance based on score situation and round number
	scoreDiff := math.Abs(float64(ownScore - opponentScore))
	totalScore := ownScore + opponentScore

	importance := 0.5 // Base importance

	// Pistol rounds are always important
	if IsFirstRoundHalf(round, halfLength, isOT, OTHalfLength) {
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
	if isLastRoundHalf(round, halfLength, isOT, OTHalfLength) {
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
