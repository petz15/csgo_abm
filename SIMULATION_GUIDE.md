# CS:GO ABM Simulation Guide

This guide covers the unified simulation system with enhanced analysis capabilities.

## Quick Start

### Basic Usage

```bash
# Single simulation with detailed output
go run ./cmd

# Small test run with enhanced reporting
go run ./cmd -n 10 -s

# Large-scale parallel simulation
go run ./cmd -n 1000 -c 4

# Export comprehensive analysis data
go run ./cmd -n 100 -e
```

## Simulation Modes

### Sequential Mode (`-s` flag)
- Runs simulations one after another
- Ideal for debugging and small runs
- Full analysis capabilities
- Consistent, deterministic execution order

```bash
# Sequential simulation with exports
go run ./cmd -n 100 -s -e
```

### Concurrent Mode (default)
- High-performance parallel processing
- Worker pool architecture
- Memory monitoring and optimization
- Scalable to 100,000+ simulations

```bash
# Concurrent simulation with custom worker count
go run ./cmd -n 10000 -c 8 -m 2000
```

## Enhanced Analysis System

### Statistical Features

The unified analysis package provides:

- **Win Rate Analysis**: Team performance with confidence intervals
- **Balance Scoring**: Strategy competitiveness (0-100% scale)
- **Statistical Significance**: Chi-square testing for result reliability
- **Game Categorization**: Close games vs blowouts analysis
- **Distribution Analysis**: Score and round frequency patterns

### Sample Output

```
============================================================
🎮 SIMULATION RESULTS SUMMARY
============================================================
Simulation Mode: Concurrent
Total Simulations: 1,000 (Completed: 1,000, Failed: 0)
Execution Time: 4s

⚡ PERFORMANCE METRICS
----------------------------------------
Processing Rate: 250.0 simulations/second
Peak Memory Usage: 245 MB

📊 GAME ANALYSIS
----------------------------------------
Team 1 Wins: 523 (52.3%)
Team 2 Wins: 477 (47.7%)
Average Rounds per Game: 28.7
Overtime Rate: 23.4%
Strategy Balance Score: 89.2% (Excellent)
Statistical Significance: High confidence (p<0.01)

Game Length: Standard (avg 28.7 rounds)
Close Games (≤3 round diff): 234 (23.4%)
Blowout Games (>10 round diff): 45 (4.5%)
Competitiveness: High

📈 SCORE DISTRIBUTIONS
----------------------------------------
Most Common Scores:
  16-14: 89 games (8.9%)
  16-13: 76 games (7.6%)
  16-12: 67 games (6.7%)
  16-11: 45 games (4.5%)
  16-10: 34 games (3.4%)

🎯 ROUND ANALYSIS
----------------------------------------
Most Common Game Lengths:
  30 rounds: 234 games (23.4%)
  28 rounds: 156 games (15.6%)
  29 rounds: 134 games (13.4%)

🔍 INSIGHTS & RECOMMENDATIONS
----------------------------------------
✅ Well-Balanced Strategies: Both teams show competitive performance
⚡ Moderate Overtime Rate: Games often decided in regulation
📊 Statistically Significant: Results are reliable with this sample size
============================================================
```

## Export Options

### Summary-Only Mode (Default)
- Rich console output with insights
- JSON summary file in results directory
- Recommended for runs >1,000 simulations

### Individual Results Export (`-e` flag)
- Each game exported as separate JSON file
- Complete game state and progression data
- Warning: Can create thousands of files

### Multi-Format Export
When exports are enabled, the system automatically generates:
- `simulation_summary.json`: Complete statistical data
- `summary.csv`: Key metrics for spreadsheet analysis
- `score_distribution.csv`: Score frequency analysis
- `round_distribution.csv`: Game length analysis

## Performance Optimization

### Recommended Settings

| Simulation Count | Workers | Memory Limit | Expected Time |
|------------------|---------|--------------|---------------|
| 1-100           | 2       | 3000MB       | <1 second     |
| 100-1,000       | 4       | 3000MB       | <10 seconds   |
| 1,000-10,000    | 4-8     | 3000MB       | 1-2 minutes   |
| 10,000-100,000  | 4-6     | 3000MB       | 10-30 minutes |
| 100,000+        | 2-4     | 3000MB       | 1+ hours      |

### Memory Management
- Automatic garbage collection when limits exceeded
- Bounded data structures prevent memory leaks
- Progress monitoring with memory usage tracking

### Large-Scale Recommendations
```bash
# For 100K+ simulations - conservative settings
go run ./cmd -n 100000 -c 4 -m 3000

# For 1M+ simulations - minimal workers
go run ./cmd -n 1000000 -c 2 -m 3000

# Avoid exports for very large runs
go run ./cmd -n 500000 -c 3  # No -e flag
```

## Strategy Configuration

### Available Strategies
- `all_in`: Invests all available funds each round
- `default_half`: Invests 50% of available funds

### Custom Strategy Testing
```bash
# Test strategy combinations
go run ./cmd -n 1000 -t1 all_in -t2 default_half

# Compare with reversed assignment
go run ./cmd -n 1000 -t1 default_half -t2 all_in
```

## Analysis Interpretation

### Balance Score Guide
- **90-100%**: Perfectly balanced strategies
- **80-89%**: Well-balanced, competitive play
- **70-79%**: Slight imbalance, acceptable
- **60-69%**: Noticeable imbalance
- **<60%**: Significant strategy advantage

### Statistical Significance
- **High confidence (p<0.01)**: Results are highly reliable
- **Moderate confidence (p<0.05)**: Results are likely reliable
- **Low confidence (p>0.05)**: Larger sample size recommended

### Competitiveness Categories
- **High**: Many close games, few blowouts
- **Moderate**: Balanced mix of game types
- **Low**: Many one-sided games

## Troubleshooting

### Common Issues

**Simulation Timeouts**
- Reduce worker count: `-c 2`
- Increase memory limit: `-m 4000`
- Use sequential mode: `-s`

**Memory Issues**
- Lower memory limit for more frequent GC: `-m 1000`
- Reduce worker count
- Avoid exports for large runs

**Slow Performance**
- Optimize worker count (usually 2-8)
- Ensure adequate memory limit
- Monitor system resources

### Debug Mode
```bash
# Run small sequential simulation for debugging
go run ./cmd -n 5 -s

# Single simulation for detailed inspection
go run ./cmd -n 1
```

## Integration Tips

### Data Analysis Workflow
1. Run simulations with exports: `-e`
2. Import CSV files into analysis tools
3. Use JSON for programmatic analysis
4. Compare balance scores across strategy combinations

### Automated Testing
```bash
# Quick validation runs
go run ./cmd -n 100 -s  # Sequential validation
go run ./cmd -n 100 -c 2  # Concurrent validation

# Performance benchmarking
time go run ./cmd -n 1000 -c 4
```

This guide covers the core functionality of the enhanced simulation system. For detailed technical implementation, see the main README.md.
