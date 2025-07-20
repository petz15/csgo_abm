# CS:GO Agent-Based Model - Parallel Simulation Guide

## Overview

The CSGO_ABM simulation system now supports massively parallel execution with optimizations for handling 1+ million simulations while managing memory usage and preventing system overload.

## Key Features for Large-Scale Simulations

### Performance Optimizations
- **Direct result returns**: Results are extracted directly from game objects instead of file I/O operations
- **No temporary file creation**: Eliminates disk writes and reads for individual simulations
- **Reduced memory allocation**: Faster garbage collection due to eliminated JSON parsing
- **Improved throughput**: Significantly faster processing rates for large simulation batches

### Memory Management
- **Automatic batch processing**: Simulations are processed in configurable batches to prevent memory buildup
- **Aggressive garbage collection**: Automatically triggers GC between batches for large simulation runs
- **Optimized result handling**: Results are passed directly through memory without disk operations
- **Memory monitoring**: Continuously monitors memory usage and forces GC when limits are exceeded

### Concurrency Control
- **Worker pool pattern**: Uses a controlled number of worker goroutines to prevent system overload
- **Timeout protection**: Each individual simulation has a 5-minute timeout to prevent hanging
- **Panic recovery**: Recovers from individual simulation panics without crashing the entire batch

### Progress Monitoring
- **Real-time progress reporting**: Shows completion percentage, processing rate, and elapsed time
- **Adaptive reporting intervals**: More frequent updates for very large simulation runs
- **Resource usage tracking**: Monitors memory usage, GC runs, and batch processing statistics

## Usage Examples

### Small-Scale Testing (1-1,000 simulations)
```bash
# Single simulation
go run ./cmd

# 100 simulations with default settings
go run ./cmd -n 100

# 1,000 simulations with 8 cores
go run ./cmd -n 1000 -c 8
```

### Medium-Scale Analysis (1,000-100,000 simulations)
```bash
# 10,000 simulations with optimized memory settings
go run ./cmd -n 10000 -c 4 -m 500 -b 200

# 50,000 simulations with conservative settings
go run ./cmd -n 50000 -c 4 -m 1000 -b 500
```

### Large-Scale Research (100,000+ simulations)
```bash
# 100,000 simulations (automatically optimized)
go run ./cmd -n 100000 -c 4 -m 1000 -b 1000

# 500,000 simulations with careful resource management
go run ./cmd -n 500000 -c 2 -m 2000 -b 1000

# 1,000,000 simulations (use conservative settings)
go run ./cmd -n 1000000 -c 2 -m 2000 -b 1000
```

### Ultra-Large-Scale Studies (1,000,000+ simulations)
```bash
# 5,000,000 simulations (recommended for research clusters)
go run ./cmd -n 5000000 -c 2 -m 3000 -b 1000

# 10,000,000 simulations (only on high-end systems)
go run ./cmd -n 10000000 -c 1 -m 4000 -b 1000
```

## Parameter Optimization Guidelines

### Number of Concurrent Workers (-c)
- **1-4 cores**: Best for million+ simulations to avoid memory pressure
- **4-8 cores**: Good for 10K-100K simulations on modern systems
- **8+ cores**: Use only for smaller batches (< 10K simulations)

### Memory Limit (-m) in MB
- **500-1000 MB**: Suitable for up to 100K simulations
- **1000-2000 MB**: Recommended for 100K-1M simulations
- **2000-4000 MB**: Use for 1M+ simulations on systems with sufficient RAM

### Batch Size (-b)
- **100-500**: Good for smaller runs (< 100K simulations)
- **500-1000**: Optimal for large runs (100K-1M simulations)
- **1000**: Maximum recommended batch size for ultra-large runs

## Automatic Optimizations

The system automatically applies optimizations for very large simulation runs:

### For 100,000+ simulations:
- Batch size is capped at 1000 to prevent memory issues
- Memory limit is increased to 2000 MB if not specified
- Progress reporting interval is reduced to 10 seconds
- More frequent garbage collection (every 5 batches)

### For 1,000,000+ simulations:
- Small delays between batches to prevent system overload
- Additional memory monitoring and cleanup
- Conservative timeout settings

## Expected Performance

Performance has been significantly improved with direct result returns. Here are typical ranges:

- **Single-core performance**: 0.2-1.0 simulations/second
- **4-core performance**: 0.8-4.0 simulations/second  
- **8-core performance**: 1.6-8.0 simulations/second

### Estimated Completion Times (Optimized)

| Simulations | 1 Core | 4 Cores | 8 Cores |
|-------------|--------|---------|---------|
| 10,000      | 3-14 hours | 42min-3.5hrs | 21min-1.7hrs |
| 100,000     | 1-6 days | 7-17 hours | 3.5-8.5 hours |
| 1,000,000   | 12-58 days | 3-14 days | 1.5-7 days |

## Memory Requirements

- **Base memory**: ~10-50 MB for the application
- **Per simulation**: ~1-5 KB during processing
- **Batch overhead**: Scales with batch size and concurrent workers
- **Peak usage**: Typically 2-10x the memory limit setting

## Output and Results

### Individual Simulation Results
- Results are now extracted directly from game objects for maximum performance
- No temporary files are created during parallel simulation processing
- Only summary statistics are retained in the final results directory

### Summary Statistics
- Exported to JSON format in timestamped results directory
- Includes win rates, score distributions, overtime rates
- Performance metrics (processing rate, memory usage, GC stats)

### Progress Monitoring
Real-time console output includes:
- Completion percentage and count
- Processing rate (simulations/second)
- Elapsed time
- Memory usage warnings (if applicable)

## Best Practices for Million+ Simulations

1. **Start small**: Test with 1K-10K simulations first to verify settings
2. **Monitor resources**: Watch CPU, memory, and disk usage during runs
3. **Use conservative settings**: Lower concurrency for stability
4. **Plan for time**: Large runs can take days or weeks
5. **Consider batching**: Break very large runs into multiple smaller runs
6. **Backup results**: Summary files contain all important statistics

## Troubleshooting

### Common Issues
- **Memory errors**: Reduce batch size and/or concurrent workers
- **Slow performance**: Check if other processes are using CPU/memory
- **Hanging simulations**: Individual timeouts prevent this, but check system resources

### Performance Tuning
- **Too slow**: Increase concurrent workers (but watch memory)
- **Memory pressure**: Decrease batch size and concurrent workers
- **Disk space**: Results are auto-cleaned, but monitor the results directory

## System Requirements

### Minimum (up to 100K simulations)
- 4 GB RAM
- 2+ CPU cores
- 1 GB free disk space

### Recommended (up to 1M simulations)
- 8 GB RAM
- 4+ CPU cores
- 5 GB free disk space

### High-end (1M+ simulations)
- 16+ GB RAM
- 8+ CPU cores
- 10+ GB free disk space
