package main

import (
	"context"
	"csgo_abm/internal/analysis"
	"csgo_abm/internal/tournament"
	"fmt"
	"strings"
	"time"
)

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

	tcfg := tournament.TournamentConfig{
		Format: tournament.Format(format),
		Series: tournament.SeriesSpec{
			NumGames:       games,
			Seed:           time.Now().UnixNano(),
			MaxConcurrent:  cfg.MaxConcurrent,
			TimeoutPerGame: 0,
		},
		Participants: list,
	}
	var matches []tournament.MatchSpec
	switch tcfg.Format {
	case tournament.FormatRoundRobin:
		matches = tournament.RoundRobinSchedule(tcfg.Participants)
	default:
		return fmt.Errorf("unsupported tournament format: %s", format)
	}

	// Simple sequential execution; can be parallelized later
	seriesResults := make([]tournament.SeriesResult, 0, len(matches))
	for _, m := range matches {
		res, err := tournament.RunMatchup(context.Background(), m, custom.GameRules, tcfg.Series)
		if err != nil {
			return err
		}
		seriesResults = append(seriesResults, res)
	}

	standings := tournament.ComputeStandings(list, seriesResults)

	// Build and print CLI matrix (wins/total and win%)
	names := make([]string, len(list))
	copy(names, list)
	n := len(names)
	idx := make(map[string]int, n)
	for i, name := range names {
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
	for _, ser := range seriesResults {
		i := idx[ser.Match.Team1Strategy]
		j := idx[ser.Match.Team2Strategy]
		for _, g := range ser.GameResults {
			totals[i][j]++
			totals[j][i]++
			if g.T1Wins {
				wins[i][j]++
			} else {
				wins[j][i]++
			}
		}
	}
	// Precompute cell strings and column widths
	headers := make([]string, n+1)
	headers[0] = "strategy"
	copy(headers[1:], names)
	colW := make([]int, n+1)
	colW[0] = len(headers[0])
	for i := 0; i < n; i++ {
		if len(names[i]) > colW[0] {
			colW[0] = len(names[i])
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
		fmt.Printf("%-*s ", colW[0], names[i])
		for j := 0; j < n; j++ {
			fmt.Printf("%-*s ", colW[j+1], cells[i][j])
		}
		fmt.Println()
	}

	// Export results under configured export path
	resdir, err := analysis.CreateResultsDirectoryAt(cfg.Exportpath)
	if err != nil {
		return err
	}
	if err := analysis.ExportTournamentSummary(resdir, matches, seriesResults, standings); err != nil {
		return err
	}
	fmt.Printf("✅ Tournament finished. Results exported to: %s\n", resdir)
	return nil
}
