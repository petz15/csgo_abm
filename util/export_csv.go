package util

import (
	"dbg_abm/internal/engine"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// joinFloatSlice converts a slice of float64 to a pipe-delimited string with 2 decimal precision.
func joinFloatSlice(vals []float64) string {
	if len(vals) == 0 {
		return ""
	}
	b := strings.Builder{}
	for i, v := range vals {
		if i > 0 {
			b.WriteByte('|')
		}
		b.WriteString(fmt.Sprintf("%.2f", v))
	}
	return b.String()
}

// 1. Export each game individually with all round data (all info per row)
func ExportGameAllDataCSV(game *engine.Game, path string) error {
	if path == "" {
		path = fmt.Sprintf("%s.csv", game.ID)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// Use semicolon as separator
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header row
	headers := []string{
		"round_number",
		"is_t1_ct",
		"is_t1_winner",
		"is_ot",
		"is_ct_winner",
		"outcome_reason_code",
		"outcome_ct_wins",
		"outcome_bomb_planted",
		"outcome_ct_survivors",
		"outcome_t_survivors",
		"csf",
		"csf_key",
		"ct_equipment_share_per_player",
		"t_equipment_share_per_player",
		"ct_equipment_per_player",
		"t_equipment_per_player",
		"rng_csf",
		"rng_round_outcome",
		"rng_bombplant",
		"rng_survivors_ct",
		"rng_survivors_t",
		"rng_equipment_ct",
		"rng_equipment_t",
		"t1_funds",
		"t1_funds_start",
		"t1_earned",
		"t1_rs_eq_value",
		"t1_fte_eq_value",
		"t1_re_eq_value",
		"t1_survivors",
		"t1_score_start",
		"t1_score_end",
		"t1_consecutive_loss",
		"t1_consecutive_losses_start",
		"t1_consecutive_wins",
		"t1_consecutive_wins_start",
		"t1_loss_bonus_level",
		"t1_spent",
		"t2_funds",
		"t2_funds_start",
		"t2_earned",
		"t2_rs_eq_value",
		"t2_fte_eq_value",
		"t2_re_eq_value",
		"t2_survivors",
		"t2_score_start",
		"t2_score_end",
		"t2_consecutive_loss",
		"t2_consecutive_losses_start",
		"t2_consecutive_wins",
		"t2_consecutive_wins_start",
		"t2_loss_bonus_level",
		"t2_spent",
		"t1_name",
		"t1_strategy",
		"t2_name",
		"t2_strategy",
		"game_id",
	}
	writer.Write(headers)

	for i := range game.Rounds {
		round := &game.Rounds[i]
		t1 := game.Team1.RoundData[i]
		t2 := game.Team2.RoundData[i]
		isCTWinner := round.IsT1CT == round.IsT1WinnerTeam
		row := []string{
			fmt.Sprintf("%d", round.RoundNumber),
			fmt.Sprintf("%t", round.IsT1CT),
			fmt.Sprintf("%t", round.IsT1WinnerTeam),
			fmt.Sprintf("%t", round.OT),
			fmt.Sprintf("%t", isCTWinner),
			fmt.Sprintf("%d", round.Calc_Outcome.ReasonCode),
			fmt.Sprintf("%t", round.Calc_Outcome.CTWins),
			fmt.Sprintf("%t", round.Calc_Outcome.BombPlanted),
			fmt.Sprintf("%d", round.Calc_Outcome.CTSurvivors),
			fmt.Sprintf("%d", round.Calc_Outcome.TSurvivors),
			fmt.Sprintf("%.6f", round.Calc_Outcome.CSF),
			round.Calc_Outcome.CSFKey,
			joinFloatSlice(round.Calc_Outcome.CTEquipmentSharePerPlayer),
			joinFloatSlice(round.Calc_Outcome.TEquipmentSharePerPlayer),
			joinFloatSlice(round.Calc_Outcome.CTEquipmentPerPlayer),
			joinFloatSlice(round.Calc_Outcome.TEquipmentPerPlayer),
			fmt.Sprintf("%.6f", round.Calc_Outcome.StochasticValues.RNG_CSF),
			fmt.Sprintf("%.6f", round.Calc_Outcome.StochasticValues.RNG_RoundOutcome),
			fmt.Sprintf("%.6f", round.Calc_Outcome.StochasticValues.RNG_Bombplant),
			fmt.Sprintf("%.6f", round.Calc_Outcome.StochasticValues.RNG_SurvivorsCT),
			fmt.Sprintf("%.6f", round.Calc_Outcome.StochasticValues.RNG_SurvivorsT),
			joinFloatSlice(round.Calc_Outcome.StochasticValues.RNG_EquipmentCT),
			joinFloatSlice(round.Calc_Outcome.StochasticValues.RNG_EquipmentT),
			fmt.Sprintf("%.2f", t1.Funds),
			fmt.Sprintf("%.2f", t1.Funds_start),
			fmt.Sprintf("%.2f", t1.Earned),
			fmt.Sprintf("%.2f", t1.RS_Eq_value),
			fmt.Sprintf("%.2f", t1.FTE_Eq_value),
			fmt.Sprintf("%.2f", t1.RE_Eq_value),
			fmt.Sprintf("%d", t1.Survivors),
			fmt.Sprintf("%d", t1.Score_Start),
			fmt.Sprintf("%d", t1.Score_End),
			fmt.Sprintf("%d", t1.Consecutiveloss),
			fmt.Sprintf("%d", t1.Consecutiveloss_start),
			fmt.Sprintf("%d", t1.Consecutivewins),
			fmt.Sprintf("%d", t1.Consecutivewins_start),
			fmt.Sprintf("%d", t1.LossBonusLevel),
			fmt.Sprintf("%.2f", t1.Spent),
			fmt.Sprintf("%.2f", t2.Funds),
			fmt.Sprintf("%.2f", t2.Funds_start),
			fmt.Sprintf("%.2f", t2.Earned),
			fmt.Sprintf("%.2f", t2.RS_Eq_value),
			fmt.Sprintf("%.2f", t2.FTE_Eq_value),
			fmt.Sprintf("%.2f", t2.RE_Eq_value),
			fmt.Sprintf("%d", t2.Survivors),
			fmt.Sprintf("%d", t2.Score_Start),
			fmt.Sprintf("%d", t2.Score_End),
			fmt.Sprintf("%d", t2.Consecutiveloss),
			fmt.Sprintf("%d", t2.Consecutiveloss_start),
			fmt.Sprintf("%d", t2.Consecutivewins),
			fmt.Sprintf("%d", t2.Consecutivewins_start),
			fmt.Sprintf("%d", t2.LossBonusLevel),
			fmt.Sprintf("%.2f", t2.Spent),
			game.Team1.Name,
			game.Team1.Strategy,
			game.Team2.Name,
			game.Team2.Strategy,
			game.ID,
		}
		writer.Write(row)
	}
	return writer.Error()
}

// 2. Export all games in one file with all round data (all info per row)
func ExportAllGamesAllDataCSV(games []*engine.Game, path string) error {
	if path == "" {
		path = "all_games_full.csv"
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header row
	headers := []string{
		"round_number",
		"is_t1_ct",
		"is_t1_winner",
		"is_ot",
		"is_ct_winner",
		"outcome_reason_code",
		"outcome_ct_wins",
		"outcome_bomb_planted",
		"outcome_ct_survivors",
		"outcome_t_survivors",
		"csf",
		"csf_key",
		"ct_equipment_share_per_player",
		"t_equipment_share_per_player",
		"ct_equipment_per_player",
		"t_equipment_per_player",
		"rng_csf",
		"rng_round_outcome",
		"rng_bombplant",
		"rng_survivors_ct",
		"rng_survivors_t",
		"rng_equipment_ct",
		"rng_equipment_t",
		"t1_funds",
		"t1_funds_start",
		"t1_earned",
		"t1_rs_eq_value",
		"t1_fte_eq_value",
		"t1_re_eq_value",
		"t1_survivors",
		"t1_score_start",
		"t1_score_end",
		"t1_consecutive_loss",
		"t1_consecutive_losses_start",
		"t1_consecutive_wins",
		"t1_consecutive_wins_start",
		"t1_loss_bonus_level",
		"t1_spent",
		"t2_funds",
		"t2_funds_start",
		"t2_earned",
		"t2_rs_eq_value",
		"t2_fte_eq_value",
		"t2_re_eq_value",
		"t2_survivors",
		"t2_score_start",
		"t2_score_end",
		"t2_consecutive_loss",
		"t2_consecutive_losses_start",
		"t2_consecutive_wins",
		"t2_consecutive_wins_start",
		"t2_loss_bonus_level",
		"t2_spent",
		"t1_name",
		"t1_strategy",
		"t2_name",
		"t2_strategy",
		"game_id",
	}
	writer.Write(headers)

	for _, game := range games {
		for i := range game.Rounds {
			round := &game.Rounds[i]
			t1 := game.Team1.RoundData[i]
			t2 := game.Team2.RoundData[i]
			isCTWinner := round.IsT1CT == round.IsT1WinnerTeam
			row := []string{
				fmt.Sprintf("%d", round.RoundNumber),
				fmt.Sprintf("%t", round.IsT1CT),
				fmt.Sprintf("%t", round.IsT1WinnerTeam),
				fmt.Sprintf("%t", round.OT),
				fmt.Sprintf("%t", isCTWinner),
				fmt.Sprintf("%d", round.Calc_Outcome.ReasonCode),
				fmt.Sprintf("%t", round.Calc_Outcome.CTWins),
				fmt.Sprintf("%t", round.Calc_Outcome.BombPlanted),
				fmt.Sprintf("%d", round.Calc_Outcome.CTSurvivors),
				fmt.Sprintf("%d", round.Calc_Outcome.TSurvivors),
				fmt.Sprintf("%.6f", round.Calc_Outcome.CSF),
				round.Calc_Outcome.CSFKey,
				joinFloatSlice(round.Calc_Outcome.CTEquipmentSharePerPlayer),
				joinFloatSlice(round.Calc_Outcome.TEquipmentSharePerPlayer),
				joinFloatSlice(round.Calc_Outcome.CTEquipmentPerPlayer),
				joinFloatSlice(round.Calc_Outcome.TEquipmentPerPlayer),
				fmt.Sprintf("%.6f", round.Calc_Outcome.StochasticValues.RNG_CSF),
				fmt.Sprintf("%.6f", round.Calc_Outcome.StochasticValues.RNG_RoundOutcome),
				fmt.Sprintf("%.6f", round.Calc_Outcome.StochasticValues.RNG_Bombplant),
				fmt.Sprintf("%.6f", round.Calc_Outcome.StochasticValues.RNG_SurvivorsCT),
				fmt.Sprintf("%.6f", round.Calc_Outcome.StochasticValues.RNG_SurvivorsT),
				joinFloatSlice(round.Calc_Outcome.StochasticValues.RNG_EquipmentCT),
				joinFloatSlice(round.Calc_Outcome.StochasticValues.RNG_EquipmentT),
				fmt.Sprintf("%.2f", t1.Funds),
				fmt.Sprintf("%.2f", t1.Funds_start),
				fmt.Sprintf("%.2f", t1.Earned),
				fmt.Sprintf("%.2f", t1.RS_Eq_value),
				fmt.Sprintf("%.2f", t1.FTE_Eq_value),
				fmt.Sprintf("%.2f", t1.RE_Eq_value),
				fmt.Sprintf("%d", t1.Survivors),
				fmt.Sprintf("%d", t1.Score_Start),
				fmt.Sprintf("%d", t1.Score_End),
				fmt.Sprintf("%d", t1.Consecutiveloss),
				fmt.Sprintf("%d", t1.Consecutiveloss_start),
				fmt.Sprintf("%d", t1.Consecutivewins),
				fmt.Sprintf("%d", t1.Consecutivewins_start),
				fmt.Sprintf("%d", t1.LossBonusLevel),
				fmt.Sprintf("%.2f", t1.Spent),
				fmt.Sprintf("%.2f", t2.Funds),
				fmt.Sprintf("%.2f", t2.Funds_start),
				fmt.Sprintf("%.2f", t2.Earned),
				fmt.Sprintf("%.2f", t2.RS_Eq_value),
				fmt.Sprintf("%.2f", t2.FTE_Eq_value),
				fmt.Sprintf("%.2f", t2.RE_Eq_value),
				fmt.Sprintf("%d", t2.Survivors),
				fmt.Sprintf("%d", t2.Score_Start),
				fmt.Sprintf("%d", t2.Score_End),
				fmt.Sprintf("%d", t2.Consecutiveloss),
				fmt.Sprintf("%d", t2.Consecutiveloss_start),
				fmt.Sprintf("%d", t2.Consecutivewins),
				fmt.Sprintf("%d", t2.Consecutivewins_start),
				fmt.Sprintf("%d", t2.LossBonusLevel),
				fmt.Sprintf("%.2f", t2.Spent),
				game.Team1.Name,
				game.Team1.Strategy,
				game.Team2.Name,
				game.Team2.Strategy,
				game.ID,
			}
			writer.Write(row)
		}
	}
	return writer.Error()
}

// 3. Export each game individually with only relevant columns (minimal info per row)
func ExportGameMinimalCSV(game *engine.Game, path string) error {
	if path == "" {
		path = fmt.Sprintf("%s_minimal.csv", game.ID)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header row
	headers := []string{
		"round_number",
		"is_t1_winner",
		"is_t1_ct",
		"is_ot",
		"t1_score_start",
		"t1_score_end",
		"t1_spent",
		"t1_earned",
		"t1_funds_start",
		"t1_rs_eq",
		"t1_fte_eq",
		"t1_re_eq",
		"t1_survivors",
		"t1_consecutive_losses",
		"t1_consecutive_losses_start",
		"t1_consecutive_wins",
		"t1_consecutive_wins_start",
		"t1_loss_bonus_level",
		"t2_score_start",
		"t2_score_end",
		"t2_spent",
		"t2_earned",
		"t2_funds_start",
		"t2_rs_eq",
		"t2_fte_eq",
		"t2_re_eq",
		"t2_survivors",
		"t2_consecutive_losses",
		"t2_consecutive_losses_start",
		"t2_consecutive_wins",
		"t2_consecutive_wins_start",
		"t2_loss_bonus_level",
	}
	writer.Write(headers)

	for i := range game.Rounds {
		round := &game.Rounds[i]
		t1 := game.Team1.RoundData[i]
		t2 := game.Team2.RoundData[i]
		row := []string{
			fmt.Sprintf("%d", round.RoundNumber),
			fmt.Sprintf("%t", round.IsT1WinnerTeam),
			fmt.Sprintf("%t", round.IsT1CT),
			fmt.Sprintf("%t", game.OT),
			fmt.Sprintf("%d", round.Calc_Outcome.ReasonCode),
			fmt.Sprintf("%t", round.Calc_Outcome.BombPlanted),
			fmt.Sprintf("%d", t1.Score_Start),
			fmt.Sprintf("%d", t1.Score_End),
			fmt.Sprintf("%.2f", t1.Spent),
			fmt.Sprintf("%.2f", t1.Earned),
			fmt.Sprintf("%.2f", t1.Funds_start),
			fmt.Sprintf("%.2f", t1.RS_Eq_value),
			fmt.Sprintf("%.2f", t1.FTE_Eq_value),
			fmt.Sprintf("%.2f", t1.RE_Eq_value),
			fmt.Sprintf("%d", t1.Survivors),
			fmt.Sprintf("%d", t1.Consecutiveloss),
			fmt.Sprintf("%d", t1.Consecutiveloss_start),
			fmt.Sprintf("%d", t1.Consecutivewins),
			fmt.Sprintf("%d", t1.Consecutivewins_start),
			fmt.Sprintf("%d", t1.LossBonusLevel),
			fmt.Sprintf("%d", t2.Score_Start),
			fmt.Sprintf("%d", t2.Score_End),
			fmt.Sprintf("%.2f", t2.Spent),
			fmt.Sprintf("%.2f", t2.Earned),
			fmt.Sprintf("%.2f", t2.Funds_start),
			fmt.Sprintf("%.2f", t2.RS_Eq_value),
			fmt.Sprintf("%.2f", t2.FTE_Eq_value),
			fmt.Sprintf("%.2f", t2.RE_Eq_value),
			fmt.Sprintf("%d", t2.Survivors),
			fmt.Sprintf("%d", t2.Consecutiveloss),
			fmt.Sprintf("%d", t2.Consecutiveloss_start),
			fmt.Sprintf("%d", t2.Consecutivewins),
			fmt.Sprintf("%d", t2.Consecutivewins_start),
			fmt.Sprintf("%d", t2.LossBonusLevel),
		}
		writer.Write(row)
	}
	return writer.Error()
}

// 4. Export all games in one file in minimal format
func ExportAllGamesMinimalCSV(games []*engine.Game, path string) error {
	if path == "" {
		path = "all_games_minimal.csv"
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	// Write header row
	headers := []string{
		"round_number",
		"is_t1_winner",
		"is_t1_ct",
		"is_ot",
		"outcome_reason_code",
		"outcome_bomb_planted",
		"t1_score_start",
		"t1_score_end",
		"t1_spent",
		"t1_earned",
		"t1_funds_start",
		"t1_rs_eq",
		"t1_fte_eq",
		"t1_re_eq",
		"t1_survivors",
		"t1_consecutive_losses",
		"t1_consecutive_losses_start",
		"t1_consecutive_wins",
		"t1_consecutive_wins_start",
		"t1_loss_bonus_level",
		"t2_score_start",
		"t2_score_end",
		"t2_spent",
		"t2_earned",
		"t2_funds_start",
		"t2_rs_eq",
		"t2_fte_eq",
		"t2_re_eq",
		"t2_survivors",
		"t2_consecutive_losses",
		"t2_consecutive_losses_start",
		"t2_consecutive_wins",
		"t2_consecutive_wins_start",
		"t2_loss_bonus_level",
		"game_id",
	}
	writer.Write(headers)

	for _, game := range games {
		for i := range game.Rounds {
			round := &game.Rounds[i]
			t1 := game.Team1.RoundData[i]
			t2 := game.Team2.RoundData[i]
			row := []string{
				fmt.Sprintf("%d", round.RoundNumber),
				fmt.Sprintf("%t", round.IsT1WinnerTeam),
				fmt.Sprintf("%t", round.IsT1CT),
				fmt.Sprintf("%t", round.OT),
				fmt.Sprintf("%d", round.Calc_Outcome.ReasonCode),
				fmt.Sprintf("%t", round.Calc_Outcome.BombPlanted),
				fmt.Sprintf("%d", t1.Score_Start),
				fmt.Sprintf("%d", t1.Score_End),
				fmt.Sprintf("%.2f", t1.Spent),
				fmt.Sprintf("%.2f", t1.Earned),
				fmt.Sprintf("%.2f", t1.Funds_start),
				fmt.Sprintf("%.2f", t1.RS_Eq_value),
				fmt.Sprintf("%.2f", t1.FTE_Eq_value),
				fmt.Sprintf("%.2f", t1.RE_Eq_value),
				fmt.Sprintf("%d", t1.Survivors),
				fmt.Sprintf("%d", t1.Consecutiveloss),
				fmt.Sprintf("%d", t1.Consecutiveloss_start),
				fmt.Sprintf("%d", t1.Consecutivewins),
				fmt.Sprintf("%d", t1.Consecutivewins_start),
				fmt.Sprintf("%d", t1.LossBonusLevel),
				fmt.Sprintf("%d", t2.Score_Start),
				fmt.Sprintf("%d", t2.Score_End),
				fmt.Sprintf("%.2f", t2.Spent),
				fmt.Sprintf("%.2f", t2.Earned),
				fmt.Sprintf("%.2f", t2.Funds_start),
				fmt.Sprintf("%.2f", t2.RS_Eq_value),
				fmt.Sprintf("%.2f", t2.FTE_Eq_value),
				fmt.Sprintf("%.2f", t2.RE_Eq_value),
				fmt.Sprintf("%d", t2.Survivors),
				fmt.Sprintf("%d", t2.Consecutiveloss),
				fmt.Sprintf("%d", t2.Consecutiveloss_start),
				fmt.Sprintf("%d", t2.Consecutivewins),
				fmt.Sprintf("%d", t2.Consecutivewins_start),
				fmt.Sprintf("%d", t2.LossBonusLevel),
				game.ID,
			}
			writer.Write(row)
		}
	}
	return writer.Error()
}
