package analysis

import (
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
		if team1Won {
			atomic.AddInt64(&s.Team1OTWins, 1)
		} else {
			atomic.AddInt64(&s.Team2OTWins, 1)
		}
	} else {
		if team1Won {
			atomic.AddInt64(&s.Team1RTWins, 1)
		} else {
			atomic.AddInt64(&s.Team2RTWins, 1)
		}
	}

}

// UpdateFailedSimulation increments the failed simulation counter
func (s *SimulationStats) UpdateFailedSimulation() {
	s.FailedSims++
}

// TeamGameEconomics holds economic statistics for a single game
type TeamGameEconomics struct {
	TotalSpent       float64
	TotalEarned      float64
	AverageFunds     float64
	AverageRSEq      float64
	AverageFTEEq     float64
	AverageREEq      float64
	AverageSurvivors float64
	MaxFunds         float64
	MinFunds         float64
	MaxConsecLosses  int
}

// CalculateFinalStats computes all derived statistics
func (s *SimulationStats) CalculateFinalStats() {
	if s.CompletedSims > 0 {
		s.Team1WinRate = float64(s.Team1Wins) / float64(s.CompletedSims) * 100
		s.Team2WinRate = float64(s.Team2Wins) / float64(s.CompletedSims) * 100
		s.Team1OTWinRate = float64(s.Team1OTWins) / float64(s.OvertimeGames) * 100
		s.Team2OTWinRate = float64(s.Team2OTWins) / float64(s.OvertimeGames) * 100
		s.Team1RTWinRate = float64(s.Team1RTWins) / float64(s.CompletedSims-s.OvertimeGames) * 100
		s.Team2RTWinRate = float64(s.Team2RTWins) / float64(s.CompletedSims-s.OvertimeGames) * 100

		s.OvertimeRate = float64(s.OvertimeGames) / float64(s.CompletedSims) * 100
		s.AverageRounds = float64(s.TotalRounds) / float64(s.CompletedSims)

		if s.ExecutionTime > 0 {
			s.ProcessingRate = float64(s.CompletedSims) / s.ExecutionTime.Seconds()
		}

	}
}
