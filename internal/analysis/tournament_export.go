package analysis

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
)

// ExportTournamentSummary writes tournament matches, series, and standings to JSON and standings CSV
func ExportTournamentSummary(dir string, matches interface{}, series interface{}, standings interface{}) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// JSON dump
	if err := writeJSON(filepath.Join(dir, "tournament_summary.json"), map[string]interface{}{
		"matches":   matches,
		"series":    series,
		"standings": standings,
	}); err != nil {
		return err
	}
	// CSV standings (best-effort if type provides GetRows)
	f, err := os.Create(filepath.Join(dir, "tournament_standings.csv"))
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	w.Write([]string{"strategy", "wins", "losses", "map_wins", "map_losses"})
	// duck-typed accessor used by internal/tournament
	if s, ok := standings.(interface {
		GetRows() []struct {
			Strategy string
			Wins     int
			Losses   int
			MapWins  int
			MapLoss  int
		}
	}); ok {
		for _, r := range s.GetRows() {
			w.Write([]string{
				r.Strategy,
				strconv.Itoa(r.Wins),
				strconv.Itoa(r.Losses),
				strconv.Itoa(r.MapWins),
				strconv.Itoa(r.MapLoss),
			})
		}
	}
	return nil
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}
