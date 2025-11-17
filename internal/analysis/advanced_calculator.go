package analysis

import (
	"csgo_abm/internal/engine"
	"math"
	"sort"
)

// AdvancedAnalyzer processes game data for advanced analysis
type AdvancedAnalyzer struct {
	analysis *AdvancedAnalysis

	// Temporary storage for calculations
	economicAdvantages  []float64
	equipmentAdvantages []float64
	allFundsTeam1       []float64
	allFundsTeam2       []float64
	allEquipmentTeam1   []float64
	allEquipmentTeam2   []float64

	// Round-by-round data for time series
	roundDataMap map[int]*RoundTimePoint

	// Streak tracking
	currentTeam1Streak int
	currentTeam2Streak int
	lastWinner         string

	// Comeback tracking
	comebackScenarios map[int][]ComebackAttempt

	// Half/side tracking
	halfSideRounds []HalfSideRound

	// Win condition tracking
	winConditionSamples []WinConditionSample
}

// ComebackAttempt tracks a potential comeback scenario
type ComebackAttempt struct {
	GameID                string
	DeficitSize           int
	StartingRound         int
	Team                  string
	Success               bool
	EconomicAdvantageAvg  float64
	EquipmentAdvantageAvg float64
}

// HalfSideRound tracks performance by half and side
type HalfSideRound struct {
	GameID         string
	RoundNumber    int
	IsFirstHalf    bool
	Team1IsCT      bool
	Team1Won       bool
	Team1Funds     float64
	Team1Equipment float64
	Team2Funds     float64
	Team2Equipment float64
}

// WinConditionSample represents a single round for win condition analysis
type WinConditionSample struct {
	Team1Won           bool
	FundsAdvantage     float64 // Team1 - Team2
	EquipmentAdvantage float64
	Team1ConsecLoss    int
	Team2ConsecLoss    int
	Team1Survivors     int
	Team2Survivors     int
	Team1Spent         float64
	Team2Spent         float64
}

// NewAdvancedAnalyzer creates a new analyzer
func NewAdvancedAnalyzer() *AdvancedAnalyzer {
	return &AdvancedAnalyzer{
		analysis:            NewAdvancedAnalysis(),
		economicAdvantages:  make([]float64, 0, 10000),
		equipmentAdvantages: make([]float64, 0, 10000),
		allFundsTeam1:       make([]float64, 0, 10000),
		allFundsTeam2:       make([]float64, 0, 10000),
		allEquipmentTeam1:   make([]float64, 0, 10000),
		allEquipmentTeam2:   make([]float64, 0, 10000),
		roundDataMap:        make(map[int]*RoundTimePoint),
		comebackScenarios:   make(map[int][]ComebackAttempt),
		halfSideRounds:      make([]HalfSideRound, 0, 10000),
		winConditionSamples: make([]WinConditionSample, 0, 10000),
		lastWinner:          "",
	}
}

// ProcessGame analyzes a single game
func (aa *AdvancedAnalyzer) ProcessGame(game *engine.Game) {
	if len(game.Rounds) == 0 {
		return
	}

	// Track score for comeback analysis
	team1Score := 0
	team2Score := 0

	// Track streaks for this game
	currentStreak := 0
	currentStreakTeam := ""
	streakStartRound := 1
	streakStartEconAdvantage := 0.0
	streakStartEquipAdvantage := 0.0

	// Process each round
	for i, round := range game.Rounds {
		roundNum := i + 1

		// Get team data
		team1Data := game.Team1.RoundData[i]
		team2Data := game.Team2.RoundData[i]

		// Calculate advantages
		economicAdvantage := team1Data.Funds - team2Data.Funds
		equipmentAdvantage := team1Data.FTE_Eq_value - team2Data.FTE_Eq_value

		// Store for distributions
		aa.economicAdvantages = append(aa.economicAdvantages, economicAdvantage)
		aa.equipmentAdvantages = append(aa.equipmentAdvantages, equipmentAdvantage)
		aa.allFundsTeam1 = append(aa.allFundsTeam1, team1Data.Funds)
		aa.allFundsTeam2 = append(aa.allFundsTeam2, team2Data.Funds)
		aa.allEquipmentTeam1 = append(aa.allEquipmentTeam1, team1Data.FTE_Eq_value)
		aa.allEquipmentTeam2 = append(aa.allEquipmentTeam2, team2Data.FTE_Eq_value)

		// Update time series
		aa.updateTimeSeries(roundNum, &team1Data, &team2Data, round.IsT1WinnerTeam)

		// Track half/side data
		isFirstHalf := roundNum <= game.GameRules.HalfLength
		aa.halfSideRounds = append(aa.halfSideRounds, HalfSideRound{
			GameID:         game.ID,
			RoundNumber:    roundNum,
			IsFirstHalf:    isFirstHalf,
			Team1IsCT:      round.IsT1CT,
			Team1Won:       round.IsT1WinnerTeam,
			Team1Funds:     team1Data.Funds,
			Team1Equipment: team1Data.FTE_Eq_value,
			Team2Funds:     team2Data.Funds,
			Team2Equipment: team2Data.FTE_Eq_value,
		})

		// Track win conditions
		aa.winConditionSamples = append(aa.winConditionSamples, WinConditionSample{
			Team1Won:           round.IsT1WinnerTeam,
			FundsAdvantage:     economicAdvantage,
			EquipmentAdvantage: equipmentAdvantage,
			Team1ConsecLoss:    team1Data.Consecutiveloss,
			Team2ConsecLoss:    team2Data.Consecutiveloss,
			Team1Survivors:     team1Data.Survivors,
			Team2Survivors:     team2Data.Survivors,
			Team1Spent:         team1Data.Spent,
			Team2Spent:         team2Data.Spent,
		})

		// Track spending decisions
		aa.analysis.Distributions.SpendingDecisions = append(aa.analysis.Distributions.SpendingDecisions,
			SpendingDecision{
				AvailableFunds: team1Data.Funds,
				AmountSpent:    team1Data.Spent,
				SpendRatio:     team1Data.Spent / math.Max(team1Data.Funds, 1),
				RoundWon:       round.IsT1WinnerTeam,
				Team:           "Team1",
			},
			SpendingDecision{
				AvailableFunds: team2Data.Funds,
				AmountSpent:    team2Data.Spent,
				SpendRatio:     team2Data.Spent / math.Max(team2Data.Funds, 1),
				RoundWon:       !round.IsT1WinnerTeam,
				Team:           "Team2",
			},
		)

		// Update scores
		if round.IsT1WinnerTeam {
			team1Score++
		} else {
			team2Score++
		}

		// Track streaks
		winner := "Team1"
		if !round.IsT1WinnerTeam {
			winner = "Team2"
		}

		if winner == currentStreakTeam {
			currentStreak++
		} else {
			// End previous streak
			if currentStreak > 0 {
				aa.recordStreak(game.ID, currentStreakTeam, currentStreak, streakStartRound, roundNum-1,
					streakStartEconAdvantage, economicAdvantage, streakStartEquipAdvantage, equipmentAdvantage)
			}
			// Start new streak
			currentStreak = 1
			currentStreakTeam = winner
			streakStartRound = roundNum
			streakStartEconAdvantage = economicAdvantage
			streakStartEquipAdvantage = equipmentAdvantage
		}
	}

	// Record final streak
	if currentStreak > 0 {
		finalEconAdvantage := 0.0
		finalEquipAdvantage := 0.0
		if len(game.Rounds) > 0 {
			lastIdx := len(game.Rounds) - 1
			team1Data := game.Team1.RoundData[lastIdx]
			team2Data := game.Team2.RoundData[lastIdx]
			finalEconAdvantage = team1Data.Funds - team2Data.Funds
			finalEquipAdvantage = team1Data.FTE_Eq_value - team2Data.FTE_Eq_value
		}
		aa.recordStreak(game.ID, currentStreakTeam, currentStreak, streakStartRound, len(game.Rounds),
			streakStartEconAdvantage, finalEconAdvantage, streakStartEquipAdvantage, finalEquipAdvantage)
	}

	// Analyze comeback scenarios
	aa.analyzeComebacks(game.ID, game.Rounds, game.Team1.RoundData, game.Team2.RoundData, game.Is_T1_Winner)
}

// updateTimeSeries updates the time series data
func (aa *AdvancedAnalyzer) updateTimeSeries(roundNum int, team1Data *engine.Team_RoundData, team2Data *engine.Team_RoundData, team1Won bool) {
	if _, exists := aa.roundDataMap[roundNum]; !exists {
		aa.roundDataMap[roundNum] = &RoundTimePoint{
			RoundNumber: roundNum,
		}
	}

	point := aa.roundDataMap[roundNum]
	point.GamesReachedThisRound++

	// Accumulate values
	point.Team1AvgFunds += team1Data.Funds
	point.Team2AvgFunds += team2Data.Funds
	point.Team1AvgEquipment += team1Data.FTE_Eq_value
	point.Team2AvgEquipment += team2Data.FTE_Eq_value
	point.AvgEconomicAdvantage += (team1Data.Funds - team2Data.Funds)
	point.AvgEquipmentAdvantage += (team1Data.FTE_Eq_value - team2Data.FTE_Eq_value)
	point.Team1AvgSpent += team1Data.Spent
	point.Team2AvgSpent += team2Data.Spent
	point.Team1AvgEarned += team1Data.Earned
	point.Team2AvgEarned += team2Data.Earned
	point.Team1AvgSurvivors += float64(team1Data.Survivors)
	point.Team2AvgSurvivors += float64(team2Data.Survivors)

	if team1Won {
		point.Team1Wins++
	} else {
		point.Team2Wins++
	}
}

// recordStreak records a streak
func (aa *AdvancedAnalyzer) recordStreak(gameID, team string, length, startRound, endRound int,
	startEconAdvantage, endEconAdvantage, startEquipAdvantage, endEquipAdvantage float64) {

	streakInfo := StreakInfo{
		Length:             length,
		StartRound:         startRound,
		EndRound:           endRound,
		StartEconomicEdge:  startEconAdvantage,
		EndEconomicEdge:    endEconAdvantage,
		StartEquipmentEdge: startEquipAdvantage,
		EndEquipmentEdge:   endEquipAdvantage,
		EconomicChange:     endEconAdvantage - startEconAdvantage,
		EquipmentChange:    endEquipAdvantage - startEquipAdvantage,
		GameID:             gameID,
	}

	if team == "Team1" {
		aa.analysis.Streaks.Team1WinStreaks = append(aa.analysis.Streaks.Team1WinStreaks, streakInfo)
	} else {
		aa.analysis.Streaks.Team2WinStreaks = append(aa.analysis.Streaks.Team2WinStreaks, streakInfo)
	}
}

// analyzeComebacks analyzes comeback scenarios in a game
func (aa *AdvancedAnalyzer) analyzeComebacks(gameID string, rounds []engine.Round, team1Data []engine.Team_RoundData, team2Data []engine.Team_RoundData, team1Won bool) {
	team1Score := 0
	team2Score := 0

	for i, round := range rounds {
		if round.IsT1WinnerTeam {
			team1Score++
		} else {
			team2Score++
		}

		// Check for comeback scenarios
		deficit := team2Score - team1Score
		if deficit > 0 && team1Won {
			// Team1 came back from deficit
			economicAdvantages := []float64{}
			equipmentAdvantages := []float64{}

			for j := i; j < len(rounds); j++ {
				economicAdvantages = append(economicAdvantages, team1Data[j].Funds-team2Data[j].Funds)
				equipmentAdvantages = append(equipmentAdvantages, team1Data[j].FTE_Eq_value-team2Data[j].FTE_Eq_value)
			}

			avgEconAdvantage := 0.0
			avgEquipAdvantage := 0.0
			if len(economicAdvantages) > 0 {
				for _, v := range economicAdvantages {
					avgEconAdvantage += v
				}
				avgEconAdvantage /= float64(len(economicAdvantages))

				for _, v := range equipmentAdvantages {
					avgEquipAdvantage += v
				}
				avgEquipAdvantage /= float64(len(equipmentAdvantages))
			}

			aa.comebackScenarios[deficit] = append(aa.comebackScenarios[deficit], ComebackAttempt{
				GameID:                gameID,
				DeficitSize:           deficit,
				StartingRound:         i + 1,
				Team:                  "Team1",
				Success:               true,
				EconomicAdvantageAvg:  avgEconAdvantage,
				EquipmentAdvantageAvg: avgEquipAdvantage,
			})
		} else if deficit < 0 && !team1Won {
			// Team2 came back from deficit
			actualDeficit := -deficit
			economicAdvantages := []float64{}
			equipmentAdvantages := []float64{}

			for j := i; j < len(rounds); j++ {
				economicAdvantages = append(economicAdvantages, team2Data[j].Funds-team1Data[j].Funds)
				equipmentAdvantages = append(equipmentAdvantages, team2Data[j].FTE_Eq_value-team1Data[j].FTE_Eq_value)
			}

			avgEconAdvantage := 0.0
			avgEquipAdvantage := 0.0
			if len(economicAdvantages) > 0 {
				for _, v := range economicAdvantages {
					avgEconAdvantage += v
				}
				avgEconAdvantage /= float64(len(economicAdvantages))

				for _, v := range equipmentAdvantages {
					avgEquipAdvantage += v
				}
				avgEquipAdvantage /= float64(len(equipmentAdvantages))
			}

			aa.comebackScenarios[actualDeficit] = append(aa.comebackScenarios[actualDeficit], ComebackAttempt{
				GameID:                gameID,
				DeficitSize:           actualDeficit,
				StartingRound:         i + 1,
				Team:                  "Team2",
				Success:               true,
				EconomicAdvantageAvg:  avgEconAdvantage,
				EquipmentAdvantageAvg: avgEquipAdvantage,
			})
		}
	}
}

// Finalize computes all final statistics
func (aa *AdvancedAnalyzer) Finalize() *AdvancedAnalysis {
	// Calculate economic momentum
	aa.calculateEconomicMomentum()

	// Calculate win conditions
	aa.calculateWinConditions()

	// Calculate comeback statistics
	aa.calculateComebackStats()

	// Calculate half/side effects
	aa.calculateHalfSideEffects()

	// Finalize time series
	aa.finalizeTimeSeries()

	// Create distributions
	aa.createDistributions()

	// Finalize streak analysis
	aa.finalizeStreakAnalysis()

	return aa.analysis
}

// calculateEconomicMomentum calculates economic momentum indicators
func (aa *AdvancedAnalyzer) calculateEconomicMomentum() {
	if len(aa.economicAdvantages) == 0 {
		return
	}

	// Average advantages
	avgEconAdv := 0.0
	avgEquipAdv := 0.0
	for i := range aa.economicAdvantages {
		avgEconAdv += aa.economicAdvantages[i]
		avgEquipAdv += aa.equipmentAdvantages[i]
	}
	avgEconAdv /= float64(len(aa.economicAdvantages))
	avgEquipAdv /= float64(len(aa.equipmentAdvantages))

	aa.analysis.EconomicMomentum.AverageEconomicAdvantage = avgEconAdv
	aa.analysis.EconomicMomentum.AverageEquipmentDifferential = avgEquipAdv

	// Volatility
	aa.analysis.EconomicMomentum.EconomicVolatility = calculateStdDev(aa.economicAdvantages)
	aa.analysis.EconomicMomentum.EquipmentVolatility = calculateStdDev(aa.equipmentAdvantages)

	// Momentum shifts (count significant changes > 10000 funds or > 5000 equipment)
	momentumShifts := 0
	momentumDurations := []int{}
	currentDuration := 0
	lastAdvantageSign := 0

	for i := range aa.economicAdvantages {
		currentSign := 0
		if aa.economicAdvantages[i] > 5000 {
			currentSign = 1
		} else if aa.economicAdvantages[i] < -5000 {
			currentSign = -1
		}

		if currentSign != 0 {
			if currentSign == lastAdvantageSign {
				currentDuration++
			} else {
				if currentDuration > 0 {
					momentumDurations = append(momentumDurations, currentDuration)
				}
				currentDuration = 1
				if lastAdvantageSign != 0 {
					momentumShifts++
				}
				lastAdvantageSign = currentSign
			}
		}
	}

	if currentDuration > 0 {
		momentumDurations = append(momentumDurations, currentDuration)
	}

	aa.analysis.EconomicMomentum.MomentumShifts = momentumShifts

	if len(momentumDurations) > 0 {
		avgDuration := 0.0
		for _, d := range momentumDurations {
			avgDuration += float64(d)
		}
		aa.analysis.EconomicMomentum.AverageMomentumDuration = avgDuration / float64(len(momentumDurations))
	}

	// Spending efficiency
	totalTeam1Wins := 0.0
	totalTeam2Wins := 0.0
	totalTeam1Spent := 0.0
	totalTeam2Spent := 0.0

	for _, sample := range aa.winConditionSamples {
		totalTeam1Spent += sample.Team1Spent
		totalTeam2Spent += sample.Team2Spent
		if sample.Team1Won {
			totalTeam1Wins++
		} else {
			totalTeam2Wins++
		}
	}

	if totalTeam1Spent > 0 {
		aa.analysis.EconomicMomentum.Team1SpendEfficiency = totalTeam1Wins / (totalTeam1Spent / 1000)
	}
	if totalTeam2Spent > 0 {
		aa.analysis.EconomicMomentum.Team2SpendEfficiency = totalTeam2Wins / (totalTeam2Spent / 1000)
	}
}

// calculateWinConditions calculates win condition statistics
func (aa *AdvancedAnalyzer) calculateWinConditions() {
	if len(aa.winConditionSamples) == 0 {
		return
	}

	// Create continuous ranges for equipment advantage
	equipRanges := []struct{ min, max float64 }{
		{-math.MaxFloat64, -10000},
		{-10000, -5000},
		{-5000, -2500},
		{-2500, 0},
		{0, 2500},
		{2500, 5000},
		{5000, 10000},
		{10000, math.MaxFloat64},
	}

	for _, rng := range equipRanges {
		team1Wins := 0
		total := 0

		for _, sample := range aa.winConditionSamples {
			inRange := false
			if rng.max == math.MaxFloat64 {
				// For the last range, include everything >= min
				inRange = sample.EquipmentAdvantage >= rng.min
			} else {
				// For other ranges, use [min, max)
				inRange = sample.EquipmentAdvantage >= rng.min && sample.EquipmentAdvantage < rng.max
			}
			if inRange {
				total++
				if sample.Team1Won {
					team1Wins++
				}
			}
		}

		if total > 0 {
			aa.analysis.WinConditions.EquipmentAdvantageRanges = append(
				aa.analysis.WinConditions.EquipmentAdvantageRanges,
				EquipmentAdvantageRange{
					MinAdvantage: rng.min,
					MaxAdvantage: rng.max,
					Team1WinRate: float64(team1Wins) / float64(total),
					SampleSize:   total,
				},
			)
		}
	}

	// Create continuous ranges for economic advantage
	econRanges := []struct{ min, max float64 }{
		{-math.MaxFloat64, -15000},
		{-15000, -10000},
		{-10000, -5000},
		{-5000, -2500},
		{-2500, 0},
		{0, 2500},
		{2500, 5000},
		{5000, 10000},
		{10000, 15000},
		{15000, math.MaxFloat64},
	}

	for _, rng := range econRanges {
		team1Wins := 0
		total := 0

		for _, sample := range aa.winConditionSamples {
			inRange := false
			if rng.max == math.MaxFloat64 {
				// For the last range, include everything >= min
				inRange = sample.FundsAdvantage >= rng.min
			} else {
				// For other ranges, use [min, max)
				inRange = sample.FundsAdvantage >= rng.min && sample.FundsAdvantage < rng.max
			}
			if inRange {
				total++
				if sample.Team1Won {
					team1Wins++
				}
			}
		}

		if total > 0 {
			aa.analysis.WinConditions.EconomicAdvantageRanges = append(
				aa.analysis.WinConditions.EconomicAdvantageRanges,
				EconomicAdvantageRange{
					MinAdvantage: rng.min,
					MaxAdvantage: rng.max,
					Team1WinRate: float64(team1Wins) / float64(total),
					SampleSize:   total,
				},
			)
		}
	}

	// Calculate ROI
	team1Wins := 0.0
	team2Wins := 0.0
	team1Spent := 0.0
	team2Spent := 0.0

	for _, sample := range aa.winConditionSamples {
		team1Spent += sample.Team1Spent
		team2Spent += sample.Team2Spent
		if sample.Team1Won {
			team1Wins++
		} else {
			team2Wins++
		}
	}

	if team1Spent > 0 {
		aa.analysis.WinConditions.Team1EquipmentROI = team1Wins / (team1Spent / 1000)
	}
	if team2Spent > 0 {
		aa.analysis.WinConditions.Team2EquipmentROI = team2Wins / (team2Spent / 1000)
	}

	// Calculate correlations
	outcomes := make([]float64, len(aa.winConditionSamples))
	funds := make([]float64, len(aa.winConditionSamples))
	equipment := make([]float64, len(aa.winConditionSamples))
	consecLoss := make([]float64, len(aa.winConditionSamples))

	for i, sample := range aa.winConditionSamples {
		if sample.Team1Won {
			outcomes[i] = 1.0
		} else {
			outcomes[i] = 0.0
		}
		funds[i] = sample.FundsAdvantage
		equipment[i] = sample.EquipmentAdvantage
		consecLoss[i] = float64(sample.Team2ConsecLoss - sample.Team1ConsecLoss)
	}

	aa.analysis.WinConditions.FundsCorrelation = calculateCorrelation(funds, outcomes)
	aa.analysis.WinConditions.EquipmentCorrelation = calculateCorrelation(equipment, outcomes)
	aa.analysis.WinConditions.ConsecLossCorrelation = calculateCorrelation(consecLoss, outcomes)
}

// calculateComebackStats calculates comeback statistics
func (aa *AdvancedAnalyzer) calculateComebackStats() {
	for deficit, attempts := range aa.comebackScenarios {
		successes := 0
		totalEconAdvantage := 0.0
		totalEquipAdvantage := 0.0

		for _, attempt := range attempts {
			if attempt.Success {
				successes++
			}
			totalEconAdvantage += attempt.EconomicAdvantageAvg
			totalEquipAdvantage += attempt.EquipmentAdvantageAvg
		}

		scenario := ComebackScenario{
			Attempts:    len(attempts),
			Successes:   successes,
			SuccessRate: float64(successes) / float64(len(attempts)),
		}

		if len(attempts) > 0 {
			scenario.AvgEconomicEdge = totalEconAdvantage / float64(len(attempts))
			scenario.AvgEquipmentEdge = totalEquipAdvantage / float64(len(attempts))
		}

		aa.analysis.ComebackAnalysis.ComebacksByDeficit[deficit] = scenario
	}
}

// calculateHalfSideEffects calculates half and side effects
func (aa *AdvancedAnalyzer) calculateHalfSideEffects() {
	for _, round := range aa.halfSideRounds {
		// First half vs second half
		if round.IsFirstHalf {
			if round.Team1Won {
				aa.analysis.HalfSideEffects.FirstHalfTeam1Wins++
			} else {
				aa.analysis.HalfSideEffects.FirstHalfTeam2Wins++
			}
		} else {
			if round.Team1Won {
				aa.analysis.HalfSideEffects.SecondHalfTeam1Wins++
			} else {
				aa.analysis.HalfSideEffects.SecondHalfTeam2Wins++
			}
		}

		// Team1 side analysis
		if round.Team1IsCT {
			aa.analysis.HalfSideEffects.Team1CTRounds++
			aa.analysis.HalfSideEffects.Team1CTAvgFunds += round.Team1Funds
			aa.analysis.HalfSideEffects.Team1CTAvgEquipment += round.Team1Equipment
			if round.Team1Won {
				aa.analysis.HalfSideEffects.Team1CTWins++
			}
		} else {
			aa.analysis.HalfSideEffects.Team1TRounds++
			aa.analysis.HalfSideEffects.Team1TAvgFunds += round.Team1Funds
			aa.analysis.HalfSideEffects.Team1TAvgEquipment += round.Team1Equipment
			if round.Team1Won {
				aa.analysis.HalfSideEffects.Team1TWins++
			}
		}

		// Team2 side analysis (opposite of Team1)
		if !round.Team1IsCT {
			aa.analysis.HalfSideEffects.Team2CTRounds++
			aa.analysis.HalfSideEffects.Team2CTAvgFunds += round.Team2Funds
			aa.analysis.HalfSideEffects.Team2CTAvgEquipment += round.Team2Equipment
			if !round.Team1Won {
				aa.analysis.HalfSideEffects.Team2CTWins++
			}
		} else {
			aa.analysis.HalfSideEffects.Team2TRounds++
			aa.analysis.HalfSideEffects.Team2TAvgFunds += round.Team2Funds
			aa.analysis.HalfSideEffects.Team2TAvgEquipment += round.Team2Equipment
			if !round.Team1Won {
				aa.analysis.HalfSideEffects.Team2TWins++
			}
		}
	}

	// Calculate averages
	firstHalfTotal := aa.analysis.HalfSideEffects.FirstHalfTeam1Wins + aa.analysis.HalfSideEffects.FirstHalfTeam2Wins
	if firstHalfTotal > 0 {
		aa.analysis.HalfSideEffects.FirstHalfTeam1WinRate = float64(aa.analysis.HalfSideEffects.FirstHalfTeam1Wins) / float64(firstHalfTotal)
	}

	secondHalfTotal := aa.analysis.HalfSideEffects.SecondHalfTeam1Wins + aa.analysis.HalfSideEffects.SecondHalfTeam2Wins
	if secondHalfTotal > 0 {
		aa.analysis.HalfSideEffects.SecondHalfTeam1WinRate = float64(aa.analysis.HalfSideEffects.SecondHalfTeam1Wins) / float64(secondHalfTotal)
	}

	if aa.analysis.HalfSideEffects.Team1CTRounds > 0 {
		aa.analysis.HalfSideEffects.Team1CTWinRate = float64(aa.analysis.HalfSideEffects.Team1CTWins) / float64(aa.analysis.HalfSideEffects.Team1CTRounds)
		aa.analysis.HalfSideEffects.Team1CTAvgFunds /= float64(aa.analysis.HalfSideEffects.Team1CTRounds)
		aa.analysis.HalfSideEffects.Team1CTAvgEquipment /= float64(aa.analysis.HalfSideEffects.Team1CTRounds)
	}

	if aa.analysis.HalfSideEffects.Team1TRounds > 0 {
		aa.analysis.HalfSideEffects.Team1TWinRate = float64(aa.analysis.HalfSideEffects.Team1TWins) / float64(aa.analysis.HalfSideEffects.Team1TRounds)
		aa.analysis.HalfSideEffects.Team1TAvgFunds /= float64(aa.analysis.HalfSideEffects.Team1TRounds)
		aa.analysis.HalfSideEffects.Team1TAvgEquipment /= float64(aa.analysis.HalfSideEffects.Team1TRounds)
	}

	if aa.analysis.HalfSideEffects.Team2CTRounds > 0 {
		aa.analysis.HalfSideEffects.Team2CTWinRate = float64(aa.analysis.HalfSideEffects.Team2CTWins) / float64(aa.analysis.HalfSideEffects.Team2CTRounds)
		aa.analysis.HalfSideEffects.Team2CTAvgFunds /= float64(aa.analysis.HalfSideEffects.Team2CTRounds)
		aa.analysis.HalfSideEffects.Team2CTAvgEquipment /= float64(aa.analysis.HalfSideEffects.Team2CTRounds)
	}

	if aa.analysis.HalfSideEffects.Team2TRounds > 0 {
		aa.analysis.HalfSideEffects.Team2TWinRate = float64(aa.analysis.HalfSideEffects.Team2TWins) / float64(aa.analysis.HalfSideEffects.Team2TRounds)
		aa.analysis.HalfSideEffects.Team2TAvgFunds /= float64(aa.analysis.HalfSideEffects.Team2TRounds)
		aa.analysis.HalfSideEffects.Team2TAvgEquipment /= float64(aa.analysis.HalfSideEffects.Team2TRounds)
	}
}

// finalizeTimeSeries converts accumulated time series data to averages
func (aa *AdvancedAnalyzer) finalizeTimeSeries() {
	for _, point := range aa.roundDataMap {
		if point.GamesReachedThisRound > 0 {
			n := float64(point.GamesReachedThisRound)
			point.Team1AvgFunds /= n
			point.Team2AvgFunds /= n
			point.Team1AvgEquipment /= n
			point.Team2AvgEquipment /= n
			point.AvgEconomicAdvantage /= n
			point.AvgEquipmentAdvantage /= n
			point.Team1AvgSpent /= n
			point.Team2AvgSpent /= n
			point.Team1AvgEarned /= n
			point.Team2AvgEarned /= n
			point.Team1AvgSurvivors /= n
			point.Team2AvgSurvivors /= n

			total := point.Team1Wins + point.Team2Wins
			if total > 0 {
				point.Team1WinRate = float64(point.Team1Wins) / float64(total)
			}
		}

		aa.analysis.TimeSeries.RoundData = append(aa.analysis.TimeSeries.RoundData, *point)
	}

	// Sort by round number
	sort.Slice(aa.analysis.TimeSeries.RoundData, func(i, j int) bool {
		return aa.analysis.TimeSeries.RoundData[i].RoundNumber < aa.analysis.TimeSeries.RoundData[j].RoundNumber
	})
}

// createDistributions creates histogram distributions
func (aa *AdvancedAnalyzer) createDistributions() {
	numBins := 20

	aa.analysis.Distributions.Team1FundsDistribution = createHistogramBins(aa.allFundsTeam1, numBins)
	aa.analysis.Distributions.Team2FundsDistribution = createHistogramBins(aa.allFundsTeam2, numBins)
	aa.analysis.Distributions.Team1EquipmentDistribution = createHistogramBins(aa.allEquipmentTeam1, numBins)
	aa.analysis.Distributions.Team2EquipmentDistribution = createHistogramBins(aa.allEquipmentTeam2, numBins)
	aa.analysis.Distributions.EconomicAdvantageDistribution = createHistogramBins(aa.economicAdvantages, numBins)
	aa.analysis.Distributions.EquipmentAdvantageDistribution = createHistogramBins(aa.equipmentAdvantages, numBins)

	// Create win probability heatmap
	aa.createWinProbabilityHeatmap()
}

// createWinProbabilityHeatmap creates a 2D heatmap of win probability
func (aa *AdvancedAnalyzer) createWinProbabilityHeatmap() {
	// Define grid size
	gridSize := 10
	minEquip := 0.0
	maxEquip := 30000.0
	binSize := (maxEquip - minEquip) / float64(gridSize)

	// Create 2D grid
	heatmap := make([][]HeatmapCell, gridSize)
	for i := range heatmap {
		heatmap[i] = make([]HeatmapCell, gridSize)
		for j := range heatmap[i] {
			heatmap[i][j] = HeatmapCell{
				Team1Equipment: minEquip + float64(i)*binSize + binSize/2,
				Team2Equipment: minEquip + float64(j)*binSize + binSize/2,
			}
		}
	}

	// Populate grid
	for _, sample := range aa.winConditionSamples {
		team1Equip := sample.EquipmentAdvantage + sample.Team2Spent // Reconstruct Team1 equipment
		team2Equip := sample.Team2Spent

		i := int((team1Equip - minEquip) / binSize)
		j := int((team2Equip - minEquip) / binSize)

		if i >= 0 && i < gridSize && j >= 0 && j < gridSize {
			heatmap[i][j].SampleSize++
			if sample.Team1Won {
				heatmap[i][j].Team1WinRate += 1.0
			}
		}
	}

	// Calculate win rates
	for i := range heatmap {
		for j := range heatmap[i] {
			if heatmap[i][j].SampleSize > 0 {
				heatmap[i][j].Team1WinRate /= float64(heatmap[i][j].SampleSize)
			}
		}
	}

	aa.analysis.Distributions.WinProbabilityHeatmap = heatmap
}

// finalizeStreakAnalysis calculates streak statistics
func (aa *AdvancedAnalyzer) finalizeStreakAnalysis() {
	// Sort streaks by length
	sortStreaksByLength(aa.analysis.Streaks.Team1WinStreaks)
	sortStreaksByLength(aa.analysis.Streaks.Team2WinStreaks)

	// Calculate average streaks
	if len(aa.analysis.Streaks.Team1WinStreaks) > 0 {
		total := 0
		for _, streak := range aa.analysis.Streaks.Team1WinStreaks {
			total += streak.Length
		}
		aa.analysis.Streaks.Team1AvgWinStreak = float64(total) / float64(len(aa.analysis.Streaks.Team1WinStreaks))
		aa.analysis.Streaks.Team1MaxWinStreak = aa.analysis.Streaks.Team1WinStreaks[0].Length
	}

	if len(aa.analysis.Streaks.Team2WinStreaks) > 0 {
		total := 0
		for _, streak := range aa.analysis.Streaks.Team2WinStreaks {
			total += streak.Length
		}
		aa.analysis.Streaks.Team2AvgWinStreak = float64(total) / float64(len(aa.analysis.Streaks.Team2WinStreaks))
		aa.analysis.Streaks.Team2MaxWinStreak = aa.analysis.Streaks.Team2WinStreaks[0].Length
	}

	// Calculate streak economic impact
	streakImpactMap := make(map[int]*StreakEconomicImpact)

	for _, streak := range aa.analysis.Streaks.Team1WinStreaks {
		if _, exists := streakImpactMap[streak.Length]; !exists {
			streakImpactMap[streak.Length] = &StreakEconomicImpact{
				StreakLength: streak.Length,
			}
		}
		impact := streakImpactMap[streak.Length]
		impact.Occurrences++
		impact.AvgEconomicAdvantageGain += streak.EconomicChange
		impact.AvgEquipmentAdvantageGain += streak.EquipmentChange
	}

	for _, streak := range aa.analysis.Streaks.Team2WinStreaks {
		if _, exists := streakImpactMap[streak.Length]; !exists {
			streakImpactMap[streak.Length] = &StreakEconomicImpact{
				StreakLength: streak.Length,
			}
		}
		impact := streakImpactMap[streak.Length]
		impact.Occurrences++
		impact.AvgEconomicAdvantageGain += -streak.EconomicChange // Negative for Team2
		impact.AvgEquipmentAdvantageGain += -streak.EquipmentChange
	}

	// Calculate averages
	for _, impact := range streakImpactMap {
		if impact.Occurrences > 0 {
			impact.AvgEconomicAdvantageGain /= float64(impact.Occurrences)
			impact.AvgEquipmentAdvantageGain /= float64(impact.Occurrences)
		}
		aa.analysis.Streaks.StreakEconomicImpact = append(aa.analysis.Streaks.StreakEconomicImpact, *impact)
	}

	// Sort by streak length
	sort.Slice(aa.analysis.Streaks.StreakEconomicImpact, func(i, j int) bool {
		return aa.analysis.Streaks.StreakEconomicImpact[i].StreakLength < aa.analysis.Streaks.StreakEconomicImpact[j].StreakLength
	})
}
