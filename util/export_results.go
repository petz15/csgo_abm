package util

import (
	"csgo_abm/internal/engine"
	"encoding/json"
	"os"
	"runtime"
)

func ExportResultsToJSON(results interface{}, path string) error {
	if path == "" {
		path = "results.json"
	}

	// Marshal to JSON in memory
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	err = os.WriteFile(path, jsonData, 0644)

	// Clear the large JSON data from memory
	jsonData = nil

	// Hint to the garbage collector
	runtime.GC()

	return err
}

// RoundExport represents detailed data for a single round
type RoundExport struct {
	RoundNumber    int                   `json:"round_number"`
	OvertimeRound  bool                  `json:"overtime_round"`
	Sideswitch     bool                  `json:"sideswitch"`
	Team1Side      string                `json:"team1_side"` // "CT" or "T"
	Team2Side      string                `json:"team2_side"`
	Winner         string                `json:"winner"` // "Team1" or "Team2"
	Team1RoundData engine.Team_RoundData `json:"team1_round_data"`
	Team2RoundData engine.Team_RoundData `json:"team2_round_data"`
	RoundOutcome   engine.RoundOutcome   `json:"round_outcome"`
}

// RoundExportSimple represents minimal round data for compact export
type RoundExportSimple struct {
	RoundNumber         int     `json:"round_number"`
	Winner              string  `json:"winner"`     // "Team1" or "Team2"
	Team1Side           string  `json:"team1_side"` // "CT" or "T"
	Team1ScoreAfter     int     `json:"team1_score_after"`
	Team2ScoreAfter     int     `json:"team2_score_after"`
	Team1Spent          float64 `json:"team1_spent"`
	Team2Spent          float64 `json:"team2_spent"`
	Team1FTEEquipment   float64 `json:"team1_fte_equipment"`
	Team2FTEEquipment   float64 `json:"team2_fte_equipment"`
	Team1Earned         float64 `json:"team1_earned"`
	Team2Earned         float64 `json:"team2_earned"`
	Team1Funds          float64 `json:"team1_funds"`
	Team2Funds          float64 `json:"team2_funds"`
	Team1ConsecLoss     int     `json:"team1_consecutive_loss"`
	Team2ConsecLoss     int     `json:"team2_consecutive_loss"`
	Team1LossBonusLevel int     `json:"team1_loss_bonus_level"`
	Team2LossBonusLevel int     `json:"team2_loss_bonus_level"`
}

// GameRoundsExport represents all rounds for a complete game
type GameRoundsExport struct {
	GameID         string        `json:"game_id"`
	Team1Name      string        `json:"team1_name"`
	Team1Strategy  string        `json:"team1_strategy"`
	Team2Name      string        `json:"team2_name"`
	Team2Strategy  string        `json:"team2_strategy"`
	FinalScore     [2]int        `json:"final_score"` // [Team1Score, Team2Score]
	WentToOvertime bool          `json:"went_to_overtime"`
	TotalRounds    int           `json:"total_rounds"`
	Rounds         []RoundExport `json:"rounds"`
}

// GameRoundsExportSimple represents all rounds with minimal data
type GameRoundsExportSimple struct {
	GameID         string              `json:"game_id"`
	Team1Name      string              `json:"team1_name"`
	Team1Strategy  string              `json:"team1_strategy"`
	Team2Name      string              `json:"team2_name"`
	Team2Strategy  string              `json:"team2_strategy"`
	FinalScore     [2]int              `json:"final_score"` // [Team1Score, Team2Score]
	WentToOvertime bool                `json:"went_to_overtime"`
	TotalRounds    int                 `json:"total_rounds"`
	Rounds         []RoundExportSimple `json:"rounds"`
}

// ExportRoundsToJSON exports detailed round-by-round data for a game
func ExportRoundsToJSON(game *engine.Game, path string) error {
	if path == "" {
		path = "rounds.json"
	}

	// Build the export structure
	export := GameRoundsExport{
		GameID:         game.ID,
		Team1Name:      game.Team1.Name,
		Team1Strategy:  game.Team1.Strategy,
		Team2Name:      game.Team2.Name,
		Team2Strategy:  game.Team2.Strategy,
		FinalScore:     game.Score,
		WentToOvertime: game.OT,
		TotalRounds:    len(game.Rounds),
		Rounds:         make([]RoundExport, 0, len(game.Rounds)),
	}

	// Process each round
	for i := range game.Rounds {
		round := &game.Rounds[i]

		// Determine sides using helper methods
		team1Side := "T"
		team2Side := "CT"
		if round.IsT1CT {
			team1Side = "CT"
			team2Side = "T"
		}

		// Determine winner
		winner := "Team2"
		if round.IsT1WinnerTeam {
			winner = "Team1"
		}

		// Extract team data for this round (round number is 1-indexed)
		team1Data := game.Team1.RoundData[i]
		team2Data := game.Team2.RoundData[i]

		roundExport := RoundExport{
			RoundNumber:    round.RoundNumber,
			OvertimeRound:  round.OT,
			Sideswitch:     round.Sideswitch,
			Team1Side:      team1Side,
			Team2Side:      team2Side,
			Winner:         winner,
			Team1RoundData: team1Data,
			Team2RoundData: team2Data,
			RoundOutcome:   round.Calc_Outcome,
		}

		export.Rounds = append(export.Rounds, roundExport)
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	err = os.WriteFile(path, jsonData, 0644)

	// Clear from memory
	jsonData = nil

	// Hint to GC
	runtime.GC()

	return err
}

// ExportRoundsToJSONSimple exports minimal round-by-round data for a game
func ExportRoundsToJSONSimple(game *engine.Game, path string) error {
	if path == "" {
		path = "rounds_simple.json"
	}

	// Build the export structure
	export := GameRoundsExportSimple{
		GameID:         game.ID,
		Team1Name:      game.Team1.Name,
		Team1Strategy:  game.Team1.Strategy,
		Team2Name:      game.Team2.Name,
		Team2Strategy:  game.Team2.Strategy,
		FinalScore:     game.Score,
		WentToOvertime: game.OT,
		TotalRounds:    len(game.Rounds),
		Rounds:         make([]RoundExportSimple, 0, len(game.Rounds)),
	}

	// Process each round
	for i := range game.Rounds {
		round := &game.Rounds[i]

		// Determine sides
		team1Side := "T"
		if round.IsT1CT {
			team1Side = "CT"
		}

		// Determine winner
		winner := "Team2"
		if round.IsT1WinnerTeam {
			winner = "Team1"
		}

		// Extract team data for this round
		team1Data := game.Team1.RoundData[i]
		team2Data := game.Team2.RoundData[i]

		roundExport := RoundExportSimple{
			RoundNumber:         round.RoundNumber,
			Winner:              winner,
			Team1Side:           team1Side,
			Team1ScoreAfter:     team1Data.Score_End,
			Team2ScoreAfter:     team2Data.Score_End,
			Team1Spent:          team1Data.Spent,
			Team2Spent:          team2Data.Spent,
			Team1FTEEquipment:   team1Data.FTE_Eq_value,
			Team2FTEEquipment:   team2Data.FTE_Eq_value,
			Team1Earned:         team1Data.Earned,
			Team2Earned:         team2Data.Earned,
			Team1Funds:          team1Data.Funds,
			Team2Funds:          team2Data.Funds,
			Team1ConsecLoss:     team1Data.Consecutiveloss,
			Team2ConsecLoss:     team2Data.Consecutiveloss,
			Team1LossBonusLevel: team1Data.LossBonusLevel,
			Team2LossBonusLevel: team2Data.LossBonusLevel,
		}

		export.Rounds = append(export.Rounds, roundExport)
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	err = os.WriteFile(path, jsonData, 0644)

	// Clear from memory
	jsonData = nil

	// Hint to GC
	runtime.GC()

	return err
}
