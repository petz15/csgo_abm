package analysis

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

// CreateResultsDirectory creates a timestamped results directory
func CreateResultsDirectory() (string, error) {
	now := time.Now()
	dirName := fmt.Sprintf("results_%s", now.Format("20060102_150405"))
	err := os.MkdirAll(dirName, 0755)
	return dirName, err
}

// CreateResultsDirectoryAt creates a timestamped results directory under a base path
func CreateResultsDirectoryAt(base string) (string, error) {
	if base == "" {
		return CreateResultsDirectory()
	}
	now := time.Now()
	dirName := fmt.Sprintf("results_%s", now.Format("20060102_150405"))
	full := filepath.Join(base, dirName)
	if err := os.MkdirAll(full, 0755); err != nil {
		return "", err
	}
	return full, nil
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

	return nil
}
