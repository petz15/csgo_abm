package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// BatchSimulationData represents a batch simulation run
type BatchSimulationData struct {
	ID               string
	BatchName        string
	Description      string
	StartedAt        time.Time
	CompletedAt      *time.Time
	TotalExperiments int
	TotalSimulations int
	Status           string
	ConfigJSON       string
}

// MatchData represents a single match/game
type MatchData struct {
	ID                string
	BatchSimulationID string
	ExperimentName    string
	SimulationNumber  int
	Team1Name         string
	Team2Name         string
	Team1Strategy     string
	Team2Strategy     string
	Team1FinalScore   int
	Team2FinalScore   int
	Winner            string
	TotalRounds       int
	OvertimeRounds    int
	WentToOvertime    bool
	DurationSeconds   float64
	Timestamp         time.Time
	GameRulesJSON     string
}

// RoundData represents a single round within a match
type RoundData struct {
	MatchID                string
	RoundNumber            int
	HalfNumber             int
	IsOvertime             bool
	Winner                 string
	WinReason              string
	Team1Name              string
	Team1FundsStart        float64
	Team1FundsEnd          float64
	Team1EquipmentSpent    float64
	Team1EquipmentValue    float64
	Team1ScoreBefore       int
	Team1ScoreAfter        int
	Team1ConsecutiveLosses int
	Team1SurvivingPlayers  int
	Team2Name              string
	Team2FundsStart        float64
	Team2FundsEnd          float64
	Team2EquipmentSpent    float64
	Team2EquipmentValue    float64
	Team2ScoreBefore       int
	Team2ScoreAfter        int
	Team2ConsecutiveLosses int
	Team2SurvivingPlayers  int
	EconomicAdvantage      float64
	SpendingDifferential   float64
	IsPistolRound          bool
	IsEcoRound             bool
	Timestamp              time.Time
}

// BatchSummaryStats represents aggregated statistics for a strategy in a batch
type BatchSummaryStats struct {
	BatchSimulationID   string
	StrategyName        string
	TotalMatches        int
	TotalWins           int
	TotalLosses         int
	WinRate             float64
	TotalRoundsPlayed   int
	RoundsWon           int
	RoundsLost          int
	RoundWinRate        float64
	AvgFundsPerRound    float64
	AvgSpendingPerRound float64
	AvgEquipmentValue   float64
	AvgMatchScore       float64
	OvertimeFrequency   float64
	AvgMatchDuration    float64
	WinRateStdDev       float64
	ConsistencyScore    float64
	Timestamp           time.Time
}

// ExportBatchSimulation exports batch-level metadata
func (pe *PostgresExporter) ExportBatchSimulation(ctx context.Context, batch BatchSimulationData) error {
	query := `
		INSERT INTO batch_simulations (
			id, batch_name, description, started_at, completed_at,
			total_experiments, total_simulations, status, config
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO UPDATE SET
			completed_at = EXCLUDED.completed_at,
			total_experiments = EXCLUDED.total_experiments,
			total_simulations = EXCLUDED.total_simulations,
			status = EXCLUDED.status,
			config = EXCLUDED.config
	`

	_, err := pe.db.ExecContext(ctx, query,
		batch.ID,
		batch.BatchName,
		batch.Description,
		batch.StartedAt,
		batch.CompletedAt,
		batch.TotalExperiments,
		batch.TotalSimulations,
		batch.Status,
		batch.ConfigJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to export batch simulation: %w", err)
	}

	return nil
}

// ExportMatch exports a single match with all its metadata
func (pe *PostgresExporter) ExportMatch(ctx context.Context, match MatchData) error {
	query := `
		INSERT INTO matches (
			id, batch_simulation_id, experiment_name, simulation_number,
			team1_name, team2_name, team1_strategy, team2_strategy,
			team1_final_score, team2_final_score, winner, total_rounds,
			overtime_rounds, went_to_overtime, duration_seconds, timestamp, game_rules
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (id) DO NOTHING
	`

	_, err := pe.db.ExecContext(ctx, query,
		match.ID,
		match.BatchSimulationID,
		match.ExperimentName,
		match.SimulationNumber,
		match.Team1Name,
		match.Team2Name,
		match.Team1Strategy,
		match.Team2Strategy,
		match.Team1FinalScore,
		match.Team2FinalScore,
		match.Winner,
		match.TotalRounds,
		match.OvertimeRounds,
		match.WentToOvertime,
		match.DurationSeconds,
		match.Timestamp,
		match.GameRulesJSON,
	)

	if err != nil {
		return fmt.Errorf("failed to export match: %w", err)
	}

	return nil
}

// ExportRound exports a single round with detailed economic data
func (pe *PostgresExporter) ExportRound(ctx context.Context, round RoundData) error {
	query := `
		INSERT INTO rounds (
			match_id, round_number, half_number, is_overtime, winner, win_reason,
			team1_name, team1_funds_start, team1_funds_end, team1_equipment_spent,
			team1_equipment_value, team1_score_before, team1_score_after,
			team1_consecutive_losses, team1_surviving_players,
			team2_name, team2_funds_start, team2_funds_end, team2_equipment_spent,
			team2_equipment_value, team2_score_before, team2_score_after,
			team2_consecutive_losses, team2_surviving_players,
			economic_advantage, spending_differential, is_pistol_round, is_eco_round, timestamp
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24,
			$25, $26, $27, $28, $29
		)
		ON CONFLICT (match_id, round_number) DO NOTHING
	`

	_, err := pe.db.ExecContext(ctx, query,
		round.MatchID,
		round.RoundNumber,
		round.HalfNumber,
		round.IsOvertime,
		round.Winner,
		round.WinReason,
		round.Team1Name,
		round.Team1FundsStart,
		round.Team1FundsEnd,
		round.Team1EquipmentSpent,
		round.Team1EquipmentValue,
		round.Team1ScoreBefore,
		round.Team1ScoreAfter,
		round.Team1ConsecutiveLosses,
		round.Team1SurvivingPlayers,
		round.Team2Name,
		round.Team2FundsStart,
		round.Team2FundsEnd,
		round.Team2EquipmentSpent,
		round.Team2EquipmentValue,
		round.Team2ScoreBefore,
		round.Team2ScoreAfter,
		round.Team2ConsecutiveLosses,
		round.Team2SurvivingPlayers,
		round.EconomicAdvantage,
		round.SpendingDifferential,
		round.IsPistolRound,
		round.IsEcoRound,
		round.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to export round: %w", err)
	}

	return nil
}

// ExportRoundsBatch exports multiple rounds in a single transaction for better performance
func (pe *PostgresExporter) ExportRoundsBatch(ctx context.Context, rounds []RoundData) error {
	if len(rounds) == 0 {
		return nil
	}

	tx, err := pe.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO rounds (
			match_id, round_number, half_number, is_overtime, winner, win_reason,
			team1_name, team1_funds_start, team1_funds_end, team1_equipment_spent,
			team1_equipment_value, team1_score_before, team1_score_after,
			team1_consecutive_losses, team1_surviving_players,
			team2_name, team2_funds_start, team2_funds_end, team2_equipment_spent,
			team2_equipment_value, team2_score_before, team2_score_after,
			team2_consecutive_losses, team2_surviving_players,
			economic_advantage, spending_differential, is_pistol_round, is_eco_round, timestamp
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24,
			$25, $26, $27, $28, $29
		)
		ON CONFLICT (match_id, round_number) DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, round := range rounds {
		_, err := stmt.ExecContext(ctx,
			round.MatchID,
			round.RoundNumber,
			round.HalfNumber,
			round.IsOvertime,
			round.Winner,
			round.WinReason,
			round.Team1Name,
			round.Team1FundsStart,
			round.Team1FundsEnd,
			round.Team1EquipmentSpent,
			round.Team1EquipmentValue,
			round.Team1ScoreBefore,
			round.Team1ScoreAfter,
			round.Team1ConsecutiveLosses,
			round.Team1SurvivingPlayers,
			round.Team2Name,
			round.Team2FundsStart,
			round.Team2FundsEnd,
			round.Team2EquipmentSpent,
			round.Team2EquipmentValue,
			round.Team2ScoreBefore,
			round.Team2ScoreAfter,
			round.Team2ConsecutiveLosses,
			round.Team2SurvivingPlayers,
			round.EconomicAdvantage,
			round.SpendingDifferential,
			round.IsPistolRound,
			round.IsEcoRound,
			round.Timestamp,
		)
		if err != nil {
			return fmt.Errorf("failed to insert round %d: %w", round.RoundNumber, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ExportBatchSummaryStats exports aggregated statistics for a batch simulation
func (pe *PostgresExporter) ExportBatchSummaryStats(ctx context.Context, stats BatchSummaryStats) error {
	query := `
		INSERT INTO batch_summary_stats (
			batch_simulation_id, strategy_name, total_matches, total_wins, total_losses,
			win_rate, total_rounds_played, rounds_won, rounds_lost, round_win_rate,
			avg_funds_per_round, avg_spending_per_round, avg_equipment_value,
			avg_match_score, overtime_frequency, avg_match_duration,
			win_rate_std_dev, consistency_score, timestamp
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19
		)
		ON CONFLICT (batch_simulation_id, strategy_name) DO UPDATE SET
			total_matches = EXCLUDED.total_matches,
			total_wins = EXCLUDED.total_wins,
			total_losses = EXCLUDED.total_losses,
			win_rate = EXCLUDED.win_rate,
			total_rounds_played = EXCLUDED.total_rounds_played,
			rounds_won = EXCLUDED.rounds_won,
			rounds_lost = EXCLUDED.rounds_lost,
			round_win_rate = EXCLUDED.round_win_rate,
			avg_funds_per_round = EXCLUDED.avg_funds_per_round,
			avg_spending_per_round = EXCLUDED.avg_spending_per_round,
			avg_equipment_value = EXCLUDED.avg_equipment_value,
			avg_match_score = EXCLUDED.avg_match_score,
			overtime_frequency = EXCLUDED.overtime_frequency,
			avg_match_duration = EXCLUDED.avg_match_duration,
			win_rate_std_dev = EXCLUDED.win_rate_std_dev,
			consistency_score = EXCLUDED.consistency_score,
			timestamp = EXCLUDED.timestamp
	`

	_, err := pe.db.ExecContext(ctx, query,
		stats.BatchSimulationID,
		stats.StrategyName,
		stats.TotalMatches,
		stats.TotalWins,
		stats.TotalLosses,
		stats.WinRate,
		stats.TotalRoundsPlayed,
		stats.RoundsWon,
		stats.RoundsLost,
		stats.RoundWinRate,
		stats.AvgFundsPerRound,
		stats.AvgSpendingPerRound,
		stats.AvgEquipmentValue,
		stats.AvgMatchScore,
		stats.OvertimeFrequency,
		stats.AvgMatchDuration,
		stats.WinRateStdDev,
		stats.ConsistencyScore,
		stats.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to export batch summary stats: %w", err)
	}

	return nil
}

// ExportMatchWithRounds exports a complete match including all rounds in a single transaction
func (pe *PostgresExporter) ExportMatchWithRounds(ctx context.Context, match MatchData, rounds []RoundData) error {
	tx, err := pe.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert match
	matchQuery := `
		INSERT INTO matches (
			id, batch_simulation_id, experiment_name, simulation_number,
			team1_name, team2_name, team1_strategy, team2_strategy,
			team1_final_score, team2_final_score, winner, total_rounds,
			overtime_rounds, went_to_overtime, duration_seconds, timestamp, game_rules
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (id) DO NOTHING
	`

	_, err = tx.ExecContext(ctx, matchQuery,
		match.ID,
		match.BatchSimulationID,
		match.ExperimentName,
		match.SimulationNumber,
		match.Team1Name,
		match.Team2Name,
		match.Team1Strategy,
		match.Team2Strategy,
		match.Team1FinalScore,
		match.Team2FinalScore,
		match.Winner,
		match.TotalRounds,
		match.OvertimeRounds,
		match.WentToOvertime,
		match.DurationSeconds,
		match.Timestamp,
		match.GameRulesJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to insert match: %w", err)
	}

	// Insert all rounds
	if len(rounds) > 0 {
		roundStmt, err := tx.PrepareContext(ctx, `
			INSERT INTO rounds (
				match_id, round_number, half_number, is_overtime, winner, win_reason,
				team1_name, team1_funds_start, team1_funds_end, team1_equipment_spent,
				team1_equipment_value, team1_score_before, team1_score_after,
				team1_consecutive_losses, team1_surviving_players,
				team2_name, team2_funds_start, team2_funds_end, team2_equipment_spent,
				team2_equipment_value, team2_score_before, team2_score_after,
				team2_consecutive_losses, team2_surviving_players,
				economic_advantage, spending_differential, is_pistol_round, is_eco_round, timestamp
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10, $11, $12, $13, $14, $15,
				$16, $17, $18, $19, $20, $21, $22, $23, $24,
				$25, $26, $27, $28, $29
			)
			ON CONFLICT (match_id, round_number) DO NOTHING
		`)
		if err != nil {
			return fmt.Errorf("failed to prepare round statement: %w", err)
		}
		defer roundStmt.Close()

		for _, round := range rounds {
			_, err := roundStmt.ExecContext(ctx,
				round.MatchID,
				round.RoundNumber,
				round.HalfNumber,
				round.IsOvertime,
				round.Winner,
				round.WinReason,
				round.Team1Name,
				round.Team1FundsStart,
				round.Team1FundsEnd,
				round.Team1EquipmentSpent,
				round.Team1EquipmentValue,
				round.Team1ScoreBefore,
				round.Team1ScoreAfter,
				round.Team1ConsecutiveLosses,
				round.Team1SurvivingPlayers,
				round.Team2Name,
				round.Team2FundsStart,
				round.Team2FundsEnd,
				round.Team2EquipmentSpent,
				round.Team2EquipmentValue,
				round.Team2ScoreBefore,
				round.Team2ScoreAfter,
				round.Team2ConsecutiveLosses,
				round.Team2SurvivingPlayers,
				round.EconomicAdvantage,
				round.SpendingDifferential,
				round.IsPistolRound,
				round.IsEcoRound,
				round.Timestamp,
			)
			if err != nil {
				return fmt.Errorf("failed to insert round %d: %w", round.RoundNumber, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GenerateUUID is a helper function to generate UUIDs for export
func GenerateUUID() string {
	return uuid.New().String()
}
