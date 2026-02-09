package main

import (
	"csgo_abm/internal/analysis"
	"csgo_abm/internal/strategy"
	"csgo_abm/internal/tournament"
	"fmt"
	"os"
	"strings"
)

type MatchSpec struct {
	Team1Strategy string
	Team2Strategy string
}

type MatchResult struct {
	Team1Wins int
	Team2Wins int
}

// RoundRobinSchedule creates all matchups for a round-robin tournament
func RoundRobinSchedule(strategies []string) []MatchSpec {
	var matches []MatchSpec
	for i := 0; i < len(strategies); i++ {
		for j := i + 1; j < len(strategies); j++ {
			// A vs B
			matches = append(matches, MatchSpec{
				Team1Strategy: strategies[i],
				Team2Strategy: strategies[j],
			})
			// B vs A (for balance, though starting side is randomized) -> not necessary, one is enough since sides are randomized in each game
			// matches = append(matches, MatchSpec{
			// 	Team1Strategy: strategies[j],
			// 	Team2Strategy: strategies[i],
			// })
		}
	}
	return matches
}

func runTournament(cfg *SimulationConfig, custom *CustomConfig, strategiesCSV string, format string, games int) error {
	// Parse strategies list
	list := []string{}
	for _, s := range strings.Split(strategiesCSV, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		list = append(list, s)
	}
	if len(list) < 2 {
		return fmt.Errorf("need at least two strategies for a tournament")
	}

	// Validate all strategies upfront
	fmt.Println("Validating strategies...")
	for _, strat := range list {
		if err := strategy.ValidateStrategy(strat); err != nil {
			return fmt.Errorf("❌ %v", err)
		}
		fmt.Printf("  ✓ %s\n", strat)
	}

	var matches []MatchSpec
	switch format {
	case "roundrobin":
		matches = RoundRobinSchedule(list)
	default:
		return fmt.Errorf("unsupported tournament format: %s", format)
	}

	fmt.Printf("Running tournament with %d strategies, %d matchups, %d games each...\n", len(list), len(matches), games)

	// Run all matchups and collect results
	matchResults := make([]MatchResult, len(matches))

	for i, m := range matches {
		fmt.Printf("\nMatchup %d/%d: %s vs %s\n", i+1, len(matches), m.Team1Strategy, m.Team2Strategy)

		// Create a unique folder for this matchup to avoid CSV file conflicts
		matchupFolder := fmt.Sprintf("%s/matchup_%03d_%s_vs_%s", cfg.Exportpath, i+1, m.Team1Strategy, m.Team2Strategy)
		if err := os.MkdirAll(matchupFolder, 0755); err != nil {
			return fmt.Errorf("failed to create matchup folder: %v", err)
		}

		// Create a simulation config for this matchup
		matchConfig := SimulationConfig{
			NumSimulations:        games,
			MaxConcurrent:         cfg.MaxConcurrent,
			MemoryLimit:           cfg.MemoryLimit,
			Team1Name:             m.Team1Strategy,
			Team1Strategy:         m.Team1Strategy,
			Team2Name:             m.Team2Strategy,
			Team2Strategy:         m.Team2Strategy,
			GameRules:             custom.GameRules,
			ExportDetailedResults: false,
			ExportRounds:          false,
			Sequential:            cfg.Sequential,
			SuppressOutput:        true,              // Suppress output during tournament
			CSVExportMode:         cfg.CSVExportMode, // Use the tournament's CSV export mode
			Exportpath:            matchupFolder,     // Each matchup gets its own folder
		}

		// Run the simulations for this matchup
		var stats *analysis.SimulationStats
		var err error
		if cfg.Sequential {
			// For sequential, we need to capture stats differently
			tempStats := analysis.NewStats(games, "sequential")
			for g := 0; g < games; g++ {
				simPrefix := fmt.Sprintf("tournament_%s_vs_%s_game_%d_", m.Team1Strategy, m.Team2Strategy, g)
				result, gameErr := StartGame_default(
					m.Team1Strategy,
					m.Team1Strategy,
					m.Team2Strategy,
					m.Team2Strategy,
					custom.GameRules,
					simPrefix,
					false,
					false,
					cfg.CSVExportMode, // Use the tournament's CSV export mode
					matchupFolder,     // Use matchup-specific folder
				)
				if gameErr != nil {
					continue
				}
				updateglobalstats(tempStats, result)
			}
			stats = tempStats
		} else {
			// Run concurrent simulations
			stats, err = RunParallelSimulations(matchConfig)
			if err != nil {
				return fmt.Errorf("matchup %d/%d (%s vs %s) failed after running %d games: %w",
					i+1, len(matches), m.Team1Strategy, m.Team2Strategy, games, err)
			}
		}

		matchResults[i] = MatchResult{
			Team1Wins: int(stats.Team1Wins),
			Team2Wins: int(stats.Team2Wins),
		}

		fmt.Printf("  Result: %s won %d, %s won %d\n",
			m.Team1Strategy, matchResults[i].Team1Wins,
			m.Team2Strategy, matchResults[i].Team2Wins)
	}

	// Compute standings
	standings := computeStandings(list, matches, matchResults)

	// Build and print CLI matrix
	printTournamentMatrix(list, matches, matchResults)

	// Export results
	resdir, err := analysis.CreateResultsDirectoryAt(cfg.Exportpath)
	if err != nil {
		return err
	}

	// Convert to tournament package types for export compatibility
	tournamentMatches := make([]tournament.MatchSpec, len(matches))
	for i, m := range matches {
		tournamentMatches[i] = tournament.MatchSpec{
			Team1Name:     m.Team1Strategy,
			Team1Strategy: m.Team1Strategy,
			Team2Name:     m.Team2Strategy,
			Team2Strategy: m.Team2Strategy,
		}
	}

	seriesResults := make([]tournament.SeriesResult, len(matches))
	for i := range matches {
		seriesResults[i] = tournament.SeriesResult{
			Match: tournamentMatches[i],
		}
		// Add game results
		for g := 0; g < matchResults[i].Team1Wins; g++ {
			seriesResults[i].GameResults = append(seriesResults[i].GameResults, tournament.GameOutcome{
				T1Wins: true,
				Score:  [2]int{16, 0}, // Placeholder scores
			})
		}
		for g := 0; g < matchResults[i].Team2Wins; g++ {
			seriesResults[i].GameResults = append(seriesResults[i].GameResults, tournament.GameOutcome{
				T1Wins: false,
				Score:  [2]int{0, 16}, // Placeholder scores
			})
		}
	}

	if err := analysis.ExportTournamentSummary(resdir, tournamentMatches, seriesResults, standings); err != nil {
		return err
	}

	fmt.Printf("\n✅ Tournament finished. Results exported to: %s\n", resdir)
	return nil
}

func computeStandings(strategies []string, matches []MatchSpec, results []MatchResult) tournament.Standings {
	idx := make(map[string]int)
	rows := make([]tournament.StandingsRow, len(strategies))
	for i, s := range strategies {
		idx[s] = i
		rows[i] = tournament.StandingsRow{Strategy: s}
	}

	for i, m := range matches {
		i1 := idx[m.Team1Strategy]
		i2 := idx[m.Team2Strategy]

		rows[i1].MapWins += results[i].Team1Wins
		rows[i1].MapLoss += results[i].Team2Wins
		rows[i2].MapWins += results[i].Team2Wins
		rows[i2].MapLoss += results[i].Team1Wins

		if results[i].Team1Wins > results[i].Team2Wins {
			rows[i1].Wins++
			rows[i2].Losses++
		} else {
			rows[i2].Wins++
			rows[i1].Losses++
		}
	}

	return tournament.Standings{Rows: rows}
}

func printTournamentMatrix(strategies []string, matches []MatchSpec, results []MatchResult) {
	n := len(strategies)
	idx := make(map[string]int, n)
	for i, name := range strategies {
		idx[name] = i
	}

	wins := make([][]int, n)
	totals := make([][]int, n)
	cells := make([][]string, n)
	for i := 0; i < n; i++ {
		wins[i] = make([]int, n)
		totals[i] = make([]int, n)
		cells[i] = make([]string, n)
	}

	for i, m := range matches {
		i1 := idx[m.Team1Strategy]
		i2 := idx[m.Team2Strategy]
		totalGames := results[i].Team1Wins + results[i].Team2Wins
		// Populate both directions of the matchup
		totals[i1][i2] += totalGames
		totals[i2][i1] += totalGames
		wins[i1][i2] += results[i].Team1Wins
		wins[i2][i1] += results[i].Team2Wins
	}

	// Precompute cell strings and column widths
	headers := make([]string, n+1)
	headers[0] = "strategy"
	copy(headers[1:], strategies)
	colW := make([]int, n+1)
	colW[0] = len(headers[0])
	for i := 0; i < n; i++ {
		if len(strategies[i]) > colW[0] {
			colW[0] = len(strategies[i])
		}
	}
	for j := 0; j < n; j++ {
		colW[j+1] = len(headers[j+1])
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i == j {
				cells[i][j] = "-"
			} else {
				t := totals[i][j]
				if t == 0 {
					cells[i][j] = "0/0 (0.00%)"
				} else {
					w := wins[i][j]
					pct := 100.0 * float64(w) / float64(t)
					cells[i][j] = fmt.Sprintf("%d/%d (%.2f%%)", w, t, pct)
				}
			}
			if len(cells[i][j]) > colW[j+1] {
				colW[j+1] = len(cells[i][j])
			}
		}
	}

	// Print header
	fmt.Println()
	fmt.Println("Tournament win-rate matrix:")
	for j := 0; j <= n; j++ {
		fmt.Printf("%-*s ", colW[j], headers[j])
	}
	fmt.Println()

	// Print rows
	for i := 0; i < n; i++ {
		fmt.Printf("%-*s ", colW[0], strategies[i])
		for j := 0; j < n; j++ {
			fmt.Printf("%-*s ", colW[j+1], cells[i][j])
		}
		fmt.Println()
	}
}
