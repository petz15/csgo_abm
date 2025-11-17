package util

import (
	"csgo_abm/internal/engine"
	"encoding/csv"
	"fmt"
	"os"
)

// ExportRoundsToCSV exports minimal round-by-round data as CSV (much smaller file size than JSON)
func ExportRoundsToCSV(game *engine.Game, path string) error {
	if path == "" {
		path = "rounds.csv"
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header with metadata as comments
	writer.Write([]string{"# GameID:", game.ID})
	writer.Write([]string{"# Team1:", game.Team1.Name, game.Team1.Strategy})
	writer.Write([]string{"# Team2:", game.Team2.Name, game.Team2.Strategy})
	writer.Write([]string{"# FinalScore:", fmt.Sprintf("%d-%d", game.Score[0], game.Score[1])})
	writer.Write([]string{"# Overtime:", fmt.Sprintf("%v", game.OT)})
	writer.Write([]string{""}) // Empty line

	// Write column headers
	headers := []string{
		"round_number",
		"winner",
		"team1_side",
		"team1_score",
		"team2_score",
		"team1_spent",
		"team2_spent",
		"team1_fte_eq",
		"team2_fte_eq",
		"team1_earned",
		"team2_earned",
		"team1_funds",
		"team2_funds",
		"team1_consec_loss",
		"team2_consec_loss",
		"team1_loss_bonus_lvl",
		"team2_loss_bonus_lvl",
	}
	writer.Write(headers)

	// Write data rows
	for i := range game.Rounds {
		round := &game.Rounds[i]
		team1Data := game.Team1.RoundData[i]
		team2Data := game.Team2.RoundData[i]

		team1Side := "T"
		if round.IsT1CT {
			team1Side = "CT"
		}

		winner := "Team2"
		if round.IsT1WinnerTeam {
			winner = "Team1"
		}

		row := []string{
			fmt.Sprintf("%d", round.RoundNumber),
			winner,
			team1Side,
			fmt.Sprintf("%d", team1Data.Score_End),
			fmt.Sprintf("%d", team2Data.Score_End),
			fmt.Sprintf("%.2f", team1Data.Spent),
			fmt.Sprintf("%.2f", team2Data.Spent),
			fmt.Sprintf("%.2f", team1Data.FTE_Eq_value),
			fmt.Sprintf("%.2f", team2Data.FTE_Eq_value),
			fmt.Sprintf("%.2f", team1Data.Earned),
			fmt.Sprintf("%.2f", team2Data.Earned),
			fmt.Sprintf("%.2f", team1Data.Funds),
			fmt.Sprintf("%.2f", team2Data.Funds),
			fmt.Sprintf("%d", team1Data.Consecutiveloss),
			fmt.Sprintf("%d", team2Data.Consecutiveloss),
			fmt.Sprintf("%d", team1Data.LossBonusLevel),
			fmt.Sprintf("%d", team2Data.LossBonusLevel),
		}
		writer.Write(row)
	}

	return writer.Error()
}
