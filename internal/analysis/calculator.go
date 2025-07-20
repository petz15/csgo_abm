package analysis

import (
	"fmt"
	"math"
	"sort"
	"sync/atomic"
	"time"
)

// UpdateGameResult updates statistics with a single game result (thread-safe)
func (s *SimulationStats) UpdateGameResult(team1Won bool, team1Score, team2Score, totalRounds int, wentToOvertime bool, responseTime time.Duration) {
	// Update core stats atomically for thread safety
	atomic.AddInt64(&s.CompletedSims, 1)
	atomic.AddInt64(&s.TotalRounds, int64(totalRounds))

	if team1Won {
		atomic.AddInt64(&s.Team1Wins, 1)
	} else {
		atomic.AddInt64(&s.Team2Wins, 1)
	}

	if wentToOvertime {
		atomic.AddInt64(&s.OvertimeGames, 1)
	}

	// Update distributions (with mutex protection)
	scoreKey := fmt.Sprintf("%d-%d", team1Score, team2Score)
	s.ScoreMutex.Lock()
	// Limit score distribution to prevent memory issues
	if len(s.ScoreDistribution) < 1000 || s.ScoreDistribution[scoreKey] > 0 {
		s.ScoreDistribution[scoreKey]++
	}
	s.ScoreMutex.Unlock()

	s.RoundMutex.Lock()
	s.RoundDistribution[totalRounds]++
	s.RoundMutex.Unlock()

	// Update advanced stats
	s.updateAdvancedStats(team1Won, team1Score, team2Score, totalRounds, responseTime)
}

// UpdateFailedSimulation increments the failed simulation counter
func (s *SimulationStats) UpdateFailedSimulation() {
	s.FailedSims++
}

// CalculateFinalStats computes all derived statistics
func (s *SimulationStats) CalculateFinalStats() {
	if s.CompletedSims > 0 {
		s.Team1WinRate = float64(s.Team1Wins) / float64(s.CompletedSims) * 100
		s.Team2WinRate = float64(s.Team2Wins) / float64(s.CompletedSims) * 100
		s.OvertimeRate = float64(s.OvertimeGames) / float64(s.CompletedSims) * 100
		s.AverageRounds = float64(s.TotalRounds) / float64(s.CompletedSims)

		if s.ExecutionTime > 0 {
			s.ProcessingRate = float64(s.CompletedSims) / s.ExecutionTime.Seconds()
		}
	}

	// Calculate advanced statistics
	s.calculateAdvancedStats()
}

// updateAdvancedStats updates advanced statistics for each game
func (s *SimulationStats) updateAdvancedStats(team1Won bool, team1Score, team2Score, totalRounds int, responseTime time.Duration) {
	if s.AdvancedStats == nil {
		return
	}

	// Track response times for concurrent mode
	if s.SimulationMode == "concurrent" && responseTime > 0 {
		s.AdvancedStats.ResponseTimes = append(s.AdvancedStats.ResponseTimes, responseTime)
	}

	// Analyze game closeness
	scoreDiff := absInt(team1Score - team2Score)
	if scoreDiff <= 3 {
		s.AdvancedStats.CloseGames++
	} else if scoreDiff > 10 {
		s.AdvancedStats.BlowoutGames++
	}
}

// calculateAdvancedStats computes final advanced statistics
func (s *SimulationStats) calculateAdvancedStats() {
	if s.AdvancedStats == nil || s.CompletedSims == 0 {
		return
	}

	// Calculate balance score
	winDiff := abs(s.Team1Wins - s.Team2Wins)
	s.AdvancedStats.BalanceScore = (1.0 - float64(winDiff)/float64(s.CompletedSims)) * 100

	// Calculate statistical significance
	s.AdvancedStats.ChiSquareValue, s.AdvancedStats.StatisticalSignificance = calculateStatisticalSignificance(s.Team1Wins, s.Team2Wins)

	// Calculate round statistics
	s.calculateRoundStats()

	// Calculate score line analysis
	s.AdvancedStats.TopScoreLines = analyzeScoreDistribution(s.ScoreDistribution, s.CompletedSims)

	// Calculate response time percentiles for concurrent mode
	if s.SimulationMode == "concurrent" && len(s.AdvancedStats.ResponseTimes) > 0 {
		s.calculateResponseTimePercentiles()
	}
}

// calculateRoundStats computes round distribution statistics
func (s *SimulationStats) calculateRoundStats() {
	if len(s.RoundDistribution) == 0 {
		return
	}

	rounds := make([]int, 0, s.CompletedSims)
	for round, count := range s.RoundDistribution {
		for i := int64(0); i < count; i++ {
			rounds = append(rounds, round)
		}
	}

	if len(rounds) > 0 {
		sort.Ints(rounds)
		s.AdvancedStats.MedianRounds = calculateMedian(rounds)
		s.AdvancedStats.StdDevRounds = calculateStdDev(rounds)
	}
}

// calculateResponseTimePercentiles computes response time percentiles
func (s *SimulationStats) calculateResponseTimePercentiles() {
	if len(s.AdvancedStats.ResponseTimes) == 0 {
		return
	}

	times := make([]time.Duration, len(s.AdvancedStats.ResponseTimes))
	copy(times, s.AdvancedStats.ResponseTimes)
	sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })

	s.AdvancedStats.P50ResponseTime = times[len(times)*50/100]
	s.AdvancedStats.P95ResponseTime = times[len(times)*95/100]
	s.AdvancedStats.P99ResponseTime = times[len(times)*99/100]
}

// Helper functions

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func calculateMedian(sorted []int) float64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	if n%2 == 0 {
		return float64(sorted[n/2-1]+sorted[n/2]) / 2
	}
	return float64(sorted[n/2])
}

func calculateStdDev(values []int) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := calculateMean(values)
	var sum float64
	for _, v := range values {
		diff := float64(v) - mean
		sum += diff * diff
	}
	return math.Sqrt(sum / float64(len(values)))
}

func calculateMean(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0
	for _, v := range values {
		sum += v
	}
	return float64(sum) / float64(len(values))
}

func analyzeScoreDistribution(scores map[string]int64, totalGames int64) []ScoreLine {
	scoreLines := make([]ScoreLine, 0, len(scores))

	for score, count := range scores {
		scoreLines = append(scoreLines, ScoreLine{
			Score:     score,
			Count:     count,
			Frequency: float64(count) / float64(totalGames) * 100,
		})
	}

	// Sort by frequency (descending)
	sort.Slice(scoreLines, func(i, j int) bool {
		return scoreLines[i].Count > scoreLines[j].Count
	})

	// Return top 10 most common scores
	if len(scoreLines) > 10 {
		scoreLines = scoreLines[:10]
	}

	return scoreLines
}

func calculateStatisticalSignificance(team1Wins, team2Wins int64) (float64, string) {
	total := team1Wins + team2Wins
	if total < 30 {
		return 0, "insufficient_data"
	}

	// Simple chi-square test for equal proportions
	expected := float64(total) / 2
	chiSquare := math.Pow(float64(team1Wins)-expected, 2)/expected +
		math.Pow(float64(team2Wins)-expected, 2)/expected

	if chiSquare > 3.84 { // p < 0.05
		return chiSquare, "significant"
	}
	return chiSquare, "not_significant"
}
