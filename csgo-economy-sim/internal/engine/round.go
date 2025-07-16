package engine

type Round struct {
	CTTeam      *Team
	TTeam       *Team
	RoundNumber int
	WinnerTeam  *Team
	WinnerSide  bool // True if CT wins, false if T wins
	BombPlanted bool
	SurvivingCT int
	SurvivingT  int
}

func initRound(ctTeam, tTeam *Team, roundNumber int) *Round {
	return &Round{
		CTTeam:      ctTeam,
		TTeam:       tTeam,
		RoundNumber: roundNumber,
	}
}

func (r *Round) determineWinner() *Team {
	// Placeholder for CSF logic
	ctSpend := r.CTTeam.GetTotalInvestment()
	tSpend := r.TTeam.GetTotalInvestment()

	if ctSpend > tSpend {
		return r.CTTeam
	} else {
		return r.TTeam
	}
}

func (r *Round) Reset() {
	r.BombPlanted = false
	r.SurvivingCT = 0
	r.SurvivingT = 0
	r.Winner = nil
}
