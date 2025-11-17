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

	// Create Python plotting script
	if err := ge.CreatePythonPlottingScript(); err != nil {
		return fmt.Errorf("failed to create plotting script: %w", err)
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

// CreatePythonPlottingScript creates a Python script for visualizations
func (ge *GraphExporter) CreatePythonPlottingScript() error {
	path := filepath.Join(ge.outputPath, "plot_analysis.py")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	script := `#!/usr/bin/env python3
"""
CS:GO Simulation Advanced Analysis Plotting Script
Auto-generated script to visualize simulation data
"""

import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import numpy as np
import json
from pathlib import Path
import os

# Set style
sns.set_style("darkgrid")
plt.rcParams['figure.figsize'] = (12, 6)

# Get script directory and set up paths
script_dir = Path(__file__).parent.resolve()
graphs_dir = script_dir / "graphs"
output_dir = graphs_dir / "plots"
output_dir.mkdir(parents=True, exist_ok=True)

print("📊 CS:GO Simulation Analysis - Generating Visualizations...")
print(f"📁 Working directory: {script_dir}")
print(f"📁 Graphs directory: {graphs_dir}")
print(f"📁 Output directory: {output_dir}")

# ============================================================================
# TIME SERIES PLOTS
# ============================================================================
print("\n1️⃣  Generating Time Series Plots...")

df_ts = pd.read_csv(graphs_dir / "time_series.csv")

# Plot 1: Economic Timeline (Dual-axis)
fig, ax1 = plt.subplots(figsize=(14, 7))
ax1.plot(df_ts['round_number'], df_ts['team1_avg_funds'], 'b-', label='Team 1 Funds', linewidth=2)
ax1.plot(df_ts['round_number'], df_ts['team2_avg_funds'], 'r-', label='Team 2 Funds', linewidth=2)
ax1.set_xlabel('Round Number', fontsize=12)
ax1.set_ylabel('Average Funds ($)', fontsize=12)
ax1.tick_params(axis='y')
ax1.legend(loc='upper left')
ax1.grid(True, alpha=0.3)

ax2 = ax1.twinx()
ax2.plot(df_ts['round_number'], df_ts['team1_avg_equipment'], 'b--', label='Team 1 Equipment', linewidth=1.5, alpha=0.7)
ax2.plot(df_ts['round_number'], df_ts['team2_avg_equipment'], 'r--', label='Team 2 Equipment', linewidth=1.5, alpha=0.7)
ax2.set_ylabel('Average Equipment Value ($)', fontsize=12)
ax2.legend(loc='upper right')

plt.title('Economic Timeline: Funds & Equipment by Round', fontsize=14, fontweight='bold')
plt.tight_layout()
plt.savefig(output_dir / 'economic_timeline.png', dpi=300, bbox_inches='tight')
plt.close()

# Plot 2: Equipment Differential Over Time
plt.figure(figsize=(14, 6))
plt.plot(df_ts['round_number'], df_ts['avg_equipment_advantage'], 'g-', linewidth=2)
plt.axhline(y=0, color='black', linestyle='--', alpha=0.5)
plt.fill_between(df_ts['round_number'], 0, df_ts['avg_equipment_advantage'], 
                 where=(df_ts['avg_equipment_advantage'] > 0), alpha=0.3, color='blue', label='Team 1 Advantage')
plt.fill_between(df_ts['round_number'], 0, df_ts['avg_equipment_advantage'], 
                 where=(df_ts['avg_equipment_advantage'] <= 0), alpha=0.3, color='red', label='Team 2 Advantage')
plt.xlabel('Round Number', fontsize=12)
plt.ylabel('Equipment Advantage ($)', fontsize=12)
plt.title('Equipment Differential Over Time', fontsize=14, fontweight='bold')
plt.legend()
plt.grid(True, alpha=0.3)
plt.tight_layout()
plt.savefig(output_dir / 'equipment_differential_timeline.png', dpi=300, bbox_inches='tight')
plt.close()

# Plot 3: Economic Momentum
plt.figure(figsize=(14, 6))
plt.plot(df_ts['round_number'], df_ts['avg_economic_advantage'], 'purple', linewidth=2)
plt.axhline(y=0, color='black', linestyle='--', alpha=0.5)
plt.fill_between(df_ts['round_number'], 0, df_ts['avg_economic_advantage'], 
                 where=(df_ts['avg_economic_advantage'] > 0), alpha=0.3, color='blue')
plt.fill_between(df_ts['round_number'], 0, df_ts['avg_economic_advantage'], 
                 where=(df_ts['avg_economic_advantage'] <= 0), alpha=0.3, color='red')
plt.xlabel('Round Number', fontsize=12)
plt.ylabel('Economic Advantage ($)', fontsize=12)
plt.title('Economic Momentum: Fund Differential by Round', fontsize=14, fontweight='bold')
plt.grid(True, alpha=0.3)
plt.tight_layout()
plt.savefig(output_dir / 'economic_momentum.png', dpi=300, bbox_inches='tight')
plt.close()

# Plot 4: Win Rate by Round
plt.figure(figsize=(14, 6))
plt.plot(df_ts['round_number'], df_ts['team1_win_rate'], 'b-', linewidth=2, marker='o', markersize=4)
plt.axhline(y=0.5, color='gray', linestyle='--', alpha=0.5, label='50% (Balanced)')
plt.xlabel('Round Number', fontsize=12)
plt.ylabel('Team 1 Win Rate', fontsize=12)
plt.title('Team 1 Win Rate by Round Number', fontsize=14, fontweight='bold')
plt.ylim(0, 1)
plt.legend()
plt.grid(True, alpha=0.3)
plt.tight_layout()
plt.savefig(output_dir / 'win_rate_by_round.png', dpi=300, bbox_inches='tight')
plt.close()

print("✅ Time series plots saved")

# ============================================================================
# DISTRIBUTION PLOTS
# ============================================================================
print("\n2️⃣  Generating Distribution Plots...")

# Plot 5: Funds Distributions
df_t1_funds = pd.read_csv(graphs_dir / "team1_funds_distribution.csv")
df_t2_funds = pd.read_csv(graphs_dir / "team2_funds_distribution.csv")

fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 6))
ax1.bar(df_t1_funds['midpoint'], df_t1_funds['percentage'], width=df_t1_funds['max']-df_t1_funds['min'], 
        alpha=0.7, color='blue', edgecolor='black')
ax1.set_xlabel('Funds ($)', fontsize=12)
ax1.set_ylabel('Percentage (%)', fontsize=12)
ax1.set_title('Team 1 Funds Distribution', fontsize=12, fontweight='bold')
ax1.grid(True, alpha=0.3)

ax2.bar(df_t2_funds['midpoint'], df_t2_funds['percentage'], width=df_t2_funds['max']-df_t2_funds['min'], 
        alpha=0.7, color='red', edgecolor='black')
ax2.set_xlabel('Funds ($)', fontsize=12)
ax2.set_ylabel('Percentage (%)', fontsize=12)
ax2.set_title('Team 2 Funds Distribution', fontsize=12, fontweight='bold')
ax2.grid(True, alpha=0.3)

plt.tight_layout()
plt.savefig(output_dir / 'funds_distributions.png', dpi=300, bbox_inches='tight')
plt.close()

# Plot 6: Equipment Distributions
df_t1_equip = pd.read_csv(graphs_dir / "team1_equipment_distribution.csv")
df_t2_equip = pd.read_csv(graphs_dir / "team2_equipment_distribution.csv")

fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 6))
ax1.bar(df_t1_equip['midpoint'], df_t1_equip['percentage'], width=df_t1_equip['max']-df_t1_equip['min'], 
        alpha=0.7, color='blue', edgecolor='black')
ax1.set_xlabel('Equipment Value ($)', fontsize=12)
ax1.set_ylabel('Percentage (%)', fontsize=12)
ax1.set_title('Team 1 Equipment Distribution', fontsize=12, fontweight='bold')
ax1.grid(True, alpha=0.3)

ax2.bar(df_t2_equip['midpoint'], df_t2_equip['percentage'], width=df_t2_equip['max']-df_t2_equip['min'], 
        alpha=0.7, color='red', edgecolor='black')
ax2.set_xlabel('Equipment Value ($)', fontsize=12)
ax2.set_ylabel('Percentage (%)', fontsize=12)
ax2.set_title('Team 2 Equipment Distribution', fontsize=12, fontweight='bold')
ax2.grid(True, alpha=0.3)

plt.tight_layout()
plt.savefig(output_dir / 'equipment_distributions.png', dpi=300, bbox_inches='tight')
plt.close()

# Plot 7: Economic Advantage Distribution
df_econ_adv = pd.read_csv(graphs_dir / "economic_advantage_distribution.csv")

plt.figure(figsize=(12, 6))
colors = ['red' if x < 0 else 'blue' for x in df_econ_adv['midpoint']]
plt.bar(df_econ_adv['midpoint'], df_econ_adv['percentage'], 
        width=df_econ_adv['max']-df_econ_adv['min'], alpha=0.7, color=colors, edgecolor='black')
plt.axvline(x=0, color='black', linestyle='--', linewidth=2)
plt.xlabel('Economic Advantage (Team1 - Team2) ($)', fontsize=12)
plt.ylabel('Percentage (%)', fontsize=12)
plt.title('Economic Advantage Distribution', fontsize=14, fontweight='bold')
plt.grid(True, alpha=0.3)
plt.tight_layout()
plt.savefig(output_dir / 'economic_advantage_distribution.png', dpi=300, bbox_inches='tight')
plt.close()

# Plot 8: Equipment Advantage Distribution
df_equip_adv = pd.read_csv(graphs_dir / "equipment_advantage_distribution.csv")

plt.figure(figsize=(12, 6))
colors = ['red' if x < 0 else 'blue' for x in df_equip_adv['midpoint']]
plt.bar(df_equip_adv['midpoint'], df_equip_adv['percentage'], 
        width=df_equip_adv['max']-df_equip_adv['min'], alpha=0.7, color=colors, edgecolor='black')
plt.axvline(x=0, color='black', linestyle='--', linewidth=2)
plt.xlabel('Equipment Advantage (Team1 - Team2) ($)', fontsize=12)
plt.ylabel('Percentage (%)', fontsize=12)
plt.title('Equipment Advantage Distribution', fontsize=14, fontweight='bold')
plt.grid(True, alpha=0.3)
plt.tight_layout()
plt.savefig(output_dir / 'equipment_advantage_distribution.png', dpi=300, bbox_inches='tight')
plt.close()

print("✅ Distribution plots saved")

# ============================================================================
# WIN CONDITION ANALYSIS
# ============================================================================
print("\n3️⃣  Generating Win Condition Plots...")

# Plot 9: Equipment Advantage vs Win Rate
df_equip_wins = pd.read_csv(graphs_dir / "equipment_advantage_win_rates.csv")

plt.figure(figsize=(12, 6))
plt.plot(df_equip_wins['midpoint'], df_equip_wins['team1_win_rate'], 'o-', linewidth=2, markersize=8, color='green')
plt.axhline(y=0.5, color='gray', linestyle='--', alpha=0.5, label='50% (Balanced)')
plt.axvline(x=0, color='black', linestyle='--', alpha=0.5)
plt.xlabel('Equipment Advantage ($)', fontsize=12)
plt.ylabel('Team 1 Win Rate', fontsize=12)
plt.title('Win Rate by Equipment Advantage (Continuous)', fontsize=14, fontweight='bold')
plt.ylim(0, 1)
plt.legend()
plt.grid(True, alpha=0.3)
plt.tight_layout()
plt.savefig(output_dir / 'equipment_advantage_win_rate.png', dpi=300, bbox_inches='tight')
plt.close()

# Plot 10: Economic Advantage vs Win Rate
df_econ_wins = pd.read_csv(graphs_dir / "economic_advantage_win_rates.csv")

plt.figure(figsize=(12, 6))
plt.plot(df_econ_wins['midpoint'], df_econ_wins['team1_win_rate'], 'o-', linewidth=2, markersize=8, color='purple')
plt.axhline(y=0.5, color='gray', linestyle='--', alpha=0.5, label='50% (Balanced)')
plt.axvline(x=0, color='black', linestyle='--', alpha=0.5)
plt.xlabel('Economic Advantage ($)', fontsize=12)
plt.ylabel('Team 1 Win Rate', fontsize=12)
plt.title('Win Rate by Economic Advantage (Continuous)', fontsize=14, fontweight='bold')
plt.ylim(0, 1)
plt.legend()
plt.grid(True, alpha=0.3)
plt.tight_layout()
plt.savefig(output_dir / 'economic_advantage_win_rate.png', dpi=300, bbox_inches='tight')
plt.close()

# Plot 11: Spending Decision Scatter
df_spending = pd.read_csv(graphs_dir / "spending_decisions.csv")
df_spending_sample = df_spending.sample(min(5000, len(df_spending)))  # Sample for performance

fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 6))

# Team 1
df_t1 = df_spending_sample[df_spending_sample['team'] == 'Team1']
colors_t1 = ['green' if x == 1 else 'red' for x in df_t1['round_won']]
ax1.scatter(df_t1['available_funds'], df_t1['amount_spent'], c=colors_t1, alpha=0.3, s=10)
ax1.set_xlabel('Available Funds ($)', fontsize=12)
ax1.set_ylabel('Amount Spent ($)', fontsize=12)
ax1.set_title('Team 1 Spending Decisions', fontsize=12, fontweight='bold')
ax1.grid(True, alpha=0.3)

# Team 2
df_t2 = df_spending_sample[df_spending_sample['team'] == 'Team2']
colors_t2 = ['green' if x == 1 else 'red' for x in df_t2['round_won']]
ax2.scatter(df_t2['available_funds'], df_t2['amount_spent'], c=colors_t2, alpha=0.3, s=10)
ax2.set_xlabel('Available Funds ($)', fontsize=12)
ax2.set_ylabel('Amount Spent ($)', fontsize=12)
ax2.set_title('Team 2 Spending Decisions', fontsize=12, fontweight='bold')
ax2.grid(True, alpha=0.3)

plt.tight_layout()
plt.savefig(output_dir / 'spending_decisions.png', dpi=300, bbox_inches='tight')
plt.close()

print("✅ Win condition plots saved")

# ============================================================================
# STREAK ANALYSIS
# ============================================================================
print("\n4️⃣  Generating Streak Analysis Plots...")

# Plot 12: Streak Length Distributions
df_t1_streaks = pd.read_csv(graphs_dir / "team1_win_streaks.csv")
df_t2_streaks = pd.read_csv(graphs_dir / "team2_win_streaks.csv")

fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 6))

ax1.hist(df_t1_streaks['length'], bins=range(1, max(df_t1_streaks['length'])+2), 
         alpha=0.7, color='blue', edgecolor='black')
ax1.set_xlabel('Streak Length (Rounds)', fontsize=12)
ax1.set_ylabel('Frequency', fontsize=12)
ax1.set_title('Team 1 Win Streak Distribution', fontsize=12, fontweight='bold')
ax1.grid(True, alpha=0.3)

ax2.hist(df_t2_streaks['length'], bins=range(1, max(df_t2_streaks['length'])+2), 
         alpha=0.7, color='red', edgecolor='black')
ax2.set_xlabel('Streak Length (Rounds)', fontsize=12)
ax2.set_ylabel('Frequency', fontsize=12)
ax2.set_title('Team 2 Win Streak Distribution', fontsize=12, fontweight='bold')
ax2.grid(True, alpha=0.3)

plt.tight_layout()
plt.savefig(output_dir / 'win_streak_distributions.png', dpi=300, bbox_inches='tight')
plt.close()

# Plot 13: Streak Economic Impact
df_streak_impact = pd.read_csv(graphs_dir / "streak_economic_impact.csv")

fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 6))

ax1.bar(df_streak_impact['streak_length'], df_streak_impact['avg_economic_advantage_gain'], 
        alpha=0.7, color='purple', edgecolor='black')
ax1.set_xlabel('Streak Length', fontsize=12)
ax1.set_ylabel('Avg Economic Advantage Gain ($)', fontsize=12)
ax1.set_title('Economic Impact of Win Streaks', fontsize=12, fontweight='bold')
ax1.grid(True, alpha=0.3)

ax2.bar(df_streak_impact['streak_length'], df_streak_impact['avg_equipment_advantage_gain'], 
        alpha=0.7, color='orange', edgecolor='black')
ax2.set_xlabel('Streak Length', fontsize=12)
ax2.set_ylabel('Avg Equipment Advantage Gain ($)', fontsize=12)
ax2.set_title('Equipment Impact of Win Streaks', fontsize=12, fontweight='bold')
ax2.grid(True, alpha=0.3)

plt.tight_layout()
plt.savefig(output_dir / 'streak_economic_impact.png', dpi=300, bbox_inches='tight')
plt.close()

print("✅ Streak analysis plots saved")

# ============================================================================
# SUMMARY STATISTICS
# ============================================================================
print("\n5️⃣  Loading Advanced Analysis Summary...")

with open(script_dir / "advanced_analysis.json", "r") as f:
    analysis = json.load(f)

print("\n" + "="*60)
print("📈 ADVANCED ANALYSIS SUMMARY")
print("="*60)

print("\n💰 Economic Momentum:")
print(f"  Average Economic Advantage: ${analysis['economic_momentum']['average_economic_advantage']:.2f}")
print(f"  Average Equipment Differential: ${analysis['economic_momentum']['average_equipment_differential']:.2f}")
print(f"  Economic Volatility: ${analysis['economic_momentum']['economic_volatility']:.2f}")
print(f"  Momentum Shifts: {analysis['economic_momentum']['momentum_shifts']}")
print(f"  Team 1 Spend Efficiency: {analysis['economic_momentum']['team1_spend_efficiency']:.4f} wins/$1K")
print(f"  Team 2 Spend Efficiency: {analysis['economic_momentum']['team2_spend_efficiency']:.4f} wins/$1K")

print("\n🎯 Win Conditions:")
print(f"  Team 1 Equipment ROI: {analysis['win_conditions']['team1_equipment_roi']:.4f} wins/$1K")
print(f"  Team 2 Equipment ROI: {analysis['win_conditions']['team2_equipment_roi']:.4f} wins/$1K")
print(f"  Equipment Correlation: {analysis['win_conditions']['equipment_correlation']:.4f}")
print(f"  Funds Correlation: {analysis['win_conditions']['funds_correlation']:.4f}")

print("\n🔄 Comeback Analysis:")
for deficit, scenario in analysis['comeback_analysis']['comebacks_by_deficit'].items():
    print(f"  Deficit {deficit}: {scenario['success_rate']*100:.1f}% success rate ({scenario['successes']}/{scenario['attempts']})")

print("\n⚖️  Half/Side Effects:")
print(f"  First Half Team 1 Win Rate: {analysis['half_side_effects']['first_half_team1_win_rate']*100:.1f}%")
print(f"  Second Half Team 1 Win Rate: {analysis['half_side_effects']['second_half_team1_win_rate']*100:.1f}%")
print(f"  Team 1 CT Win Rate: {analysis['half_side_effects']['team1_ct_win_rate']*100:.1f}%")
print(f"  Team 1 T Win Rate: {analysis['half_side_effects']['team1_t_win_rate']*100:.1f}%")

print("\n🏆 Streaks:")
print(f"  Team 1 Max Win Streak: {analysis['streaks']['team1_max_win_streak']} rounds")
print(f"  Team 1 Avg Win Streak: {analysis['streaks']['team1_avg_win_streak']:.2f} rounds")
print(f"  Team 2 Max Win Streak: {analysis['streaks']['team2_max_win_streak']} rounds")
print(f"  Team 2 Avg Win Streak: {analysis['streaks']['team2_avg_win_streak']:.2f} rounds")

print("\n" + "="*60)
print(f"✅ All visualizations saved to: {output_dir}")
print("="*60)
`

	_, err = file.WriteString(script)
	return err
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
