package engine

type Round struct {
	RoundNumber    int
	Team1          *Team
	Team2          *Team
	CTTeam         bool // false if Team1 is CT
	WinnerTeam     bool //False if Team1 wins, true if Team2 wins
	WinnerSide     bool // false if CT wins, true if T wins
	BombPlanted    bool
	RoundEndReason string // "Elimination", "Defused", "Exploded"
	Sideswitch     bool   // True if sideswitch has occurred
	Bombplanted    bool
	gameRules      *GameRules
	SurvivingT1    int     // Number of surviving T1 players
	SurvivingT2    int     // Number of surviving T2 players
	EquipmentT1    float64 // Total equipment value for T1 team
	EquipmentT2    float64 // Total equipment value for T2 team
	FundsearnedT1  float64 // Total funds earned by T1 team
	FundsearnedT2  float64 // Total funds earned by T2 team
	OT             bool    // True if this is an overtime round
}

func NewRound(T1 *Team, T2 *Team, roundNumber int, ctteam bool, sidewitch bool, gamerules *GameRules, ot bool) *Round {

	T1.NewRound()
	T2.NewRound()

	return &Round{
		Team1:         T1,
		Team2:         T2,
		CTTeam:        ctteam,
		RoundNumber:   roundNumber,
		WinnerTeam:    false,
		WinnerSide:    false,
		BombPlanted:   false,
		Sideswitch:    sidewitch,
		Bombplanted:   false,
		gameRules:     gamerules,
		SurvivingT1:   0,
		SurvivingT2:   0,
		EquipmentT1:   0,
		EquipmentT2:   0,
		FundsearnedT1: 0,
		FundsearnedT2: 0,
		OT:            ot,
	}
}

// solution for what information gets passed to the teams, could be a json file in the gamerules which specifies
// which variables are passed to the teams. No clue how this will be done yet, but it is an idea.
func (r *Round) BuyPhase() {
	r.Team1.BuyPhase(r.RoundNumber, r.OT, r.Team2, *r.gameRules) // Call the team's buy phase logic
	r.Team2.BuyPhase(r.RoundNumber, r.OT, r.Team1, *r.gameRules)

}

func (r *Round) RoundStart() {
	// Simulate the round using the comprehensive DetermineRoundOutcome function
	r.determineRoundOutcome()

	//round end logic
	r.RoundEnd()
}

func (r *Round) RoundEnd() {
	//determine funds earned
	r.determineFundsEarned()

	//update Teams
	r.Team1.RoundEnd(!r.WinnerTeam, r.FundsearnedT1, r.SurvivingT1, r.EquipmentT1)
	r.Team2.RoundEnd(r.WinnerTeam, r.FundsearnedT2, r.SurvivingT2, r.EquipmentT2)
}

func (r *Round) determineFundsEarned() {

	// Determine which team is CT and which is T
	winnerFunds := 3250 * 5
	loserBonusFunds := 0

	if r.BombPlanted {
		winnerFunds = 3500 * 5
	}

	// Funds earned by killing players
	r.FundsearnedT1 = float64((5 - r.SurvivingT2) * 300) // 300 per player killed
	r.FundsearnedT2 = float64((5 - r.SurvivingT1) * 300) // 300 per player killed

	if r.BombPlanted && r.WinnerSide {
		winnerFunds += 300 // Add bonus for winning team
	} else if r.BombPlanted && !r.WinnerSide {
		winnerFunds += 300    // Add bonus for winning team
		loserBonusFunds = 300 // Add bonus for losing team
	}

	if !r.WinnerTeam {
		r.FundsearnedT1 += float64(winnerFunds)
		loserBonusFunds += r.LossBonusCalculation(r.Team2) // Calculate loss bonus for T2
		r.FundsearnedT2 += float64(loserBonusFunds * 5)
	} else {
		r.FundsearnedT2 += float64(winnerFunds)
		loserBonusFunds += r.LossBonusCalculation(r.Team1) // Calculate loss bonus for T1
		r.FundsearnedT1 += float64(loserBonusFunds * 5)
	}

	// Ensure funds do not exceed maximum allowed
	if r.Team1.Funds+r.FundsearnedT1 > r.gameRules.MaxFunds {
		r.FundsearnedT1 = r.gameRules.MaxFunds - r.Team1.Funds
	}
	if r.Team2.Funds+r.FundsearnedT2 > r.gameRules.MaxFunds {
		r.FundsearnedT2 = r.gameRules.MaxFunds - r.Team2.Funds
	}

}

func (r *Round) LossBonusCalculation(loserteam *Team) int {
	// Calculate loss bonus based on consecutive losses
	lossBonus := 0
	if loserteam.Consecutiveloss >= 4 {
		lossBonus = 3400 //Loss bonus for 5th loss and beyond
	} else if loserteam.Consecutiveloss == 3 {
		lossBonus = 2900 // Loss bonus for fourth loss
	} else if loserteam.Consecutiveloss == 2 {
		lossBonus = 2400 // Loss bonus for third loss
	} else if loserteam.Consecutiveloss == 1 {
		lossBonus = 1900 // Loss bonus for second loss
	} else {
		lossBonus = 1400 // Loss bonus for first loss
	}
	return lossBonus
}

// determineRoundOutcome uses the comprehensive ABM-based probability function
// to determine all aspects of the round outcome in one call
func (r *Round) determineRoundOutcome() {
	// Get equipment values for CSF calculation
	team1equipment := r.Team1.Equipment
	team2equipment := r.Team2.Equipment

	// Calculate CSF probability - we need CT win probability for the distributions
	var ctWinProb float64
	if r.CTTeam {
		// Team1 is CT, so use team1's win probability
		ctWinProb = ContestSuccessFunction_simples(team1equipment, team2equipment)
	} else {
		// Team1 is T, so CT win probability is 1 - team1's win probability
		ctWinProb = 1.0 - ContestSuccessFunction_simples(team1equipment, team2equipment)
	}

	// Get comprehensive round outcome from ABM distributions (uses CT win probability)
	outcome := DetermineRoundOutcome(ctWinProb)

	// Determine which team won
	r.WinnerTeam = !outcome.CTWins
	if r.CTTeam {
		// If Team1 is CT, they win when CT wins
		r.WinnerTeam = !outcome.CTWins
	}

	// Set winner side (false = CT, true = T)
	r.WinnerSide = !outcome.CTWins

	// Set bomb planted status
	r.BombPlanted = outcome.BombPlanted
	r.Bombplanted = outcome.BombPlanted

	// Set round end reason
	r.RoundEndReason = outcome.ReasonName

	// Assign survivors to correct teams based on side assignment
	if !r.CTTeam {
		// Team1 is CT, Team2 is T
		r.SurvivingT1 = outcome.CTSurvivors
		r.SurvivingT2 = outcome.TSurvivors
		r.EquipmentT1 = outcome.CTEquipmentSaved
		r.EquipmentT2 = outcome.TEquipmentSaved
	} else {
		// Team1 is T, Team2 is CT
		r.SurvivingT1 = outcome.TSurvivors
		r.SurvivingT2 = outcome.CTSurvivors
		r.EquipmentT1 = outcome.TEquipmentSaved
		r.EquipmentT2 = outcome.CTEquipmentSaved
	}
}
