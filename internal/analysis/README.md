# Analysis Package Documentation

The `internal/analysis` package provides unified statistical analysis for CS:GO economic simulations.

## Core Components

### SimulationStats (`stats.go`)
Unified statistics structure for both sequential and concurrent simulation modes.

```go
type SimulationStats struct {
    // Core metrics
    TotalSimulations  int64
    CompletedSims     int64
    FailedSims        int64
    Team1Wins         int64
    Team2Wins         int64
    
    // Advanced analysis
    AdvancedStats *AdvancedStats
    
    // Thread safety
    ScoreMutex sync.Mutex
    RoundMutex sync.Mutex
}
```

### Statistical Calculator (`calculator.go`)
Thread-safe computation engine with atomic operations.

**Key Methods:**
- `UpdateGameResult()`: Thread-safe game result processing
- `CalculateFinalStats()`: Comprehensive statistical analysis
- `UpdateFailedSimulation()`: Error tracking

### Enhanced Reporter (`reporter.go`)
Rich visual output with strategic insights.

**Features:**
- Balance scoring (0-100% scale)
- Statistical significance testing
- Game categorization (close vs blowouts)
- Strategic recommendations

### Multi-Format Exporter (`exporter.go`)
Comprehensive data export capabilities.

**Export Formats:**
- JSON: Complete statistical data
- CSV: Key metrics for analysis
- Distribution files: Score and round frequency

## Usage Examples

### Basic Integration
```go
// Initialize statistics
stats := analysis.NewStats(1000, "concurrent")

// Update with game results
stats.UpdateGameResult(true, 16, 12, 28, false, 0)

// Calculate final analysis
stats.CalculateFinalStats()

// Generate enhanced report
analysis.PrintEnhancedStats(stats)

// Export results
analysis.ExportAllFormats(stats, "results_dir")
```

### Configuration
```go
config := analysis.SimulationConfig{
    NumSimulations: 1000,
    MaxConcurrent:  4,
    Team1Strategy:  "all_in",
    Team2Strategy:  "default_half",
    ExportResults:  true,
    Sequential:     false,
}
```

## Thread Safety

All statistical operations are thread-safe:
- Atomic operations for core counters
- Mutex protection for map operations
- Safe for concurrent access across goroutines

## Advanced Features

### Balance Scoring
Calculates strategy competitiveness on 0-100% scale:
- 90-100%: Perfectly balanced
- 80-89%: Well-balanced
- <60%: Significant imbalance

### Statistical Significance
Chi-square testing for result reliability:
- p<0.01: High confidence
- p<0.05: Moderate confidence
- p>0.05: Larger sample recommended

### Game Categorization
- Close games: ≤3 round difference
- Blowouts: >10 round difference
- Overtime analysis and patterns

## Performance

- Minimal overhead: <1% performance impact
- Memory efficient: Bounded data structures
- Scalable: Tested with 1M+ simulations
- Thread-safe: Full concurrent support

## Integration Points

The analysis package integrates with:
- `simulation_concurrent.go`: Parallel processing
- `simulation_sequential.go`: Sequential execution
- `main.go`: Unified configuration
- Export systems: Multiple output formats
