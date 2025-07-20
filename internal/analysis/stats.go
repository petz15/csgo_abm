package analysis

import (
	"fmt"
	"runtime"
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

	// Calculated metrics
	Team1WinRate  float64 `json:"team1_win_rate"`
	Team2WinRate  float64 `json:"team2_win_rate"`
	AverageRounds float64 `json:"average_rounds"`
	OvertimeRate  float64 `json:"overtime_rate"`

	// Performance metrics (optional for sequential)
	ExecutionTime   time.Duration `json:"execution_time"`
	ProcessingRate  float64       `json:"simulations_per_second,omitempty"`
	PeakMemoryUsage uint64        `json:"peak_memory_usage_mb,omitempty"`
	TotalGCRuns     uint32        `json:"total_gc_runs,omitempty"`

	// Analysis data
	ScoreDistribution map[string]int64 `json:"score_distribution"`
	RoundDistribution map[int]int64    `json:"round_distribution"`

	// Advanced analysis
	AdvancedStats *AdvancedStats `json:"advanced_stats,omitempty"`

	// Metadata
	SimulationMode string    `json:"simulation_mode"` // "sequential" or "concurrent"
	StartTime      time.Time `json:"start_time"`

	// Thread safety (only needed for concurrent)
	ScoreMutex sync.Mutex `json:"-"`
	RoundMutex sync.Mutex `json:"-"`
}

// AdvancedStats contains detailed statistical analysis
type AdvancedStats struct {
	// Win streak analysis
	MaxWinStreak map[string]int     `json:"max_win_streaks"`
	AvgWinStreak map[string]float64 `json:"avg_win_streaks"`

	// Round analysis
	MedianRounds float64 `json:"median_rounds"`
	StdDevRounds float64 `json:"std_dev_rounds"`

	// Score analysis
	TopScoreLines []ScoreLine `json:"top_score_lines"`
	BlowoutGames  int64       `json:"blowout_games"` // >10 round difference
	CloseGames    int64       `json:"close_games"`   // <=3 round difference

	// Game balance analysis
	BalanceScore            float64 `json:"balance_score"`            // 0-100, higher = more balanced
	StatisticalSignificance string  `json:"statistical_significance"` // "significant", "not_significant", "insufficient_data"
	ChiSquareValue          float64 `json:"chi_square_value"`

	// Performance analysis (for concurrent mode)
	P50ResponseTime time.Duration   `json:"p50_response_time,omitempty"`
	P95ResponseTime time.Duration   `json:"p95_response_time,omitempty"`
	P99ResponseTime time.Duration   `json:"p99_response_time,omitempty"`
	ResponseTimes   []time.Duration `json:"-"` // For percentile calculations
}

// ScoreLine represents a score with its frequency
type ScoreLine struct {
	Score     string  `json:"score"`
	Count     int64   `json:"count"`
	Frequency float64 `json:"frequency"`
}

// RoundStats contains round distribution statistics
type RoundStats struct {
	Min          int           `json:"min"`
	Max          int           `json:"max"`
	Median       float64       `json:"median"`
	Mean         float64       `json:"mean"`
	StdDev       float64       `json:"std_dev"`
	Distribution map[int]int64 `json:"distribution"`
}

// SimulationConfig unified configuration for all simulation types
type SimulationConfig struct {
	NumSimulations int    `json:"num_simulations"`
	MaxConcurrent  int    `json:"max_concurrent,omitempty"` // Only for concurrent
	MemoryLimit    int    `json:"memory_limit,omitempty"`   // Only for concurrent
	Team1Name      string `json:"team1_name"`
	Team1Strategy  string `json:"team1_strategy"`
	Team2Name      string `json:"team2_name"`
	Team2Strategy  string `json:"team2_strategy"`
	GameRules      string `json:"game_rules"`
	ExportResults  bool   `json:"export_results"`
	Sequential     bool   `json:"sequential"`
}

// Validate validates the simulation configuration
func (c *SimulationConfig) Validate() error {
	if c.NumSimulations <= 0 {
		return fmt.Errorf("number of simulations must be positive")
	}
	if !c.Sequential && c.MaxConcurrent <= 0 {
		c.MaxConcurrent = runtime.NumCPU()
	}
	if c.Team1Strategy == "" {
		c.Team1Strategy = "all_in"
	}
	if c.Team2Strategy == "" {
		c.Team2Strategy = "default_half"
	}
	if c.GameRules == "" {
		c.GameRules = "default"
	}
	return nil
}

// NewStats creates a new SimulationStats instance
func NewStats(numSimulations int, mode string) *SimulationStats {
	return &SimulationStats{
		TotalSimulations:  int64(numSimulations),
		ScoreDistribution: make(map[string]int64),
		RoundDistribution: make(map[int]int64),
		SimulationMode:    mode,
		StartTime:         time.Now(),
		AdvancedStats: &AdvancedStats{
			MaxWinStreak:  make(map[string]int),
			AvgWinStreak:  make(map[string]float64),
			TopScoreLines: make([]ScoreLine, 0),
			ResponseTimes: make([]time.Duration, 0),
		},
	}
}

// NewSimulationStats creates a new SimulationStats from a config
func NewSimulationStats(config SimulationConfig) *SimulationStats {
	mode := "concurrent"
	if config.Sequential {
		mode = "sequential"
	}

	stats := NewStats(config.NumSimulations, mode)

	// Store configuration metadata
	if stats.AdvancedStats != nil {
		// Could store additional config metadata here if needed
	}

	return stats
}
