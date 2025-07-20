package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// SimulationConfig holds configuration for parallel simulations
type SimulationConfig struct {
	NumSimulations int
	MaxConcurrent  int
	MemoryLimit    int // Memory limit in MB before forcing GC
	Team1Name      string
	Team1Strategy  string
	Team2Name      string
	Team2Strategy  string
	GameRules      string
	ExportResults  bool // Whether to export individual game results
}

// SimulationStats tracks statistics across all simulations
type SimulationStats struct {
	TotalSimulations  int64            `json:"total_simulations"`
	CompletedSims     int64            `json:"completed_simulations"`
	Team1Wins         int64            `json:"team1_wins"`
	Team2Wins         int64            `json:"team2_wins"`
	TotalRounds       int64            `json:"total_rounds"`
	OvertimeGames     int64            `json:"overtime_games"`
	AverageRounds     float64          `json:"average_rounds"`
	Team1WinRate      float64          `json:"team1_win_rate"`
	Team2WinRate      float64          `json:"team2_win_rate"`
	OvertimeRate      float64          `json:"overtime_rate"`
	ExecutionTime     time.Duration    `json:"execution_time"`
	ScoreDistribution map[string]int64 `json:"score_distribution"`
	ProcessingRate    float64          `json:"simulations_per_second"`
	PeakMemoryUsage   uint64           `json:"peak_memory_usage_mb"`
	TotalGCRuns       uint32           `json:"total_gc_runs"`
	scoreMutex        sync.Mutex       // Protects ScoreDistribution map
}

// SimulationResult holds the result of a single simulation
type SimulationResult struct {
	GameID         string
	Team1Won       bool
	Team1Score     int
	Team2Score     int
	TotalRounds    int
	WentToOvertime bool
	Error          error
}

// WorkerPool manages a pool of workers for running simulations
type WorkerPool struct {
	workers int
	jobs    chan SimulationJob
	results chan SimulationResult
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	stats   *SimulationStats
}

// SimulationJob represents a single simulation job
type SimulationJob struct {
	SimID  int
	Config SimulationConfig
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, stats *SimulationStats) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers: workers,
		jobs:    make(chan SimulationJob, workers*2), // Buffer jobs
		results: make(chan SimulationResult, workers*2),
		ctx:     ctx,
		cancel:  cancel,
		stats:   stats,
	}
}

// Start initializes and starts the worker pool
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	close(wp.jobs)    // Signal workers to stop accepting jobs
	wp.cancel()       // Cancel context to stop any blocked operations
	wp.wg.Wait()      // Wait for all workers to finish
	close(wp.results) // Close results channel after workers are done
}

// AddJob adds a simulation job to the queue, returns false if pool is shutting down
func (wp *WorkerPool) AddJob(job SimulationJob) bool {
	select {
	case wp.jobs <- job:
		return true
	case <-wp.ctx.Done():
		return false
	}
}

// worker processes simulation jobs
func (wp *WorkerPool) worker(workerID int) {
	defer wp.wg.Done()

	for {
		select {
		case job, ok := <-wp.jobs:
			if !ok {
				return // Channel closed, exit worker
			}
			result := wp.processSingleSimulation(job)

			select {
			case wp.results <- result:
			case <-wp.ctx.Done():
				return
			}

		case <-wp.ctx.Done():
			return
		}
	}
}

// processSingleSimulation runs a single simulation and returns the result
func (wp *WorkerPool) processSingleSimulation(job SimulationJob) SimulationResult {
	simPrefix := fmt.Sprintf("sim_%d_", job.SimID)

	// Use worker pool context with timeout for proper cancellation
	ctx, cancel := context.WithTimeout(wp.ctx, 5*time.Minute)
	defer cancel()

	resultChan := make(chan SimulationResult, 1)

	// Run the simulation in a goroutine with proper context handling
	go func() {
		defer func() {
			if r := recover(); r != nil {
				select {
				case resultChan <- SimulationResult{
					GameID: fmt.Sprintf("%s_panic", simPrefix),
					Error:  fmt.Errorf("simulation panicked: %v", r),
				}:
				case <-ctx.Done():
					// Context cancelled, don't block
				}
			}
		}()

		// Run the simulation and get results directly
		gameResult, err := StartGameWithResults(
			job.Config.Team1Name,
			job.Config.Team1Strategy,
			job.Config.Team2Name,
			job.Config.Team2Strategy,
			job.Config.GameRules,
			simPrefix,
		)

		var result SimulationResult
		if err != nil {
			result = SimulationResult{
				GameID: fmt.Sprintf("%s_error", simPrefix),
				Error:  err,
			}
		} else {
			result = SimulationResult{
				GameID:         gameResult.GameID,
				Team1Won:       gameResult.Team1Won,
				Team1Score:     gameResult.Team1Score,
				Team2Score:     gameResult.Team2Score,
				TotalRounds:    gameResult.TotalRounds,
				WentToOvertime: gameResult.WentToOvertime,
				Error:          nil,
			}
		}

		// Ensure we don't block on send if context is cancelled
		select {
		case resultChan <- result:
		case <-ctx.Done():
			// Context cancelled, don't block
		}
	}()

	// Wait for either completion or timeout
	select {
	case result := <-resultChan:
		return result
	case <-ctx.Done():
		return SimulationResult{
			GameID: fmt.Sprintf("%s_timeout", simPrefix),
			Error:  fmt.Errorf("simulation timed out after 5 minutes"),
		}
	}
} // RunParallelSimulations orchestrates the execution of multiple simulations
func RunParallelSimulations(config SimulationConfig) error {
	startTime := time.Now()

	// Optimize settings for very large simulations
	if config.NumSimulations >= 100000 {
		// For very large simulations, reduce memory pressure
		if config.MemoryLimit > 2000 {
			config.MemoryLimit = 2000
		}
	}

	// Initialize statistics
	stats := &SimulationStats{
		TotalSimulations:  int64(config.NumSimulations),
		ScoreDistribution: make(map[string]int64),
	}

	// Create results directory
	resultsDir := fmt.Sprintf("results_%s", time.Now().Format("20060102_150405"))
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return fmt.Errorf("failed to create results directory: %v", err)
	}

	fmt.Printf("Starting %d simulations with %d concurrent workers...\n",
		config.NumSimulations, config.MaxConcurrent)
	fmt.Printf("Memory limit: %d MB\n", config.MemoryLimit)

	if config.ExportResults {
		fmt.Printf("Individual result export: ENABLED (results will be saved to %s/)\n", resultsDir)
		if config.NumSimulations > 10000 {
			fmt.Printf("WARNING: Exporting %d individual results may create filesystem pressure\n", config.NumSimulations)
		}
	} else {
		fmt.Println("Individual result export: DISABLED (summary-only mode)")
	}

	// Create worker pool
	pool := NewWorkerPool(config.MaxConcurrent, stats)
	pool.Start()

	// Create shutdown context for monitoring goroutines
	monitorCtx, monitorCancel := context.WithCancel(context.Background())
	defer monitorCancel()

	// Start result collector goroutine
	resultsDone := make(chan bool)
	go func() {
		defer close(resultsDone)
		collectResults(pool.results, stats, config.NumSimulations, config.ExportResults, resultsDir)
	}()

	// Track memory usage
	memoryMonitorDone := make(chan bool)
	go func() {
		defer close(memoryMonitorDone)
		monitorMemoryUsageWithContext(monitorCtx, stats, config.MemoryLimit)
	}()

	// Progress reporter
	progressDone := make(chan bool)
	progressInterval := 30 * time.Second
	// For very large simulations, report progress more frequently
	if config.NumSimulations >= 100000 {
		progressInterval = 10 * time.Second
	}
	go func() {
		defer close(progressDone)
		reportProgressWithContext(monitorCtx, stats, config.NumSimulations, startTime, progressInterval)
	}()

	// Submit all jobs to the worker pool
	jobsSubmitted := 0
	for simID := 0; simID < config.NumSimulations; simID++ {
		job := SimulationJob{
			SimID:  simID + 1,
			Config: config,
		}
		if pool.AddJob(job) {
			jobsSubmitted++
		} else {
			fmt.Printf("Warning: Failed to submit job %d, pool is shutting down\n", simID+1)
			break
		}
	}

	fmt.Printf("Successfully submitted %d/%d jobs\n", jobsSubmitted, config.NumSimulations)

	// Stop the worker pool and wait for completion
	pool.Stop()

	// Signal monitoring goroutines to stop
	monitorCancel()

	// Wait for all results to be collected with timeout
	timeout := time.After(30 * time.Second)
	select {
	case <-resultsDone:
		// Results collected successfully
	case <-timeout:
		fmt.Println("Warning: Timeout waiting for results collection")
	}

	select {
	case <-memoryMonitorDone:
		// Memory monitor finished
	case <-timeout:
		fmt.Println("Warning: Timeout waiting for memory monitor")
	}

	select {
	case <-progressDone:
		// Progress reporter finished
	case <-timeout:
		fmt.Println("Warning: Timeout waiting for progress reporter")
	}

	// Calculate final statistics with atomic reads
	stats.ExecutionTime = time.Since(startTime)
	completedSims := atomic.LoadInt64(&stats.CompletedSims)
	if completedSims > 0 {
		team1Wins := atomic.LoadInt64(&stats.Team1Wins)
		team2Wins := atomic.LoadInt64(&stats.Team2Wins)
		overtimeGames := atomic.LoadInt64(&stats.OvertimeGames)
		totalRounds := atomic.LoadInt64(&stats.TotalRounds)

		stats.Team1WinRate = float64(team1Wins) / float64(completedSims) * 100
		stats.Team2WinRate = float64(team2Wins) / float64(completedSims) * 100
		stats.OvertimeRate = float64(overtimeGames) / float64(completedSims) * 100
		stats.AverageRounds = float64(totalRounds) / float64(completedSims)
		stats.ProcessingRate = float64(completedSims) / stats.ExecutionTime.Seconds()
	}

	// Final garbage collection and memory stats
	runtime.GC()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	finalMemoryMB := memStats.Sys / (1024 * 1024)

	// Update peak memory if final memory is higher (atomic compare-and-swap)
	for {
		currentPeak := atomic.LoadUint64(&stats.PeakMemoryUsage)
		if finalMemoryMB <= currentPeak {
			break
		}
		if atomic.CompareAndSwapUint64(&stats.PeakMemoryUsage, currentPeak, finalMemoryMB) {
			break
		}
	}

	// Export summary statistics
	summaryPath := filepath.Join(resultsDir, "simulation_summary.json")
	if err := exportSummary(stats, summaryPath); err != nil {
		fmt.Printf("Warning: Failed to export summary: %v\n", err)
	}

	// Print final results
	printFinalStats(stats)

	// Print export information
	fmt.Printf("\nResults exported to: %s/\n", resultsDir)
	if config.ExportResults {
		fmt.Printf("- Individual game results: %d JSON files\n", atomic.LoadInt64(&stats.CompletedSims))
	}
	fmt.Println("- Summary statistics: simulation_summary.json")

	return nil
} // collectResults processes simulation results and updates statistics
func collectResults(results <-chan SimulationResult, stats *SimulationStats, totalSims int, exportResults bool, resultsDir string) {
	processedCount := int64(0)

	for result := range results {
		atomic.AddInt64(&processedCount, 1)

		if result.Error != nil {
			fmt.Printf("Simulation %s failed: %v\n", result.GameID, result.Error)
			// Still count failed simulations as completed for monitoring purposes
			atomic.AddInt64(&stats.CompletedSims, 1)
			continue
		}

		// Export individual result if requested
		if exportResults && result.GameID != "" {
			go func(res SimulationResult) {
				filename := filepath.Join(resultsDir, fmt.Sprintf("%s.json", res.GameID))
				data, err := json.MarshalIndent(res, "", "  ")
				if err == nil {
					os.WriteFile(filename, data, 0644)
				}
			}(result)
		}

		// Update statistics for successful simulations
		atomic.AddInt64(&stats.CompletedSims, 1)
		atomic.AddInt64(&stats.TotalRounds, int64(result.TotalRounds))

		if result.Team1Won {
			atomic.AddInt64(&stats.Team1Wins, 1)
		} else {
			atomic.AddInt64(&stats.Team2Wins, 1)
		}

		if result.WentToOvertime {
			atomic.AddInt64(&stats.OvertimeGames, 1)
		}

		// Update score distribution with bounds checking to prevent unlimited growth
		scoreKey := fmt.Sprintf("%d-%d", result.Team1Score, result.Team2Score)
		stats.scoreMutex.Lock()
		// Limit score distribution to prevent memory issues
		if len(stats.ScoreDistribution) < 1000 || stats.ScoreDistribution[scoreKey] > 0 {
			stats.ScoreDistribution[scoreKey]++
		}
		stats.scoreMutex.Unlock()
	}
}

// monitorMemoryUsageWithContext monitors memory usage and forces GC when needed
func monitorMemoryUsageWithContext(ctx context.Context, stats *SimulationStats, memoryLimit int) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var memStats runtime.MemStats

	for {
		select {
		case <-ticker.C:
			runtime.ReadMemStats(&memStats)
			currentMB := memStats.Alloc / (1024 * 1024)

			// Update peak memory usage with atomic compare-and-swap
			for {
				currentPeak := atomic.LoadUint64(&stats.PeakMemoryUsage)
				if currentMB <= currentPeak {
					break
				}
				if atomic.CompareAndSwapUint64(&stats.PeakMemoryUsage, currentPeak, currentMB) {
					break
				}
			}

			// Force GC if memory usage is too high
			if currentMB > uint64(memoryLimit) {
				runtime.GC()
				atomic.AddUint32(&stats.TotalGCRuns, 1)
			}

			// Break if simulations are complete
			if atomic.LoadInt64(&stats.CompletedSims) >= stats.TotalSimulations {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// reportProgressWithContext periodically reports simulation progress
func reportProgressWithContext(ctx context.Context, stats *SimulationStats, totalSims int, startTime time.Time, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			completed := atomic.LoadInt64(&stats.CompletedSims)
			elapsed := time.Since(startTime)
			rate := float64(completed) / elapsed.Seconds()

			fmt.Printf("Progress: %d/%d (%.1f%%) - Rate: %.1f sims/sec - Elapsed: %s\n",
				completed, totalSims,
				float64(completed)/float64(totalSims)*100,
				rate, elapsed.Round(time.Second))

			// Break if simulations are complete
			if completed >= int64(totalSims) {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// exportSummary exports simulation statistics to JSON
func exportSummary(stats *SimulationStats, path string) error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// printFinalStats prints final simulation statistics
func printFinalStats(stats *SimulationStats) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("SIMULATION RESULTS SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total Simulations: %d\n", stats.CompletedSims)
	fmt.Printf("Execution Time: %s\n", stats.ExecutionTime.Round(time.Second))
	fmt.Printf("Processing Rate: %.1f simulations/second\n", stats.ProcessingRate)
	fmt.Printf("Peak Memory Usage: %d MB\n", stats.PeakMemoryUsage)
	fmt.Printf("Total GC Runs: %d\n", stats.TotalGCRuns)
	fmt.Println()
	fmt.Printf("Team 1 Wins: %d (%.1f%%)\n", stats.Team1Wins, stats.Team1WinRate)
	fmt.Printf("Team 2 Wins: %d (%.1f%%)\n", stats.Team2Wins, stats.Team2WinRate)
	fmt.Printf("Average Rounds per Game: %.1f\n", stats.AverageRounds)
	fmt.Printf("Overtime Rate: %.1f%%\n", stats.OvertimeRate)

	if len(stats.ScoreDistribution) > 0 {
		fmt.Println("\nTop Score Distributions:")
		// Print top 5 most common scores
		// This is a simplified display - in production you might want to sort properly
		count := 0
		for score, freq := range stats.ScoreDistribution {
			if count >= 5 {
				break
			}
			percentage := float64(freq) / float64(stats.CompletedSims) * 100
			fmt.Printf("  %s: %d games (%.1f%%)\n", score, freq, percentage)
			count++
		}
	}
	fmt.Println(strings.Repeat("=", 60))
}
