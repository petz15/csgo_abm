package main

import (
	"csgo_abm/internal/engine"
	"fmt"
	"os"
	"path/filepath"
)

// CustomConfig holds validated custom configurations
type CustomConfig struct {
	GameRules   engine.GameRules
	ExportPath  string
	IsValidated bool
}

// ValidateAndPrepareCustomizations validates all custom configurations before simulations start
func ValidateAndPrepareCustomizations(gameRulesPath, abmModelsPath, exportPath string) (*CustomConfig, error) {
	config := &CustomConfig{}

	// Load ABM models first (mandatory for simulations to run)
	fmt.Println("🔧 Loading ABM probability models...")
	if err := engine.LoadABMModels(abmModelsPath); err != nil {
		return nil, fmt.Errorf("failed to load ABM models (required): %w", err)
	}
	if abmModelsPath != "" {
		fmt.Printf("✅ ABM probability models loaded successfully from: %s\n", abmModelsPath)
	} else {
		fmt.Println("✅ ABM probability models loaded successfully from default location")
	}

	// Validate and load game rules
	fmt.Println("🔧 Validating game configuration...")

	imported := false
	// Load game rules (handles default fallback internally)
	config.GameRules, imported = engine.NewGameRules(gameRulesPath)
	if imported {
		fmt.Printf("✅ Custom game rules loaded successfully. Custom game rules loaded from: %s\n", gameRulesPath)
	} else {
		fmt.Println("✅ Using default game rules.")
	}

	// Validate export path
	if err := validateExportPath(exportPath); err != nil {
		return nil, fmt.Errorf("export path validation failed: %w", err)
	}
	config.ExportPath = exportPath
	fmt.Printf("✅ Export path validated: %s\n", exportPath)

	// Mark configuration as validated
	config.IsValidated = true
	fmt.Println("✅ All configurations validated successfully!")
	fmt.Println()

	return config, nil
}

// validateExportPath ensures the export path is valid and accessible
func validateExportPath(exportPath string) error {
	// Check if path is absolute or relative
	if !filepath.IsAbs(exportPath) {
		// Convert to absolute path
		absPath, err := filepath.Abs(exportPath)
		if err != nil {
			return fmt.Errorf("could not resolve absolute path for '%s': %w", exportPath, err)
		}
		exportPath = absPath
	}

	// Check if parent directory exists or can be created
	parentDir := filepath.Dir(exportPath)
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		// Try to create parent directories
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return fmt.Errorf("could not create parent directories for '%s': %w", exportPath, err)
		}
	}

	// Test write permissions by creating the target directory
	if err := os.MkdirAll(exportPath, 0755); err != nil {
		return fmt.Errorf("could not create or access export directory '%s': %w", exportPath, err)
	}

	// Test write permissions by creating a temporary file
	testFile := filepath.Join(exportPath, ".write_test")
	if file, err := os.Create(testFile); err != nil {
		return fmt.Errorf("no write permission in export directory '%s': %w", exportPath, err)
	} else {
		file.Close()
		os.Remove(testFile) // Clean up test file
	}

	return nil
}

// ValidateStrategies checks if the specified strategies are available
func ValidateStrategies(team1Strategy, team2Strategy string) error {
	availableStrategies := map[string]bool{
		"all_in":          true,
		"default_half":    true,
		"adaptive_eco_v1": true,
		"yolo":            true,
		"scrooge":         true,
		"adaptive_eco_v2": true,
	}

	if !availableStrategies[team1Strategy] {
		return fmt.Errorf("unknown strategy for team 1: '%s'. Available strategies: all_in, default_half, adaptive_eco_v1, yolo, scrooge", team1Strategy)
	}

	if !availableStrategies[team2Strategy] {
		return fmt.Errorf("unknown strategy for team 2: '%s'. Available strategies: all_in, default_half, adaptive_eco_v1, yolo, scrooge", team2Strategy)
	}

	fmt.Printf("✅ Team strategies validated: %s vs %s\n", team1Strategy, team2Strategy)
	return nil
}
