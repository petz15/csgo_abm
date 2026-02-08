package analysis

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"csgo_abm/internal/tournament"
)

// ExportTournamentSummary writes tournament matches, series, and standings to JSON and standings CSV
func ExportTournamentSummary(dir string, matches []tournament.MatchSpec, series []tournament.SeriesResult, standings tournament.Standings) error {
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
	for _, r := range standings.Rows {
		w.Write([]string{
			r.Strategy,
			strconv.Itoa(r.Wins),
			strconv.Itoa(r.Losses),
			strconv.Itoa(r.MapWins),
			strconv.Itoa(r.MapLoss),
		})
	}

	// Matrix CSV (win percentage per matchup, using standings order for rows/cols)
	names := make([]string, 0, len(standings.Rows))
	for _, r := range standings.Rows {
		names = append(names, r.Strategy)
	}
	n := len(names)
	idx := make(map[string]int, n)
	for i, name := range names {
		idx[name] = i
	}
	wins := make([][]int, n)
	totals := make([][]int, n)
	for i := 0; i < n; i++ {
		wins[i] = make([]int, n)
		totals[i] = make([]int, n)
	}

	for _, ser := range series {
		i := idx[ser.Match.Team1Strategy]
		j := idx[ser.Match.Team2Strategy]
		team1Wins := 0
		team2Wins := 0
		for _, g := range ser.GameResults {
			if g.T1Wins {
				team1Wins++
			} else {
				team2Wins++
			}
		}
		totalGames := team1Wins + team2Wins
		// Populate both cells symmetrically
		totals[i][j] = totalGames
		totals[j][i] = totalGames
		wins[i][j] = team1Wins
		wins[j][i] = team2Wins
	}

	mf, err := os.Create(filepath.Join(dir, "tournament_matrix.csv"))
	if err != nil {
		return err
	}
	defer mf.Close()
	mw := csv.NewWriter(mf)
	defer mw.Flush()

	header := make([]string, 0, n+1)
	header = append(header, "strategy")
	header = append(header, names...)
	mw.Write(header)

	for i := 0; i < n; i++ {
		row := make([]string, 0, n+1)
		row = append(row, names[i])
		for j := 0; j < n; j++ {
			if i == j {
				row = append(row, "-")
				continue
			}
			t := totals[i][j]
			if t == 0 {
				row = append(row, "0/0 (0.00%)")
				continue
			}
			wv := wins[i][j]
			pct := 100.0 * float64(wv) / float64(t)
			row = append(row, fmt.Sprintf("%d/%d (%.2f%%)", wv, t, pct))
		}
		mw.Write(row)
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
