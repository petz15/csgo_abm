package main

import (
	"context"
	"csgo_abm/internal/analysis"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

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
	stats   *analysis.SimulationStats
}

// SimulationJob represents a single simulation job
type SimulationJob struct {
	SimID  int
	Config analysis.SimulationConfig
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int, stats *analysis.SimulationStats) *WorkerPool {
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
		gameResult, err := StartGame_default(
			job.Config.Team1Name,
			job.Config.Team1Strategy,
			job.Config.Team2Name,
			job.Config.Team2Strategy,
			job.Config.GameRules,
			simPrefix,
			job.Config.ExportDetailedResults,
			job.Config.Exportpath,
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
}

// RunParallelSimulations orchestrates the execution of multiple simulations
func RunParallelSimulations(config analysis.SimulationConfig) error {
	startTime := time.Now()

	// Initialize statistics
	stats := analysis.NewSimulationStats(analysis.SimulationConfig{
		NumSimulations:        config.NumSimulations,
		Team1Name:             config.Team1Name,
		Team2Name:             config.Team2Name,
		Team1Strategy:         config.Team1Strategy,
		Team2Strategy:         config.Team2Strategy,
		GameRules:             config.GameRules,
		ExportDetailedResults: config.ExportDetailedResults,
		Sequential:            false,             // This is concurrent mode
		Exportpath:            config.Exportpath, // Use the export path from main config
	})

	fmt.Printf("Starting %d simulations with %d concurrent workers...\n",
		config.NumSimulations, config.MaxConcurrent)
	fmt.Printf("Memory limit: %d MB\n", config.MemoryLimit)

	if config.ExportDetailedResults {
		fmt.Printf("Individual result export: ENABLED (results will be saved to %s/)\n", config.Exportpath)
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
		collectResults(pool.results, stats, config.NumSimulations)
	}()

	// Track memory usage
	memoryMonitorDone := make(chan bool)
	go func() {
		defer close(memoryMonitorDone)
		monitorMemoryUsageWithContext(monitorCtx, stats, config.MemoryLimit)
	}()

	// Progress reporter
	progressDone := make(chan bool)
	progressInterval := 10 * time.Second
	// For very large simulations, report progress less frequently
	if config.NumSimulations >= 100000 {
		progressInterval = 30 * time.Second
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

	// Calculate final statistics using unified analysis package
	stats.ExecutionTime = time.Since(startTime)
	stats.CalculateFinalStats()

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
	summaryPath := filepath.Join(config.Exportpath, "simulation_summary.json")
	if err := exportSummary(stats, summaryPath); err != nil {
		fmt.Printf("Warning: Failed to export summary: %v\n", err)
	}

	// Print final results
	analysis.PrintEnhancedStats(stats)

	// Print export information
	fmt.Printf("\nResults exported to: %s/\n", config.Exportpath)
	if config.ExportDetailedResults {
		fmt.Printf("- Individual game results: %d JSON files\n", atomic.LoadInt64(&stats.CompletedSims))
	}
	fmt.Println("- Summary statistics: simulation_summary.json")

	return nil
} // collectResults processes simulation results and updates statistics
func collectResults(results <-chan SimulationResult, stats *analysis.SimulationStats, totalSims int) {
	processedCount := int64(0)

	for result := range results {
		atomic.AddInt64(&processedCount, 1)

		if result.Error != nil {
			fmt.Printf("Simulation %s failed: %v\n", result.GameID, result.Error)
			// Count failed simulations
			stats.UpdateFailedSimulation()
			continue
		}

		// Update statistics for successful simulations
		stats.UpdateGameResult(
			result.Team1Won,
			result.Team1Score,
			result.Team2Score,
			result.TotalRounds,
			result.WentToOvertime,
			0, // responseTime - not tracked in current implementation
		)
	}
}

// monitorMemoryUsageWithContext monitors memory usage and forces GC when needed
func monitorMemoryUsageWithContext(ctx context.Context, stats *analysis.SimulationStats, memoryLimit int) {
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
func reportProgressWithContext(ctx context.Context, stats *analysis.SimulationStats, totalSims int, startTime time.Time, interval time.Duration) {
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
func exportSummary(stats *analysis.SimulationStats, path string) error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
