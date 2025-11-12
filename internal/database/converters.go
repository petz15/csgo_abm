package database

import (
	"encoding/json"
	"fmt"
	"time"
)

// Converter functions to transform simulation data into database-ready formats

// ConvertGameToMatchData converts your game/simulation result to MatchData
// You'll need to adapt this to your actual Game struct
func ConvertGameToMatchData(
	matchID string,
	batchID string,
	experimentName string,
	simNumber int,
	team1Name, team2Name string,
	team1Strategy, team2Strategy string,
	team1Score, team2Score int,
	totalRounds int,
	overtimeRounds int,
	durationSeconds float64,
	gameRules interface{},
) (MatchData, error) {

	winner := "draw"
	if team1Score > team2Score {
		winner = team1Name
	} else if team2Score > team1Score {
		winner = team2Name
	}

	wentToOvertime := overtimeRounds > 0

	// Convert game rules to JSON
	gameRulesJSON := ""
	if gameRules != nil {
		jsonBytes, err := json.Marshal(gameRules)
		if err != nil {
			return MatchData{}, fmt.Errorf("failed to marshal game rules: %w", err)
		}
		gameRulesJSON = string(jsonBytes)
	}

	return MatchData{
		ID:                matchID,
		BatchSimulationID: batchID,
		ExperimentName:    experimentName,
		SimulationNumber:  simNumber,
		Team1Name:         team1Name,
		Team2Name:         team2Name,
		Team1Strategy:     team1Strategy,
		Team2Strategy:     team2Strategy,
		Team1FinalScore:   team1Score,
		Team2FinalScore:   team2Score,
		Winner:            winner,
		TotalRounds:       totalRounds,
		OvertimeRounds:    overtimeRounds,
		WentToOvertime:    wentToOvertime,
		DurationSeconds:   durationSeconds,
		Timestamp:         time.Now(),
		GameRulesJSON:     gameRulesJSON,
	}, nil
}

// ConvertRoundToRoundData converts a round from your simulation to RoundData
func ConvertRoundToRoundData(
	matchID string,
	roundNum int,
	halfNum int,
	isOT bool,
	winner string,
	team1Name string,
	team1FundsStart, team1FundsEnd float64,
	team1Spent, team1EquipValue float64,
	team1ScoreBefore, team1ScoreAfter int,
	team1Losses int,
	team1Survivors int,
	team2Name string,
	team2FundsStart, team2FundsEnd float64,
	team2Spent, team2EquipValue float64,
	team2ScoreBefore, team2ScoreAfter int,
	team2Losses int,
	team2Survivors int,
) RoundData {

	economicAdv := team1FundsStart - team2FundsStart
	spendingDiff := team1Spent - team2Spent

	// Detect pistol rounds (round 1 and 16 in regular play)
	isPistolRound := (roundNum == 1) || (roundNum == 16 && !isOT)

	// Detect eco rounds (low spending, e.g., < $2000 per team)
	isEcoRound := (team1Spent < 2000 || team2Spent < 2000)

	return RoundData{
		MatchID:                matchID,
		RoundNumber:            roundNum,
		HalfNumber:             halfNum,
		IsOvertime:             isOT,
		Winner:                 winner,
		WinReason:              "", // Can be extended
		Team1Name:              team1Name,
		Team1FundsStart:        team1FundsStart,
		Team1FundsEnd:          team1FundsEnd,
		Team1EquipmentSpent:    team1Spent,
		Team1EquipmentValue:    team1EquipValue,
		Team1ScoreBefore:       team1ScoreBefore,
		Team1ScoreAfter:        team1ScoreAfter,
		Team1ConsecutiveLosses: team1Losses,
		Team1SurvivingPlayers:  team1Survivors,
		Team2Name:              team2Name,
		Team2FundsStart:        team2FundsStart,
		Team2FundsEnd:          team2FundsEnd,
		Team2EquipmentSpent:    team2Spent,
		Team2EquipmentValue:    team2EquipValue,
		Team2ScoreBefore:       team2ScoreBefore,
		Team2ScoreAfter:        team2ScoreAfter,
		Team2ConsecutiveLosses: team2Losses,
		Team2SurvivingPlayers:  team2Survivors,
		EconomicAdvantage:      economicAdv,
		SpendingDifferential:   spendingDiff,
		IsPistolRound:          isPistolRound,
		IsEcoRound:             isEcoRound,
		Timestamp:              time.Now(),
	}
}

// CalculateBatchSummaryStats calculates aggregate statistics from match results
// This is a helper to compute summary statistics across multiple matches
func CalculateBatchSummaryStats(
	batchID string,
	strategyName string,
	matches []MatchData,
	allRounds []RoundData,
) BatchSummaryStats {

	totalMatches := len(matches)
	totalWins := 0
	totalRoundsPlayed := 0
	roundsWon := 0
	totalOvertimeMatches := 0
	totalDuration := 0.0

	// Calculate win statistics
	for _, match := range matches {
		totalRoundsPlayed += match.TotalRounds
		if match.Winner == strategyName {
			totalWins++
		}
		if match.WentToOvertime {
			totalOvertimeMatches++
		}
		totalDuration += match.DurationSeconds
	}

	// Calculate round-level statistics
	var totalFunds, totalSpending, totalEquipValue float64
	roundCount := 0

	for _, round := range allRounds {
		// Count rounds where this strategy won
		if round.Winner == strategyName {
			roundsWon++
		}

		// Aggregate economic metrics based on which team is the strategy
		if round.Team1Name == strategyName {
			totalFunds += round.Team1FundsStart
			totalSpending += round.Team1EquipmentSpent
			totalEquipValue += round.Team1EquipmentValue
			roundCount++
		} else if round.Team2Name == strategyName {
			totalFunds += round.Team2FundsStart
			totalSpending += round.Team2EquipmentSpent
			totalEquipValue += round.Team2EquipmentValue
			roundCount++
		}
	}

	// Calculate averages
	winRate := 0.0
	if totalMatches > 0 {
		winRate = float64(totalWins) / float64(totalMatches)
	}

	roundWinRate := 0.0
	if totalRoundsPlayed > 0 {
		roundWinRate = float64(roundsWon) / float64(totalRoundsPlayed)
	}

	avgFunds := 0.0
	avgSpending := 0.0
	avgEquipValue := 0.0
	if roundCount > 0 {
		avgFunds = totalFunds / float64(roundCount)
		avgSpending = totalSpending / float64(roundCount)
		avgEquipValue = totalEquipValue / float64(roundCount)
	}

	avgDuration := 0.0
	if totalMatches > 0 {
		avgDuration = totalDuration / float64(totalMatches)
	}

	overtimeFreq := 0.0
	if totalMatches > 0 {
		overtimeFreq = float64(totalOvertimeMatches) / float64(totalMatches)
	}

	// Calculate average match score (rounds won per match)
	avgMatchScore := 0.0
	if totalMatches > 0 {
		avgMatchScore = float64(roundsWon) / float64(totalMatches)
	}

	return BatchSummaryStats{
		BatchSimulationID:   batchID,
		StrategyName:        strategyName,
		TotalMatches:        totalMatches,
		TotalWins:           totalWins,
		TotalLosses:         totalMatches - totalWins,
		WinRate:             winRate,
		TotalRoundsPlayed:   totalRoundsPlayed,
		RoundsWon:           roundsWon,
		RoundsLost:          totalRoundsPlayed - roundsWon,
		RoundWinRate:        roundWinRate,
		AvgFundsPerRound:    avgFunds,
		AvgSpendingPerRound: avgSpending,
		AvgEquipmentValue:   avgEquipValue,
		AvgMatchScore:       avgMatchScore,
		OvertimeFrequency:   overtimeFreq,
		AvgMatchDuration:    avgDuration,
		WinRateStdDev:       0.0, // Can be calculated with more sophisticated stats
		ConsistencyScore:    0.0, // Can be calculated based on variance
		Timestamp:           time.Now(),
	}
}
