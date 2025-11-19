package analysis

import (
	"math"
	"sort"
)

// AdvancedAnalysis contains deep economic and tactical analysis
type AdvancedAnalysis struct {
	// Economic Momentum
	EconomicMomentum EconomicMomentumAnalysis `json:"economic_momentum"`

	// Win Conditions
	WinConditions WinConditionAnalysis `json:"win_conditions"`

	// Comeback Analysis
	ComebackAnalysis ComebackAnalysis `json:"comeback_analysis"`

	// Half/Side Effects
	HalfSideEffects HalfSideAnalysis `json:"half_side_effects"`

	// Time Series Data
	TimeSeries TimeSeriesData `json:"time_series"`

	// Distribution Data
	Distributions DistributionData `json:"distributions"`

	// Streak Analysis
	Streaks StreakAnalysis `json:"streaks"`

	// Investment Ratio Analysis
	InvestmentRatio InvestmentRatioAnalysis `json:"investment_ratio"`
}

// EconomicMomentumAnalysis tracks economic flow and momentum
type EconomicMomentumAnalysis struct {
	// Average economic advantage per round (Team1 - Team2)
	AverageEconomicAdvantage float64 `json:"average_economic_advantage"`

	// Average equipment differential
	AverageEquipmentDifferential float64 `json:"average_equipment_differential"`

	// Economic volatility (standard deviation of advantage)
	EconomicVolatility  float64 `json:"economic_volatility"`
	EquipmentVolatility float64 `json:"equipment_volatility"`

	// Recovery metrics
	AverageRecoveryTime float64 `json:"average_recovery_time"` // Rounds to recover after loss
	RecoverySuccessRate float64 `json:"recovery_success_rate"` // % of successful recoveries

	// Momentum shifts
	MomentumShifts          int     `json:"momentum_shifts"`           // Major economic advantage changes
	AverageMomentumDuration float64 `json:"average_momentum_duration"` // Rounds of sustained advantage

	// Spending efficiency
	Team1SpendEfficiency float64 `json:"team1_spend_efficiency"` // Wins per dollar spent
	Team2SpendEfficiency float64 `json:"team2_spend_efficiency"`
}

// WinConditionAnalysis examines what leads to round victories
type WinConditionAnalysis struct {
	// Continuous equipment advantage analysis
	EquipmentAdvantageRanges []EquipmentAdvantageRange `json:"equipment_advantage_ranges"`

	// Economic advantage impact
	EconomicAdvantageRanges []EconomicAdvantageRange `json:"economic_advantage_ranges"`

	// ROI on equipment
	Team1EquipmentROI float64 `json:"team1_equipment_roi"` // Rounds won per $1000 spent
	Team2EquipmentROI float64 `json:"team2_equipment_roi"`

	// Win probability regression coefficients
	FundsCoefficient      float64 `json:"funds_coefficient"`
	EquipmentCoefficient  float64 `json:"equipment_coefficient"`
	ConsecLossCoefficient float64 `json:"consec_loss_coefficient"`
	SurvivorsCoefficient  float64 `json:"survivors_coefficient"`

	// Correlation strengths
	FundsCorrelation      float64 `json:"funds_correlation"`
	EquipmentCorrelation  float64 `json:"equipment_correlation"`
	ConsecLossCorrelation float64 `json:"consec_loss_correlation"`
}

// EquipmentAdvantageRange represents win rate for equipment differential ranges
type EquipmentAdvantageRange struct {
	MinAdvantage float64 `json:"min_advantage"`
	MaxAdvantage float64 `json:"max_advantage"`
	Team1WinRate float64 `json:"team1_win_rate"`
	SampleSize   int     `json:"sample_size"`
}

// EconomicAdvantageRange represents win rate for economic advantage ranges
type EconomicAdvantageRange struct {
	MinAdvantage float64 `json:"min_advantage"`
	MaxAdvantage float64 `json:"max_advantage"`
	Team1WinRate float64 `json:"team1_win_rate"`
	SampleSize   int     `json:"sample_size"`
}

// ComebackAnalysis tracks comeback scenarios
type ComebackAnalysis struct {
	// Comeback scenarios by score deficit
	ComebacksByDeficit map[int]ComebackScenario `json:"comebacks_by_deficit"` // key: deficit (1-5+)

	// Economic position during comebacks
	AverageComebackEconomicAdvantage  float64 `json:"average_comeback_economic_advantage"`
	AverageComebackEquipmentAdvantage float64 `json:"average_comeback_equipment_advantage"`

	// Overall comeback rates
	Team1ComebackRate float64 `json:"team1_comeback_rate"` // % of games won after being behind
	Team2ComebackRate float64 `json:"team2_comeback_rate"`

	// Critical round analysis
	CriticalRoundWinRate map[string]float64 `json:"critical_round_win_rate"` // Map round type to win rate
}

// ComebackScenario tracks comeback statistics for a specific deficit
type ComebackScenario struct {
	Attempts         int     `json:"attempts"`
	Successes        int     `json:"successes"`
	SuccessRate      float64 `json:"success_rate"`
	AvgEconomicEdge  float64 `json:"avg_economic_edge"`
	AvgEquipmentEdge float64 `json:"avg_equipment_edge"`
}

// HalfSideAnalysis examines performance by half and side
type HalfSideAnalysis struct {
	// First half performance
	FirstHalfTeam1Wins    int     `json:"first_half_team1_wins"`
	FirstHalfTeam2Wins    int     `json:"first_half_team2_wins"`
	FirstHalfTeam1WinRate float64 `json:"first_half_team1_win_rate"`

	// Second half performance
	SecondHalfTeam1Wins    int     `json:"second_half_team1_wins"`
	SecondHalfTeam2Wins    int     `json:"second_half_team2_wins"`
	SecondHalfTeam1WinRate float64 `json:"second_half_team1_win_rate"`

	// Side-specific performance (when Team1 is CT)
	Team1CTRounds       int     `json:"team1_ct_rounds"`
	Team1CTWins         int     `json:"team1_ct_wins"`
	Team1CTWinRate      float64 `json:"team1_ct_win_rate"`
	Team1CTAvgFunds     float64 `json:"team1_ct_avg_funds"`
	Team1CTAvgEquipment float64 `json:"team1_ct_avg_equipment"`

	// Side-specific performance (when Team1 is T)
	Team1TRounds       int     `json:"team1_t_rounds"`
	Team1TWins         int     `json:"team1_t_wins"`
	Team1TWinRate      float64 `json:"team1_t_win_rate"`
	Team1TAvgFunds     float64 `json:"team1_t_avg_funds"`
	Team1TAvgEquipment float64 `json:"team1_t_avg_equipment"`

	// Side-specific performance (when Team2 is CT)
	Team2CTRounds       int     `json:"team2_ct_rounds"`
	Team2CTWins         int     `json:"team2_ct_wins"`
	Team2CTWinRate      float64 `json:"team2_ct_win_rate"`
	Team2CTAvgFunds     float64 `json:"team2_ct_avg_funds"`
	Team2CTAvgEquipment float64 `json:"team2_ct_avg_equipment"`

	// Side-specific performance (when Team2 is T)
	Team2TRounds       int     `json:"team2_t_rounds"`
	Team2TWins         int     `json:"team2_t_wins"`
	Team2TWinRate      float64 `json:"team2_t_win_rate"`
	Team2TAvgFunds     float64 `json:"team2_t_avg_funds"`
	Team2TAvgEquipment float64 `json:"team2_t_avg_equipment"`

	// Half switch impact
	SideSwitchMomentumChange float64 `json:"side_switch_momentum_change"` // Economic advantage change at half
}

// TimeSeriesData contains round-by-round aggregated data across all games
type TimeSeriesData struct {
	RoundData []RoundTimePoint `json:"round_data"`
}

// RoundTimePoint represents aggregated data for a specific round number across games
type RoundTimePoint struct {
	RoundNumber int `json:"round_number"`

	// Economic metrics
	Team1AvgFunds     float64 `json:"team1_avg_funds"`
	Team2AvgFunds     float64 `json:"team2_avg_funds"`
	Team1AvgEquipment float64 `json:"team1_avg_equipment"`
	Team2AvgEquipment float64 `json:"team2_avg_equipment"`

	// Differentials
	AvgEconomicAdvantage  float64 `json:"avg_economic_advantage"` // Team1 - Team2
	AvgEquipmentAdvantage float64 `json:"avg_equipment_advantage"`

	// Outcomes
	Team1Wins    int     `json:"team1_wins"`
	Team2Wins    int     `json:"team2_wins"`
	Team1WinRate float64 `json:"team1_win_rate"`

	// Spending patterns
	Team1AvgSpent  float64 `json:"team1_avg_spent"`
	Team2AvgSpent  float64 `json:"team2_avg_spent"`
	Team1AvgEarned float64 `json:"team1_avg_earned"`
	Team2AvgEarned float64 `json:"team2_avg_earned"`

	// Survivors
	Team1AvgSurvivors float64 `json:"team1_avg_survivors"`
	Team2AvgSurvivors float64 `json:"team2_avg_survivors"`

	// Sample size
	GamesReachedThisRound int `json:"games_reached_this_round"`
}

// DistributionData contains distribution analysis
type DistributionData struct {
	// Funds distributions
	Team1FundsDistribution []FrequencyBin `json:"team1_funds_distribution"`
	Team2FundsDistribution []FrequencyBin `json:"team2_funds_distribution"`

	// Equipment distributions
	Team1EquipmentDistribution []FrequencyBin `json:"team1_equipment_distribution"`
	Team2EquipmentDistribution []FrequencyBin `json:"team2_equipment_distribution"`

	// Economic advantage distribution
	EconomicAdvantageDistribution  []FrequencyBin `json:"economic_advantage_distribution"`
	EquipmentAdvantageDistribution []FrequencyBin `json:"equipment_advantage_distribution"`

	// Win probability heatmap data (2D distribution)
	WinProbabilityHeatmap [][]HeatmapCell `json:"win_probability_heatmap"`

	// Spending decision scatter data
	SpendingDecisions []SpendingDecision `json:"spending_decisions"`
}

// FrequencyBin represents a histogram bin
type FrequencyBin struct {
	Min        float64 `json:"min"`
	Max        float64 `json:"max"`
	Frequency  int     `json:"frequency"`
	Percentage float64 `json:"percentage"`
}

// HeatmapCell represents a cell in the win probability heatmap
type HeatmapCell struct {
	Team1Equipment float64 `json:"team1_equipment"`
	Team2Equipment float64 `json:"team2_equipment"`
	Team1WinRate   float64 `json:"team1_win_rate"`
	SampleSize     int     `json:"sample_size"`
}

// SpendingDecision tracks a spending decision and its outcome
type SpendingDecision struct {
	AvailableFunds float64 `json:"available_funds"`
	AmountSpent    float64 `json:"amount_spent"`
	SpendRatio     float64 `json:"spend_ratio"`
	RoundWon       bool    `json:"round_won"`
	Team           string  `json:"team"`
}

// StreakAnalysis examines winning and losing streaks
type StreakAnalysis struct {
	// Win streaks
	Team1WinStreaks []StreakInfo `json:"team1_win_streaks"`
	Team2WinStreaks []StreakInfo `json:"team2_win_streaks"`

	// Loss streaks
	Team1LossStreaks []StreakInfo `json:"team1_loss_streaks"`
	Team2LossStreaks []StreakInfo `json:"team2_loss_streaks"`

	// Aggregate statistics
	Team1AvgWinStreak float64 `json:"team1_avg_win_streak"`
	Team1MaxWinStreak int     `json:"team1_max_win_streak"`
	Team2AvgWinStreak float64 `json:"team2_avg_win_streak"`
	Team2MaxWinStreak int     `json:"team2_max_win_streak"`

	Team1AvgLossStreak float64 `json:"team1_avg_loss_streak"`
	Team1MaxLossStreak int     `json:"team1_max_loss_streak"`
	Team2AvgLossStreak float64 `json:"team2_avg_loss_streak"`
	Team2MaxLossStreak int     `json:"team2_max_loss_streak"`

	// Streak impact on economics
	StreakEconomicImpact []StreakEconomicImpact `json:"streak_economic_impact"`

	// Momentum analysis
	MomentumShifts []MomentumShift `json:"momentum_shifts"`
}

// StreakInfo contains detailed information about a streak
type StreakInfo struct {
	Length             int     `json:"length"`
	StartRound         int     `json:"start_round"`
	EndRound           int     `json:"end_round"`
	StartEconomicEdge  float64 `json:"start_economic_edge"`
	EndEconomicEdge    float64 `json:"end_economic_edge"`
	StartEquipmentEdge float64 `json:"start_equipment_edge"`
	EndEquipmentEdge   float64 `json:"end_equipment_edge"`
	EconomicChange     float64 `json:"economic_change"`
	EquipmentChange    float64 `json:"equipment_change"`
	GameID             string  `json:"game_id"`
}

// StreakEconomicImpact shows how streak length affects economics
type StreakEconomicImpact struct {
	StreakLength              int     `json:"streak_length"`
	Occurrences               int     `json:"occurrences"`
	AvgEconomicAdvantageGain  float64 `json:"avg_economic_advantage_gain"`
	AvgEquipmentAdvantageGain float64 `json:"avg_equipment_advantage_gain"`
	NextRoundWinRate          float64 `json:"next_round_win_rate"`
}

// MomentumShift represents a significant momentum change
type MomentumShift struct {
	GameID           string  `json:"game_id"`
	RoundNumber      int     `json:"round_number"`
	FromTeam         string  `json:"from_team"`
	ToTeam           string  `json:"to_team"`
	EconomicSwing    float64 `json:"economic_swing"`
	EquipmentSwing   float64 `json:"equipment_swing"`
	ScoreBeforeShift [2]int  `json:"score_before_shift"`
	ScoreAfterShift  [2]int  `json:"score_after_shift"`
}

// InvestmentRatioAnalysis tracks investment ratio (spent/funds_start) for each team
type InvestmentRatioAnalysis struct {
	// Overall statistics
	Team1AverageRatio float64 `json:"team1_average_ratio"`
	Team1MedianRatio  float64 `json:"team1_median_ratio"`
	Team2AverageRatio float64 `json:"team2_average_ratio"`
	Team2MedianRatio  float64 `json:"team2_median_ratio"`

	// Distribution by 10% steps (0-10%, 10-20%, ..., 90-100%)
	Team1Distribution []InvestmentRatioBin `json:"team1_distribution"`
	Team2Distribution []InvestmentRatioBin `json:"team2_distribution"`

	// Round-by-round investment ratio time series
	RoundInvestmentData []RoundInvestmentPoint `json:"round_investment_data"`
}

// InvestmentRatioBin represents a 10% range bin for investment ratios
type InvestmentRatioBin struct {
	MinRatio   float64 `json:"min_ratio"`  // e.g., 0.0, 0.1, 0.2, ..., 0.9
	MaxRatio   float64 `json:"max_ratio"`  // e.g., 0.1, 0.2, 0.3, ..., 1.0
	Frequency  int     `json:"frequency"`  // Number of rounds in this bin
	Percentage float64 `json:"percentage"` // Percentage of total rounds
	AvgRatio   float64 `json:"avg_ratio"`  // Average ratio within this bin
	WinRate    float64 `json:"win_rate"`   // Win rate for rounds in this bin
}

// RoundInvestmentPoint represents aggregated investment data for a specific round number
type RoundInvestmentPoint struct {
	RoundNumber                int     `json:"round_number"`
	Team1AvgInvestmentRatio    float64 `json:"team1_avg_investment_ratio"`
	Team2AvgInvestmentRatio    float64 `json:"team2_avg_investment_ratio"`
	Team1MedianInvestmentRatio float64 `json:"team1_median_investment_ratio"`
	Team2MedianInvestmentRatio float64 `json:"team2_median_investment_ratio"`
	GamesReachedThisRound      int     `json:"games_reached_this_round"`
}

// NewAdvancedAnalysis creates a new advanced analysis structure
func NewAdvancedAnalysis() *AdvancedAnalysis {
	return &AdvancedAnalysis{
		EconomicMomentum: EconomicMomentumAnalysis{},
		WinConditions: WinConditionAnalysis{
			EquipmentAdvantageRanges: make([]EquipmentAdvantageRange, 0),
			EconomicAdvantageRanges:  make([]EconomicAdvantageRange, 0),
		},
		ComebackAnalysis: ComebackAnalysis{
			ComebacksByDeficit:   make(map[int]ComebackScenario),
			CriticalRoundWinRate: make(map[string]float64),
		},
		HalfSideEffects: HalfSideAnalysis{},
		TimeSeries: TimeSeriesData{
			RoundData: make([]RoundTimePoint, 0),
		},
		Distributions: DistributionData{
			Team1FundsDistribution:         make([]FrequencyBin, 0),
			Team2FundsDistribution:         make([]FrequencyBin, 0),
			Team1EquipmentDistribution:     make([]FrequencyBin, 0),
			Team2EquipmentDistribution:     make([]FrequencyBin, 0),
			EconomicAdvantageDistribution:  make([]FrequencyBin, 0),
			EquipmentAdvantageDistribution: make([]FrequencyBin, 0),
			WinProbabilityHeatmap:          make([][]HeatmapCell, 0),
			SpendingDecisions:              make([]SpendingDecision, 0),
		},
		Streaks: StreakAnalysis{
			Team1WinStreaks:      make([]StreakInfo, 0),
			Team2WinStreaks:      make([]StreakInfo, 0),
			Team1LossStreaks:     make([]StreakInfo, 0),
			Team2LossStreaks:     make([]StreakInfo, 0),
			StreakEconomicImpact: make([]StreakEconomicImpact, 0),
			MomentumShifts:       make([]MomentumShift, 0),
		},
		InvestmentRatio: InvestmentRatioAnalysis{
			Team1Distribution:   make([]InvestmentRatioBin, 0),
			Team2Distribution:   make([]InvestmentRatioBin, 0),
			RoundInvestmentData: make([]RoundInvestmentPoint, 0),
		},
	}
}

// Helper function to calculate standard deviation
func calculateStdDev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values))

	return math.Sqrt(variance)
}

// Helper function to calculate correlation coefficient
func calculateCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}

	n := float64(len(x))

	// Calculate means
	meanX, meanY := 0.0, 0.0
	for i := range x {
		meanX += x[i]
		meanY += y[i]
	}
	meanX /= n
	meanY /= n

	// Calculate correlation
	numerator := 0.0
	denomX := 0.0
	denomY := 0.0

	for i := range x {
		diffX := x[i] - meanX
		diffY := y[i] - meanY
		numerator += diffX * diffY
		denomX += diffX * diffX
		denomY += diffY * diffY
	}

	if denomX == 0 || denomY == 0 {
		return 0
	}

	return numerator / math.Sqrt(denomX*denomY)
}

// Helper function to create histogram bins
func createHistogramBins(values []float64, numBins int) []FrequencyBin {
	if len(values) == 0 {
		return []FrequencyBin{}
	}

	// Find min and max
	minVal := values[0]
	maxVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	// Create bins
	binWidth := (maxVal - minVal) / float64(numBins)
	bins := make([]FrequencyBin, numBins)

	for i := 0; i < numBins; i++ {
		bins[i] = FrequencyBin{
			Min: minVal + float64(i)*binWidth,
			Max: minVal + float64(i+1)*binWidth,
		}
	}

	// Count frequencies
	for _, v := range values {
		for i := range bins {
			if v >= bins[i].Min && (v < bins[i].Max || (i == numBins-1 && v == bins[i].Max)) {
				bins[i].Frequency++
				break
			}
		}
	}

	// Calculate percentages
	total := float64(len(values))
	for i := range bins {
		bins[i].Percentage = float64(bins[i].Frequency) / total * 100
	}

	return bins
}

// Helper function to sort streak info by length descending
func sortStreaksByLength(streaks []StreakInfo) {
	sort.Slice(streaks, func(i, j int) bool {
		return streaks[i].Length > streaks[j].Length
	})
}
