package analysis

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ExportOptions controls what data to export
type ExportOptions struct {
	ExportIndividualResults bool
	ExportSummary           bool
	ResultsDirectory        string
	SummaryFilename         string
}

// ExportResults exports simulation results based on options
func ExportResults(stats *SimulationStats, options ExportOptions) error {
	if options.ExportSummary {
		summaryPath := options.SummaryFilename
		if options.ResultsDirectory != "" {
			summaryPath = filepath.Join(options.ResultsDirectory, "simulation_summary.json")
		}

		if err := exportSummaryStats(stats, summaryPath); err != nil {
			return fmt.Errorf("failed to export summary: %v", err)
		}
	}

	return nil
}

// exportSummaryStats exports the complete statistics to JSON including configuration
func exportSummaryStats(stats *SimulationStats, filename string) error {
	// Ensure the directory exists
	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Create enhanced export data with metadata
	exportData := struct {
		*SimulationStats
		ExportTimestamp time.Time `json:"export_timestamp"`
		ExportVersion   string    `json:"export_version"`
	}{
		SimulationStats: stats,
		ExportTimestamp: time.Now(),
		ExportVersion:   "2.0",
	}

	data, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// ExportCSVSummary exports a simplified CSV summary for easy analysis
func ExportCSVSummary(stats *SimulationStats, filename string) error {
	if stats.CompletedSims == 0 {
		return fmt.Errorf("no completed simulations to export")
	}

	lines := []string{
		"metric,value",
		fmt.Sprintf("total_simulations,%d", stats.TotalSimulations),
		fmt.Sprintf("completed_simulations,%d", stats.CompletedSims),
		fmt.Sprintf("failed_simulations,%d", stats.FailedSims),
		fmt.Sprintf("team1_wins,%d", stats.Team1Wins),
		fmt.Sprintf("team2_wins,%d", stats.Team2Wins),
		fmt.Sprintf("team1_win_rate,%.2f", stats.Team1WinRate),
		fmt.Sprintf("team2_win_rate,%.2f", stats.Team2WinRate),
		fmt.Sprintf("average_rounds,%.2f", stats.AverageRounds),
		fmt.Sprintf("overtime_rate,%.2f", stats.OvertimeRate),
		fmt.Sprintf("execution_time_seconds,%.2f", stats.ExecutionTime.Seconds()),
		fmt.Sprintf("processing_rate,%.2f", stats.ProcessingRate),
	}

	if stats.AdvancedStats != nil {
		lines = append(lines,
			fmt.Sprintf("balance_score,%.2f", stats.AdvancedStats.BalanceScore),
			fmt.Sprintf("statistical_significance,%s", stats.AdvancedStats.StatisticalSignificance),
			fmt.Sprintf("close_games,%d", stats.AdvancedStats.CloseGames),
			fmt.Sprintf("blowout_games,%d", stats.AdvancedStats.BlowoutGames),
		)

		if stats.AdvancedStats.MedianRounds > 0 {
			lines = append(lines,
				fmt.Sprintf("median_rounds,%.2f", stats.AdvancedStats.MedianRounds),
				fmt.Sprintf("std_dev_rounds,%.2f", stats.AdvancedStats.StdDevRounds),
			)
		}
	}

	// Add configuration information if available
	if stats.Config != nil {
		lines = append(lines, "# Configuration")
		lines = append(lines,
			fmt.Sprintf("team1_name,%s", stats.Config.Team1Name),
			fmt.Sprintf("team1_strategy,%s", stats.Config.Team1Strategy),
			fmt.Sprintf("team2_name,%s", stats.Config.Team2Name),
			fmt.Sprintf("team2_strategy,%s", stats.Config.Team2Strategy),
			fmt.Sprintf("simulation_mode,%s", stats.SimulationMode),
		)

		if !stats.Config.Sequential {
			lines = append(lines,
				fmt.Sprintf("max_concurrent,%d", stats.Config.MaxConcurrent),
				fmt.Sprintf("memory_limit_mb,%d", stats.Config.MemoryLimit),
			)
		}
	}

	content := strings.Join(lines, "\n")
	return os.WriteFile(filename, []byte(content), 0644)
}

// ExportScoreDistribution exports score distribution as CSV
func ExportScoreDistribution(stats *SimulationStats, filename string) error {
	if len(stats.ScoreDistribution) == 0 {
		return fmt.Errorf("no score distribution data to export")
	}

	lines := []string{"score,count,frequency"}

	if stats.AdvancedStats != nil && len(stats.AdvancedStats.TopScoreLines) > 0 {
		for _, scoreLine := range stats.AdvancedStats.TopScoreLines {
			lines = append(lines, fmt.Sprintf("%s,%d,%.2f",
				scoreLine.Score, scoreLine.Count, scoreLine.Frequency))
		}
	} else {
		// Fallback to raw distribution
		for score, count := range stats.ScoreDistribution {
			frequency := float64(count) / float64(stats.CompletedSims) * 100
			lines = append(lines, fmt.Sprintf("%s,%d,%.2f", score, count, frequency))
		}
	}

	content := strings.Join(lines, "\n")
	return os.WriteFile(filename, []byte(content), 0644)
}

// ExportRoundDistribution exports round distribution as CSV
func ExportRoundDistribution(stats *SimulationStats, filename string) error {
	if len(stats.RoundDistribution) == 0 {
		return fmt.Errorf("no round distribution data to export")
	}

	lines := []string{"rounds,count,frequency"}

	for rounds, count := range stats.RoundDistribution {
		frequency := float64(count) / float64(stats.CompletedSims) * 100
		lines = append(lines, fmt.Sprintf("%d,%d,%.2f", rounds, count, frequency))
	}

	content := strings.Join(lines, "\n")
	return os.WriteFile(filename, []byte(content), 0644)
}

// CreateResultsDirectory creates a timestamped results directory
func CreateResultsDirectory() (string, error) {
	now := time.Now()
	dirName := fmt.Sprintf("results_%s", now.Format("20060102_150405"))
	err := os.MkdirAll(dirName, 0755)
	return dirName, err
}

// ExportAllFormats exports results in multiple formats for comprehensive analysis
func ExportAllFormats(stats *SimulationStats, baseDir string) error {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return err
	}

	// Export JSON summary
	if err := exportSummaryStats(stats, filepath.Join(baseDir, "simulation_summary.json")); err != nil {
		return err
	}

	// Export CSV summary
	if err := ExportCSVSummary(stats, filepath.Join(baseDir, "summary.csv")); err != nil {
		return err
	}

	// Export score distribution
	if len(stats.ScoreDistribution) > 0 {
		if err := ExportScoreDistribution(stats, filepath.Join(baseDir, "score_distribution.csv")); err != nil {
			return err
		}
	}

	// Export round distribution
	if len(stats.RoundDistribution) > 0 {
		if err := ExportRoundDistribution(stats, filepath.Join(baseDir, "round_distribution.csv")); err != nil {
			return err
		}
	}

	fmt.Printf("📁 Exported comprehensive results to: %s/\n", baseDir)
	fmt.Println("   - simulation_summary.json (complete data)")
	fmt.Println("   - summary.csv (key metrics)")
	if len(stats.ScoreDistribution) > 0 {
		fmt.Println("   - score_distribution.csv")
	}
	if len(stats.RoundDistribution) > 0 {
		fmt.Println("   - round_distribution.csv")
	}

	return nil
}
