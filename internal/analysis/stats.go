package analysis

import (
	"csgo_abm/internal/engine"
	"sync"
	"time"
)

// SimulationStats unified statistics structure for both sequential and concurrent simulations
type SimulationStats struct {
	// Core metrics
	TotalSimulations int64 `json:"total_simulations"`
	CompletedSims    int64 `json:"completed_simulations"`
	FailedSims       int64 `json:"failed_simulations"`

	// Game results
	Team1Wins     int64 `json:"team1_wins"`
	Team2Wins     int64 `json:"team2_wins"`
	TotalRounds   int64 `json:"total_rounds"`
	OvertimeGames int64 `json:"overtime_games"`
	Team1OTWins   int64 `json:"team1_overtime_wins"`
	Team2OTWins   int64 `json:"team2_overtime_wins"`
	Team1RTWins   int64 `json:"team1_regular_time_wins"`
	Team2RTWins   int64 `json:"team2_regular_time_wins"`

	// Calculated metrics
	Team1WinRate   float64 `json:"team1_win_rate"`
	Team2WinRate   float64 `json:"team2_win_rate"`
	AverageRounds  float64 `json:"average_rounds"`
	OvertimeRate   float64 `json:"overtime_rate"`
	Team1OTWinRate float64 `json:"team1_overtime_win_rate"`
	Team2OTWinRate float64 `json:"team2_overtime_win_rate"`
	Team1RTWinRate float64 `json:"team1_regular_time_win_rate"`
	Team2RTWinRate float64 `json:"team2_regular_time_win_rate"`

	// Performance metrics (optional for sequential)
	ExecutionTime   time.Duration `json:"execution_time"`
	ProcessingRate  float64       `json:"simulations_per_second,omitempty"`
	PeakMemoryUsage uint64        `json:"peak_memory_usage_mb,omitempty"`
	TotalGCRuns     uint32        `json:"total_gc_runs,omitempty"`

	// Metadata
	SimulationMode string            `json:"simulation_mode"` // "sequential" or "concurrent"
	StartTime      time.Time         `json:"start_time"`
	Config         *SimulationConfig `json:"simulation_config,omitempty"` // Configuration used for this simulation

	// Thread safety (only needed for concurrent)
	ScoreMutex sync.Mutex `json:"-"`
	RoundMutex sync.Mutex `json:"-"`
}

// ScoreLine represents a score with its frequency
type ScoreLine struct {
	Score     string  `json:"score"`
	Count     int64   `json:"count"`
	Frequency float64 `json:"frequency"`
}

type SimulationConfig struct {
	NumSimulations        int              `json:"num_simulations"`
	MaxConcurrent         int              `json:"max_concurrent,omitempty"` // Only for concurrent
	MemoryLimit           int              `json:"memory_limit,omitempty"`   // Only for concurrent
	Team1Name             string           `json:"team1_name"`
	Team1Strategy         string           `json:"team1_strategy"`
	Team2Name             string           `json:"team2_name"`
	Team2Strategy         string           `json:"team2_strategy"`
	GameRules             engine.GameRules `json:"game_rules"`
	ExportDetailedResults bool             `json:"export_detailed_results"`
	ExportRounds          bool             `json:"export_rounds"`
	Sequential            bool             `json:"sequential"`
	Exportpath            string           `json:"export_path,omitempty"` // Path for exporting results
}

// NewStats creates a new SimulationStats instance
func NewStats(NumSimulations int, mode string) *SimulationStats {
	return &SimulationStats{
		// Core metrics
		TotalSimulations: int64(NumSimulations),
		CompletedSims:    0,
		FailedSims:       0,

		// Game results
		Team1Wins:     0,
		Team2Wins:     0,
		TotalRounds:   0,
		OvertimeGames: 0,
		Team1OTWins:   0,
		Team2OTWins:   0,
		Team1RTWins:   0,
		Team2RTWins:   0,

		// Calculated metrics
		Team1WinRate:   0.0,
		Team2WinRate:   0.0,
		AverageRounds:  0.0,
		OvertimeRate:   0.0,
		Team1OTWinRate: 0.0,
		Team2OTWinRate: 0.0,
		Team1RTWinRate: 0.0,
		Team2RTWinRate: 0.0,

		// Performance metrics
		ProcessingRate:  0.0,
		PeakMemoryUsage: 0,
		TotalGCRuns:     0,

		// Metadata
		SimulationMode: mode,
		StartTime:      time.Now(),
	}
}

// NewSimulationStats creates a new SimulationStats from a config
func NewSimulationStats(config SimulationConfig) *SimulationStats {
	mode := "concurrent"
	if config.Sequential {
		mode = "sequential"
	}

	stats := NewStats(config.NumSimulations, mode)

	// Store configuration for export
	configCopy := config
	stats.Config = &configCopy

	return stats
}
