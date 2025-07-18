package engine

import (
	"math"
)

type Round struct {
	RoundNumber   int
	Team1         *Team
	Team2         *Team
	CTTeam        bool // false if Team1 is CT
	WinnerTeam    bool //False if Team1 wins, true if Team2 wins
	WinnerSide    bool // false if CT wins, true if T wins
	BombPlanted   bool
	Sideswitch    bool // True if sideswitch has occurred
	Bombplanted   bool
	gameRules     *GameRules
	SurvivingT1   int     // Number of surviving T1 players
	SurvivingT2   int     // Number of surviving T2 players
	EquipmentT1   float64 // Total equipment value for T1 team
	EquipmentT2   float64 // Total equipment value for T2 team
	FundsearnedT1 float64 // Total funds earned by T1 team
	FundsearnedT2 float64 // Total funds earned by T2 team
}

func NewRound(T1 *Team, T2 *Team, roundNumber int, ctteam bool, sidewitch bool, gamerules *GameRules) *Round {

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
	}
}

// solution for what information gets passed to the teams, could be a json file in the gamerules which specifies
// which variables are passed to the teams. No clue how this will be done yet, but it is an idea.
func (r *Round) BuyPhase() {
	r.Team1.BuyPhase(r.RoundNumber, r.Team2.Score)
	r.Team2.BuyPhase(r.RoundNumber, r.Team1.Score)
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
	/* ctteam_equipment := r.Team1.Equipment
	tteam_equipment := r.Team2.Equipment

	if r.CTTeam {
		ctteam_equipment = r.Team1.Equipment
		tteam_equipment = r.Team2.Equipment
	} */

	team1equipment := r.Team1.Equipment
	team2equipment := r.Team2.Equipment

	r.SurvivingT1 = int(math.Round(CSFNormalDistribution_std_4(float64(team1equipment), float64(team2equipment), r.gameRules.CSF_r, 0, 5)))
	r.SurvivingT2 = int(math.Round(CSFNormalDistribution_std_4(float64(team2equipment), float64(team1equipment), r.gameRules.CSF_r, 0, 5)))

}

func (r *Round) determineRemainingEquipment() {
	team1equipment := r.Team1.Equipment
	team2equipment := r.Team2.Equipment

	// Careful! Here the inverse is used. Because the remaining equipement is largely by the other
	// teams equipment and their chances of winning fights, i.e.
	//TODO: This should potentiall be changed i.e. to factor in the equipment lost by the other team
	r.EquipmentT1 = (CSFNormalDistribution_std_4(team2equipment, team1equipment, r.gameRules.CSF_r, (team1equipment / 5 * 0.8), ((math.Max(team1equipment, team2equipment) / 5) * 1.2))) * float64(r.SurvivingT1)
	r.EquipmentT2 = (CSFNormalDistribution_std_4(team1equipment, team2equipment, r.gameRules.CSF_r, (team2equipment / 5 * 0.8), ((math.Max(team1equipment, team2equipment) / 5) * 1.2))) * float64(r.SurvivingT2)
}

func (r *Round) determineBombplant() {
	// Placeholder for CSF logic

	ctteam_equipment := r.Team1.Equipment
	tteam_equipment := r.Team2.Equipment

	if r.CTTeam {
		ctteam_equipment = r.Team1.Equipment
		tteam_equipment = r.Team2.Equipment
	}

	r.BombPlanted = bool_CSF_simple(ctteam_equipment, tteam_equipment, r.gameRules.CSF_r)

}

func (r *Round) determineWinner() {
	// Placeholder for CSF logic
	team1equipment := r.Team1.Equipment
	team2equipment := r.Team2.Equipment

	r.WinnerTeam = bool_CSF_simple(team1equipment, team2equipment, r.gameRules.CSF_r)

	if r.WinnerTeam == r.CTTeam {
		r.WinnerSide = false // CT wins
	} else {
		r.WinnerSide = true // T wins
	}

}
