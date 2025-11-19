package strategy

// StrategyContext holds all relevant information for economic decision making
type StrategyContext_simple struct {
	Funds             float64
	Equipment         float64
	PlayersAlive      int     // Number of players alive in the team
	RoundImportance   float64 // Importance of the current round
	CurrentRound      int
	OpponentScore     int
	OwnScore          int
	ConsecutiveLosses int
	ConsecutiveWins   int
	Side              bool // true = CT, false = T
	IsPistolRound     bool
	IsLastRoundHalf   bool
	IsOvertime        bool
	IsAfterPistol     bool
	HalfLength        int  // Length of a half in rounds
	OTHalfLength      int  // Length of overtime half in rounds
	OwnSurvivors      int  // Number of survivors in the previous RoundEndReason
	EnemySurvivors    int  // Number of enemy survivors in the previous RoundEndReason
	RoundEndReason    int  // Reason for the end of the last round
	Is_BombPlanted    bool // Whether the bomb was planted in the last round
	Max_Funds         float64
	DefaultEquipment  float64
	OTFunds           float64
	OTEquipment       float64
}

func empty() {
	// This function is intentionally left empty to ensure the package compiles correctly.
	// It serves as a placeholder for future enhancements or additional logic.

}

// ReLU activation function
func relu(x float64) float64 {
	if x > 0 {
		return x
	}
	return 0
}
