package models

// StrategyContext holds all relevant information for economic decision making
type StrategyContext_simple struct {
	Funds             float64
	Equipment         float64
	PlayersAlive      int     // Number of players alive in the team
	RoundImportance   float64 // Importance of the current round
	EconomicAdvantage float64 // Economic advantage over the opponent
	WinProbability    float64 // Probability of winning the round based on current state
	CurrentRound      int
	OpponentScore     int
	OwnScore          int
	ConsecutiveLosses int
	Side              bool // true = CT, false = T
	IsPistolRound     bool
	IsLastRoundHalf   bool
	IsOvertime        bool
	IsEcoAfterPistol  bool
	HalfLength        int // Length of a half in rounds
	OTHalfLength      int // Length of overtime half in rounds
}

func empty() {
	// This function is intentionally left empty to ensure the package compiles correctly.
	// It serves as a placeholder for future enhancements or additional logic.

}
