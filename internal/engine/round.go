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
	// Simulate the round logic here
	// For now, we will just randomly determine a winner
	r.determineWinner()

	// Determine if bomb was planted
	r.determineBombplant()

	//round end logic
	r.RoundEnd()
}

func (r *Round) RoundEnd() {
	//TODO: Create better probabilities. e.g. if no bomb plant and T wins, no CT can survive

	//determine surviving players
	r.determineSurvivors()

	//determine remaining equipment
	r.determineRemainingEquipment()

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

func (r *Round) determineSurvivors() {

	//this is necessary if bomb plant chances are used
	ctteam_equipment := r.Team1.Equipment
	tteam_equipment := r.Team2.Equipment

	if r.CTTeam {
		ctteam_equipment = r.Team2.Equipment
		tteam_equipment = r.Team1.Equipment
	}

	// Map round end reason to code for ABM lookup
	reasonCode := r.mapRoundEndReasonToCode()

	// Calculate CSF probability for winner
	var csfProb float64
	if r.WinnerSide { // T wins
		csfProb = ContestSuccessFunction_simples(tteam_equipment, ctteam_equipment, r.gameRules.CSF_r)
	} else { // CT wins
		csfProb = ContestSuccessFunction_simples(ctteam_equipment, tteam_equipment, r.gameRules.CSF_r)
	}

	// Use ABM models to sample survivors
	surviving_CT := 0
	surviving_T := 0

	if r.WinnerSide { // T wins
		surviving_T = SampleSurvivorsFromABM("T", reasonCode, csfProb)
		// Losing CT team has fewer/no survivors based on reason
		if reasonCode == "8" { // Elimination
			surviving_CT = 0
		} else {
			// For bomb explosion scenarios, CT may have some survivors
			surviving_CT = SampleSurvivorsFromABM("CT", reasonCode, 1.0-csfProb)
		}
	} else { // CT wins
		surviving_CT = SampleSurvivorsFromABM("CT", reasonCode, csfProb)
		// Losing T team
		if reasonCode == "8" { // Elimination
			surviving_T = 0
		} else {
			surviving_T = SampleSurvivorsFromABM("T", reasonCode, 1.0-csfProb)
		}
	}

	if !r.CTTeam {
		r.SurvivingT1 = surviving_CT
		r.SurvivingT2 = surviving_T
	} else {
		r.SurvivingT1 = surviving_T
		r.SurvivingT2 = surviving_CT
	}

}

// mapRoundEndReasonToCode converts the round end reason string to the code used in ABM models
func (r *Round) mapRoundEndReasonToCode() string {
	switch r.RoundEndReason {
	case "Defused":
		return "7"
	case "Exploded":
		return "9"
	case "Elimination":
		return "8"
	default:
		return "8" // Default to elimination
	}
}

func (r *Round) determineRemainingEquipment() {
	// Use ABM models to determine equipment saved based on survivors
	ctteam_equipment := r.Team1.Equipment
	tteam_equipment := r.Team2.Equipment

	ctSurvivors := r.SurvivingT1
	tSurvivors := r.SurvivingT2

	if r.CTTeam {
		ctteam_equipment = r.Team2.Equipment
		tteam_equipment = r.Team1.Equipment
		ctSurvivors = r.SurvivingT2
		tSurvivors = r.SurvivingT1
	}

	reasonCode := r.mapRoundEndReasonToCode()

	// Sample equipment saved from ABM models
	equipmentCT := SampleEquipmentSavedFromABM("CT", reasonCode, ctSurvivors)
	equipmentT := SampleEquipmentSavedFromABM("T", reasonCode, tSurvivors)

	// If ABM returns 0, use fallback calculation based on average equipment per player
	if equipmentCT == 0 && ctSurvivors > 0 {
		equipmentCT = (ctteam_equipment / 5) * float64(ctSurvivors) * 0.7 // 70% of average equipment
	}
	if equipmentT == 0 && tSurvivors > 0 {
		equipmentT = (tteam_equipment / 5) * float64(tSurvivors) * 0.7
	}

	if !r.CTTeam {
		r.EquipmentT1 = equipmentCT
		r.EquipmentT2 = equipmentT
	} else {
		r.EquipmentT1 = equipmentT
		r.EquipmentT2 = equipmentCT
	}
}

func (r *Round) determineBombplant() {
	// Placeholder for CSF logic

	ctteam_equipment := r.Team1.Equipment
	tteam_equipment := r.Team2.Equipment

	if r.CTTeam {
		ctteam_equipment = r.Team2.Equipment
		tteam_equipment = r.Team1.Equipment
	}

	r.BombPlanted = bool_CSF_simple(tteam_equipment, ctteam_equipment, r.gameRules.CSF_r)

}

func (r *Round) determineWinner() {
	// Use CSF with r = 1.0855 to determine winner
	team1equipment := r.Team1.Equipment
	team2equipment := r.Team2.Equipment

	// Determine if Team1 wins using CSF
	r.WinnerTeam = bool_CSF_simple(team1equipment, team2equipment, r.gameRules.CSF_r)

	// Determine which side won (CT or T)
	ctWins := false
	if r.WinnerTeam == r.CTTeam {
		r.WinnerSide = false // CT wins
		ctWins = true
	} else {
		r.WinnerSide = true // T wins
		ctWins = false
	}

	// Calculate CSF probability for ABM sampling
	csfProb := ContestSuccessFunction_simples(team1equipment, team2equipment, r.gameRules.CSF_r)
	if !ctWins {
		csfProb = 1.0 - csfProb
	}

	// Determine the round end reason using ABM distribution
	r.RoundEndReason = SampleRoundEndFromABM(ctWins, csfProb)
}
