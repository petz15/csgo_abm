package main

import (
	"context"
	"csgo_abm/internal/analysis"
	"csgo_abm/internal/tournament"
	"fmt"
	"strings"
	"time"
)

func runTournament(cfg *SimulationConfig, custom *CustomConfig, strategiesCSV string, format string, bestOf int) error {
	// Parse strategies list
	list := []string{}
	for _, s := range strings.Split(strategiesCSV, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if err := ValidateStrategies(s, s); err != nil {
			return fmt.Errorf("unknown strategy '%s': %w", s, err)
		}
		list = append(list, s)
	}
	if len(list) < 2 {
		return fmt.Errorf("need at least two strategies for a tournament")
	}

	tcfg := tournament.TournamentConfig{
		Format: tournament.Format(format),
		Series: tournament.SeriesSpec{
			BestOf:         bestOf,
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
		res, err := tournament.RunSeries(context.Background(), m, custom.GameRules, tcfg.Series)
		if err != nil {
			return err
		}
		seriesResults = append(seriesResults, res)
	}

	standings := tournament.ComputeStandings(list, seriesResults)

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
