package analysis

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GraphExporter exports analysis data in formats suitable for graphing
type GraphExporter struct {
	analysis   *AdvancedAnalysis
	outputPath string
}

// NewGraphExporter creates a new graph exporter
func NewGraphExporter(analysis *AdvancedAnalysis, outputPath string) *GraphExporter {
	return &GraphExporter{
		analysis:   analysis,
		outputPath: outputPath,
	}
}

// ExportAll exports all graph data
func (ge *GraphExporter) ExportAll() error {
	// Create graphs subdirectory
	graphsDir := filepath.Join(ge.outputPath, "graphs")
	if err := os.MkdirAll(graphsDir, 0755); err != nil {
		return fmt.Errorf("failed to create graphs directory: %w", err)
	}

	// Export time series data
	if err := ge.ExportTimeSeriesCSV(); err != nil {
		return fmt.Errorf("failed to export time series: %w", err)
	}

	// Export distributions
	if err := ge.ExportDistributionsCSV(); err != nil {
		return fmt.Errorf("failed to export distributions: %w", err)
	}

	// Export streak analysis
	if err := ge.ExportStreakAnalysisCSV(); err != nil {
		return fmt.Errorf("failed to export streak analysis: %w", err)
	}

	// Export win conditions
	if err := ge.ExportWinConditionsCSV(); err != nil {
		return fmt.Errorf("failed to export win conditions: %w", err)
	}

	// Export spending decisions
	if err := ge.ExportSpendingDecisionsCSV(); err != nil {
		return fmt.Errorf("failed to export spending decisions: %w", err)
	}

	// Export heatmap
	if err := ge.ExportWinProbabilityHeatmapCSV(); err != nil {
		return fmt.Errorf("failed to export heatmap: %w", err)
	}

	// Export complete analysis as JSON
	if err := ge.ExportCompleteAnalysisJSON(); err != nil {
		return fmt.Errorf("failed to export complete analysis: %w", err)
	}

	return nil
}

// ExportTimeSeriesCSV exports time series data
func (ge *GraphExporter) ExportTimeSeriesCSV() error {
	path := filepath.Join(ge.outputPath, "graphs", "time_series.csv")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"round_number",
		"team1_avg_funds",
		"team2_avg_funds",
		"team1_avg_equipment",
		"team2_avg_equipment",
		"avg_economic_advantage",
		"avg_equipment_advantage",
		"team1_wins",
		"team2_wins",
		"team1_win_rate",
		"team1_avg_spent",
		"team2_avg_spent",
		"team1_avg_earned",
		"team2_avg_earned",
		"team1_avg_survivors",
		"team2_avg_survivors",
		"games_reached_this_round",
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, point := range ge.analysis.TimeSeries.RoundData {
		row := []string{
			fmt.Sprintf("%d", point.RoundNumber),
			fmt.Sprintf("%.2f", point.Team1AvgFunds),
			fmt.Sprintf("%.2f", point.Team2AvgFunds),
			fmt.Sprintf("%.2f", point.Team1AvgEquipment),
			fmt.Sprintf("%.2f", point.Team2AvgEquipment),
			fmt.Sprintf("%.2f", point.AvgEconomicAdvantage),
			fmt.Sprintf("%.2f", point.AvgEquipmentAdvantage),
			fmt.Sprintf("%d", point.Team1Wins),
			fmt.Sprintf("%d", point.Team2Wins),
			fmt.Sprintf("%.4f", point.Team1WinRate),
			fmt.Sprintf("%.2f", point.Team1AvgSpent),
			fmt.Sprintf("%.2f", point.Team2AvgSpent),
			fmt.Sprintf("%.2f", point.Team1AvgEarned),
			fmt.Sprintf("%.2f", point.Team2AvgEarned),
			fmt.Sprintf("%.2f", point.Team1AvgSurvivors),
			fmt.Sprintf("%.2f", point.Team2AvgSurvivors),
			fmt.Sprintf("%d", point.GamesReachedThisRound),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// ExportDistributionsCSV exports distribution data
func (ge *GraphExporter) ExportDistributionsCSV() error {
	// Export funds distributions
	if err := ge.exportSingleDistribution("team1_funds_distribution.csv", ge.analysis.Distributions.Team1FundsDistribution); err != nil {
		return err
	}
	if err := ge.exportSingleDistribution("team2_funds_distribution.csv", ge.analysis.Distributions.Team2FundsDistribution); err != nil {
		return err
	}

	// Export equipment distributions
	if err := ge.exportSingleDistribution("team1_equipment_distribution.csv", ge.analysis.Distributions.Team1EquipmentDistribution); err != nil {
		return err
	}
	if err := ge.exportSingleDistribution("team2_equipment_distribution.csv", ge.analysis.Distributions.Team2EquipmentDistribution); err != nil {
		return err
	}

	// Export advantage distributions
	if err := ge.exportSingleDistribution("economic_advantage_distribution.csv", ge.analysis.Distributions.EconomicAdvantageDistribution); err != nil {
		return err
	}
	if err := ge.exportSingleDistribution("equipment_advantage_distribution.csv", ge.analysis.Distributions.EquipmentAdvantageDistribution); err != nil {
		return err
	}

	return nil
}

// exportSingleDistribution exports a single distribution
func (ge *GraphExporter) exportSingleDistribution(filename string, bins []FrequencyBin) error {
	path := filepath.Join(ge.outputPath, "graphs", filename)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"min", "max", "midpoint", "frequency", "percentage"}); err != nil {
		return err
	}

	// Write data
	for _, bin := range bins {
		midpoint := (bin.Min + bin.Max) / 2
		row := []string{
			fmt.Sprintf("%.2f", bin.Min),
			fmt.Sprintf("%.2f", bin.Max),
			fmt.Sprintf("%.2f", midpoint),
			fmt.Sprintf("%d", bin.Frequency),
			fmt.Sprintf("%.2f", bin.Percentage),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// ExportStreakAnalysisCSV exports streak analysis data
func (ge *GraphExporter) ExportStreakAnalysisCSV() error {
	// Export win streaks
	if err := ge.exportStreaks("team1_win_streaks.csv", ge.analysis.Streaks.Team1WinStreaks); err != nil {
		return err
	}
	if err := ge.exportStreaks("team2_win_streaks.csv", ge.analysis.Streaks.Team2WinStreaks); err != nil {
		return err
	}

	// Export loss streaks
	if err := ge.exportStreaks("team1_loss_streaks.csv", ge.analysis.Streaks.Team1LossStreaks); err != nil {
		return err
	}
	if err := ge.exportStreaks("team2_loss_streaks.csv", ge.analysis.Streaks.Team2LossStreaks); err != nil {
		return err
	}

	// Export streak economic impact
	path := filepath.Join(ge.outputPath, "graphs", "streak_economic_impact.csv")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{
		"streak_length",
		"occurrences",
		"avg_economic_advantage_gain",
		"avg_equipment_advantage_gain",
		"next_round_win_rate",
	}); err != nil {
		return err
	}

	for _, impact := range ge.analysis.Streaks.StreakEconomicImpact {
		row := []string{
			fmt.Sprintf("%d", impact.StreakLength),
			fmt.Sprintf("%d", impact.Occurrences),
			fmt.Sprintf("%.2f", impact.AvgEconomicAdvantageGain),
			fmt.Sprintf("%.2f", impact.AvgEquipmentAdvantageGain),
			fmt.Sprintf("%.4f", impact.NextRoundWinRate),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// exportStreaks exports streak data
func (ge *GraphExporter) exportStreaks(filename string, streaks []StreakInfo) error {
	path := filepath.Join(ge.outputPath, "graphs", filename)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{
		"length",
		"start_round",
		"end_round",
		"start_economic_edge",
		"end_economic_edge",
		"start_equipment_edge",
		"end_equipment_edge",
		"economic_change",
		"equipment_change",
		"game_id",
	}); err != nil {
		return err
	}

	for _, streak := range streaks {
		row := []string{
			fmt.Sprintf("%d", streak.Length),
			fmt.Sprintf("%d", streak.StartRound),
			fmt.Sprintf("%d", streak.EndRound),
			fmt.Sprintf("%.2f", streak.StartEconomicEdge),
			fmt.Sprintf("%.2f", streak.EndEconomicEdge),
			fmt.Sprintf("%.2f", streak.StartEquipmentEdge),
			fmt.Sprintf("%.2f", streak.EndEquipmentEdge),
			fmt.Sprintf("%.2f", streak.EconomicChange),
			fmt.Sprintf("%.2f", streak.EquipmentChange),
			streak.GameID,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// ExportWinConditionsCSV exports win condition analysis
func (ge *GraphExporter) ExportWinConditionsCSV() error {
	// Export equipment advantage ranges
	path := filepath.Join(ge.outputPath, "graphs", "equipment_advantage_win_rates.csv")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"min_advantage", "max_advantage", "midpoint", "team1_win_rate", "sample_size"}); err != nil {
		return err
	}

	for _, rng := range ge.analysis.WinConditions.EquipmentAdvantageRanges {
		// Cap extreme values for better readability in CSV
		minVal := rng.MinAdvantage
		maxVal := rng.MaxAdvantage
		if minVal < -100000 {
			minVal = -100000
		}
		if maxVal > 100000 {
			maxVal = 100000
		}
		midpoint := (minVal + maxVal) / 2
		row := []string{
			fmt.Sprintf("%.2f", minVal),
			fmt.Sprintf("%.2f", maxVal),
			fmt.Sprintf("%.2f", midpoint),
			fmt.Sprintf("%.4f", rng.Team1WinRate),
			fmt.Sprintf("%d", rng.SampleSize),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}

	// Export economic advantage ranges
	path = filepath.Join(ge.outputPath, "graphs", "economic_advantage_win_rates.csv")
	file, err = os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer = csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"min_advantage", "max_advantage", "midpoint", "team1_win_rate", "sample_size"}); err != nil {
		return err
	}

	for _, rng := range ge.analysis.WinConditions.EconomicAdvantageRanges {
		// Cap extreme values for better readability in CSV
		minVal := rng.MinAdvantage
		maxVal := rng.MaxAdvantage
		if minVal < -100000 {
			minVal = -100000
		}
		if maxVal > 100000 {
			maxVal = 100000
		}
		midpoint := (minVal + maxVal) / 2
		row := []string{
			fmt.Sprintf("%.2f", minVal),
			fmt.Sprintf("%.2f", maxVal),
			fmt.Sprintf("%.2f", midpoint),
			fmt.Sprintf("%.4f", rng.Team1WinRate),
			fmt.Sprintf("%d", rng.SampleSize),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// ExportSpendingDecisionsCSV exports spending decision scatter data
func (ge *GraphExporter) ExportSpendingDecisionsCSV() error {
	path := filepath.Join(ge.outputPath, "graphs", "spending_decisions.csv")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"available_funds", "amount_spent", "spend_ratio", "round_won", "team"}); err != nil {
		return err
	}

	for _, decision := range ge.analysis.Distributions.SpendingDecisions {
		wonStr := "0"
		if decision.RoundWon {
			wonStr = "1"
		}
		row := []string{
			fmt.Sprintf("%.2f", decision.AvailableFunds),
			fmt.Sprintf("%.2f", decision.AmountSpent),
			fmt.Sprintf("%.4f", decision.SpendRatio),
			wonStr,
			decision.Team,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// ExportWinProbabilityHeatmapCSV exports heatmap data
func (ge *GraphExporter) ExportWinProbabilityHeatmapCSV() error {
	path := filepath.Join(ge.outputPath, "graphs", "win_probability_heatmap.csv")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write([]string{"team1_equipment", "team2_equipment", "team1_win_rate", "sample_size"}); err != nil {
		return err
	}

	for _, row := range ge.analysis.Distributions.WinProbabilityHeatmap {
		for _, cell := range row {
			if cell.SampleSize > 0 {
				dataRow := []string{
					fmt.Sprintf("%.2f", cell.Team1Equipment),
					fmt.Sprintf("%.2f", cell.Team2Equipment),
					fmt.Sprintf("%.4f", cell.Team1WinRate),
					fmt.Sprintf("%d", cell.SampleSize),
				}
				if err := writer.Write(dataRow); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// ExportCompleteAnalysisJSON exports the complete analysis as JSON
func (ge *GraphExporter) ExportCompleteAnalysisJSON() error {
	path := filepath.Join(ge.outputPath, "advanced_analysis.json")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(ge.analysis)
}

// PrintAnalysisSummary prints a text summary of the analysis
func PrintAnalysisSummary(analysis *AdvancedAnalysis) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 ADVANCED ANALYSIS SUMMARY")
	fmt.Println(strings.Repeat("=", 80))

	// Economic Momentum
	fmt.Println("\n💰 Economic Momentum Indicators:")
	fmt.Printf("  • Average Economic Advantage: $%.2f\n", analysis.EconomicMomentum.AverageEconomicAdvantage)
	fmt.Printf("  • Average Equipment Differential: $%.2f\n", analysis.EconomicMomentum.AverageEquipmentDifferential)
	fmt.Printf("  • Economic Volatility (StdDev): $%.2f\n", analysis.EconomicMomentum.EconomicVolatility)
	fmt.Printf("  • Equipment Volatility (StdDev): $%.2f\n", analysis.EconomicMomentum.EquipmentVolatility)
	fmt.Printf("  • Momentum Shifts: %d\n", analysis.EconomicMomentum.MomentumShifts)
	fmt.Printf("  • Average Momentum Duration: %.2f rounds\n", analysis.EconomicMomentum.AverageMomentumDuration)
	fmt.Printf("  • Team 1 Spend Efficiency: %.4f wins per $1K spent\n", analysis.EconomicMomentum.Team1SpendEfficiency)
	fmt.Printf("  • Team 2 Spend Efficiency: %.4f wins per $1K spent\n", analysis.EconomicMomentum.Team2SpendEfficiency)

	// Win Conditions
	fmt.Println("\n🎯 Win Condition Analysis:")
	fmt.Printf("  • Team 1 Equipment ROI: %.4f wins per $1K spent\n", analysis.WinConditions.Team1EquipmentROI)
	fmt.Printf("  • Team 2 Equipment ROI: %.4f wins per $1K spent\n", analysis.WinConditions.Team2EquipmentROI)
	fmt.Printf("  • Equipment-Outcome Correlation: %.4f\n", analysis.WinConditions.EquipmentCorrelation)
	fmt.Printf("  • Funds-Outcome Correlation: %.4f\n", analysis.WinConditions.FundsCorrelation)
	fmt.Printf("  • Consecutive Loss Correlation: %.4f\n", analysis.WinConditions.ConsecLossCorrelation)

	// Comeback Analysis
	fmt.Println("\n🔄 Comeback Analysis:")
	for deficit := 1; deficit <= 5; deficit++ {
		if scenario, exists := analysis.ComebackAnalysis.ComebacksByDeficit[deficit]; exists {
			fmt.Printf("  • %d-Round Deficit: %.1f%% success rate (%d/%d attempts)\n",
				deficit, scenario.SuccessRate*100, scenario.Successes, scenario.Attempts)
			if scenario.Attempts > 0 {
				fmt.Printf("    - Avg Economic Edge during comeback: $%.2f\n", scenario.AvgEconomicEdge)
				fmt.Printf("    - Avg Equipment Edge during comeback: $%.2f\n", scenario.AvgEquipmentEdge)
			}
		}
	}

	// Half/Side Effects
	fmt.Println("\n⚖️  Half & Side Effects:")
	fmt.Printf("  First Half:\n")
	fmt.Printf("    • Team 1 Win Rate: %.2f%%\n", analysis.HalfSideEffects.FirstHalfTeam1WinRate*100)
	fmt.Printf("  Second Half:\n")
	fmt.Printf("    • Team 1 Win Rate: %.2f%%\n", analysis.HalfSideEffects.SecondHalfTeam1WinRate*100)
	fmt.Printf("  Team 1 as CT:\n")
	fmt.Printf("    • Win Rate: %.2f%% (%d/%d)\n", analysis.HalfSideEffects.Team1CTWinRate*100,
		analysis.HalfSideEffects.Team1CTWins, analysis.HalfSideEffects.Team1CTRounds)
	fmt.Printf("    • Avg Funds: $%.2f\n", analysis.HalfSideEffects.Team1CTAvgFunds)
	fmt.Printf("  Team 1 as T:\n")
	fmt.Printf("    • Win Rate: %.2f%% (%d/%d)\n", analysis.HalfSideEffects.Team1TWinRate*100,
		analysis.HalfSideEffects.Team1TWins, analysis.HalfSideEffects.Team1TRounds)
	fmt.Printf("    • Avg Funds: $%.2f\n", analysis.HalfSideEffects.Team1TAvgFunds)

	// Streak Analysis
	fmt.Println("\n🏆 Streak Analysis:")
	fmt.Printf("  Team 1:\n")
	fmt.Printf("    • Max Win Streak: %d rounds\n", analysis.Streaks.Team1MaxWinStreak)
	fmt.Printf("    • Avg Win Streak: %.2f rounds\n", analysis.Streaks.Team1AvgWinStreak)
	fmt.Printf("  Team 2:\n")
	fmt.Printf("    • Max Win Streak: %d rounds\n", analysis.Streaks.Team2MaxWinStreak)
	fmt.Printf("    • Avg Win Streak: %.2f rounds\n", analysis.Streaks.Team2AvgWinStreak)

	fmt.Println("\n" + strings.Repeat("=", 80))
}
