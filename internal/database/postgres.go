package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// PostgresConfig holds database connection configuration
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// PostgresExporter handles exporting simulation data to PostgreSQL
type PostgresExporter struct {
	db     *sql.DB
	config PostgresConfig
}

// NewPostgresExporter creates a new PostgreSQL exporter with connection
func NewPostgresExporter(config PostgresConfig) (*PostgresExporter, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &PostgresExporter{
		db:     db,
		config: config,
	}, nil
}

// Close closes the database connection
func (pe *PostgresExporter) Close() error {
	if pe.db != nil {
		return pe.db.Close()
	}
	return nil
}

// InitSchema creates all necessary tables if they don't exist
func (pe *PostgresExporter) InitSchema() error {
	ctx := context.Background()

	// Create tables in order of dependencies
	schemas := []string{
		createBatchSimulationsTable,
		createMatchesTable,
		createRoundsTable,
		createBatchSummaryStatsTable,
	}

	for _, schema := range schemas {
		if _, err := pe.db.ExecContext(ctx, schema); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	// Create indexes for better query performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_matches_batch_id ON matches(batch_simulation_id);",
		"CREATE INDEX IF NOT EXISTS idx_matches_timestamp ON matches(timestamp);",
		"CREATE INDEX IF NOT EXISTS idx_rounds_match_id ON rounds(match_id);",
		"CREATE INDEX IF NOT EXISTS idx_rounds_round_number ON rounds(round_number);",
		"CREATE INDEX IF NOT EXISTS idx_rounds_winner ON rounds(winner);",
	}

	for _, index := range indexes {
		if _, err := pe.db.ExecContext(ctx, index); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// Table schemas
const createBatchSimulationsTable = `
CREATE TABLE IF NOT EXISTS batch_simulations (
	id UUID PRIMARY KEY,
	batch_name TEXT NOT NULL,
	description TEXT,
	started_at TIMESTAMP NOT NULL,
	completed_at TIMESTAMP,
	total_experiments INTEGER,
	total_simulations INTEGER,
	status TEXT NOT NULL,
	config JSONB
);`

const createMatchesTable = `
CREATE TABLE IF NOT EXISTS matches (
	id UUID PRIMARY KEY,
	batch_simulation_id UUID REFERENCES batch_simulations(id) ON DELETE CASCADE,
	experiment_name TEXT,
	simulation_number INTEGER NOT NULL,
	team1_name TEXT NOT NULL,
	team2_name TEXT NOT NULL,
	team1_strategy TEXT NOT NULL,
	team2_strategy TEXT NOT NULL,
	team1_final_score INTEGER NOT NULL,
	team2_final_score INTEGER NOT NULL,
	winner TEXT NOT NULL,
	total_rounds INTEGER NOT NULL,
	overtime_rounds INTEGER,
	went_to_overtime BOOLEAN NOT NULL,
	duration_seconds REAL,
	timestamp TIMESTAMP NOT NULL,
	game_rules JSONB
);`

const createRoundsTable = `
CREATE TABLE IF NOT EXISTS rounds (
	id BIGSERIAL PRIMARY KEY,
	match_id UUID NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
	round_number INTEGER NOT NULL,
	half_number INTEGER NOT NULL,
	is_overtime BOOLEAN NOT NULL,
	
	-- Round outcome
	winner TEXT NOT NULL,
	win_reason TEXT,
	
	-- Team 1 economic state
	team1_name TEXT NOT NULL,
	team1_funds_start REAL NOT NULL,
	team1_funds_end REAL NOT NULL,
	team1_equipment_spent REAL NOT NULL,
	team1_equipment_value REAL NOT NULL,
	team1_score_before INTEGER NOT NULL,
	team1_score_after INTEGER NOT NULL,
	team1_consecutive_losses INTEGER NOT NULL,
	team1_surviving_players INTEGER NOT NULL,
	
	-- Team 2 economic state
	team2_name TEXT NOT NULL,
	team2_funds_start REAL NOT NULL,
	team2_funds_end REAL NOT NULL,
	team2_equipment_spent REAL NOT NULL,
	team2_equipment_value REAL NOT NULL,
	team2_score_before INTEGER NOT NULL,
	team2_score_after INTEGER NOT NULL,
	team2_consecutive_losses INTEGER NOT NULL,
	team2_surviving_players INTEGER NOT NULL,
	
	-- Calculated fields
	economic_advantage REAL NOT NULL,
	spending_differential REAL NOT NULL,
	is_pistol_round BOOLEAN NOT NULL,
	is_eco_round BOOLEAN NOT NULL,
	
	-- Metadata
	timestamp TIMESTAMP NOT NULL,
	
	UNIQUE(match_id, round_number)
);`

const createBatchSummaryStatsTable = `
CREATE TABLE IF NOT EXISTS batch_summary_stats (
	id BIGSERIAL PRIMARY KEY,
	batch_simulation_id UUID NOT NULL REFERENCES batch_simulations(id) ON DELETE CASCADE,
	strategy_name TEXT NOT NULL,
	
	-- Win statistics
	total_matches INTEGER NOT NULL,
	total_wins INTEGER NOT NULL,
	total_losses INTEGER NOT NULL,
	win_rate REAL NOT NULL,
	
	-- Round statistics
	total_rounds_played INTEGER NOT NULL,
	rounds_won INTEGER NOT NULL,
	rounds_lost INTEGER NOT NULL,
	round_win_rate REAL NOT NULL,
	
	-- Economic statistics
	avg_funds_per_round REAL,
	avg_spending_per_round REAL,
	avg_equipment_value REAL,
	
	-- Performance metrics
	avg_match_score REAL,
	overtime_frequency REAL,
	avg_match_duration REAL,
	
	-- Advanced metrics
	win_rate_std_dev REAL,
	consistency_score REAL,
	
	timestamp TIMESTAMP NOT NULL,
	
	UNIQUE(batch_simulation_id, strategy_name)
);`
