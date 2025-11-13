package engine

import "math"

type Round struct {
	RoundNumber      int
	is_T1_CT         bool // true if Team1 is CT, false if Team1 is T
	Calc_Outcome     RoundOutcome
	is_T1_WinnerTeam bool //true if Team1 wins, false if Team2 wins
	Sideswitch       bool // True if sideswitch has occurred
	gameRules        *GameRules
	OT               bool // True if this is an overtime round
}

func NewRound(T1 *Team, T2 *Team, roundNumber int, ctteam bool, sidewitch bool, gamerules *GameRules, ot bool) *Round {

	T1.NewRound()
	T2.NewRound()

	return &Round{
		RoundNumber: roundNumber,
		is_T1_CT:    ctteam,
		Sideswitch:  sidewitch,
		gameRules:   gamerules,
		OT:          ot,
	}
}

// solution for what information gets passed to the teams, could be a json file in the gamerules which specifies
// which variables are passed to the teams. No clue how this will be done yet, but it is an idea.
func (r *Round) BuyPhase(Team1 *Team, Team2 *Team) {
	Team1.BuyPhase(r.RoundNumber, r.OT, Team2, *r.gameRules) // Call the team's buy phase logic
	Team2.BuyPhase(r.RoundNumber, r.OT, Team1, *r.gameRules)

}

func (r *Round) RoundStart(Team1 *Team, Team2 *Team) {
	// Simulate the round using the comprehensive DetermineRoundOutcome function
	r.determineRoundOutcome(Team1, Team2)

}

func (r *Round) RoundEnd(Team1 *Team, Team2 *Team) {
	//determine funds earned
	r.determineFundsEarned(Team1, Team2)

	//update Teams
	Team1.RoundEnd(r.is_T1_WinnerTeam)
	Team2.RoundEnd(!r.is_T1_WinnerTeam)
}

func (r *Round) determineFundsEarned(Team1 *Team, Team2 *Team) {
	// Determine funds earned based on round outcome
	// https://counterstrike.fandom.com/wiki/Money -> be sure to check CSGO values
	winnerFunds := 0.0
	loserFunds := 0.0

	switch r.Calc_Outcome.ReasonCode {
	case 1:
		winnerFunds += 3500 * 5 //bomb exploded
		winnerFunds += 300      //plant bonus
	case 2:
		winnerFunds += 3250 * 5
		if r.Calc_Outcome.BombPlanted { //plant bonus
			winnerFunds += 300
		}
	case 3:
		winnerFunds += 3500 * 5
		winnerFunds += 300            //defuse bonus
		loserFunds += (800 * 5) + 300 //bomb planted bonus
	case 4:
		winnerFunds += 3250 * 5
	}

	var ctteam *Team
	var tteam *Team
	if r.is_T1_CT {
		ctteam = Team1
		tteam = Team2
		//while we're here, set RE EQ values i.e. saved equipment values
		Team1.SetREEqValue(sumArray(r.Calc_Outcome.CTEquipmentPerPlayer))
		Team2.SetREEqValue(sumArray(r.Calc_Outcome.TEquipmentPerPlayer))
	} else {
		ctteam = Team2
		tteam = Team1
		//while we're here, set RE EQ values i.e. saved equipment values
		Team1.SetREEqValue(sumArray(r.Calc_Outcome.TEquipmentPerPlayer))
		Team2.SetREEqValue(sumArray(r.Calc_Outcome.CTEquipmentPerPlayer))
	}

	var lossbonus int
	//Loser bonus evaluation
	if r.Calc_Outcome.CTWins {
		lossbonus = r.LossBonusCalculation(tteam)
	} else {
		lossbonus = r.LossBonusCalculation(ctteam)
	}

	if r.gameRules.LossBonusCalc == true {
		if r.Calc_Outcome.CTWins {
			tteam.SetCurrentLossBonusLevel(tteam.GetCurrentLossBonusLevel() + 1)
			ctteam.SetCurrentLossBonusLevel(ctteam.GetCurrentLossBonusLevel() - 1)
		} else {
			ctteam.SetCurrentLossBonusLevel(ctteam.GetCurrentLossBonusLevel() + 1)
			tteam.SetCurrentLossBonusLevel(tteam.GetCurrentLossBonusLevel() - 1)
		}
	} else {
		if r.Calc_Outcome.CTWins {
			tteam.SetCurrentLossBonusLevel(tteam.GetCurrentLossBonusLevel() + 1)
			ctteam.SetCurrentLossBonusLevel(0)
		} else {
			ctteam.SetCurrentLossBonusLevel(ctteam.GetCurrentLossBonusLevel() + 1)
			tteam.SetCurrentLossBonusLevel(0)
		}
	}

	// Kills and loss bonus

	if r.Calc_Outcome.CTWins {
		loserFunds += (5 - float64(r.Calc_Outcome.CTSurvivors)) * 300
		winnerFunds += (5 - float64(r.Calc_Outcome.TSurvivors)) * 300

		// Add loss bonus to losing team
		// Reduction for surviving T players if round end reason is 4
		reduction := 0
		if r.Calc_Outcome.ReasonCode == 4 {
			reduction = r.Calc_Outcome.TSurvivors
		}
		loserFunds += float64(lossbonus) * float64(5-reduction)
	} else {
		loserFunds += (5 - float64(r.Calc_Outcome.TSurvivors)) * 300
		winnerFunds += (5 - float64(r.Calc_Outcome.CTSurvivors)) * 300

		// Add loss bonus to losing team
		loserFunds += float64(lossbonus) * 5
	}

	FundsearnedT1 := 0.0
	FundsearnedT2 := 0.0

	if r.is_T1_WinnerTeam {
		FundsearnedT1 = winnerFunds
		FundsearnedT2 = loserFunds
	} else {
		FundsearnedT1 = loserFunds
		FundsearnedT2 = winnerFunds
	}
	// Ensure funds do not exceed maximum allowed

	Team1.SetEarned(FundsearnedT1)
	Team1.SetFunds(math.Min(Team1.GetCurrentFunds(), r.gameRules.MaxFunds))

	Team2.SetEarned(FundsearnedT2)
	Team2.SetFunds(math.Min(Team2.GetCurrentFunds(), r.gameRules.MaxFunds))
}

func (r *Round) LossBonusCalculation(loserteam *Team) int {
	// Calculate loss bonus based on consecutive losses

	lossBonus := 0
	lossbonuslevel := loserteam.GetCurrentLossBonusLevel()
	if lossbonuslevel >= 4 {
		lossBonus = 3400 //Loss bonus for 5th loss and beyond
	} else if lossbonuslevel == 3 {
		lossBonus = 2900 // Loss bonus for fourth loss
	} else if lossbonuslevel == 2 {
		lossBonus = 2400 // Loss bonus for third loss
	} else if lossbonuslevel == 1 {
		lossBonus = 1900 // Loss bonus for second loss
	} else {
		lossBonus = 1400 // Loss bonus for first loss
	}
	return lossBonus
}

// determineRoundOutcome uses the comprehensive ABM-based probability function
// to determine all aspects of the round outcome in one call
func (r *Round) determineRoundOutcome(Team1 *Team, Team2 *Team) {
	// Get equipment values for CSF calculation

	ctequipment := 0.0
	tequipment := 0.0

	if r.is_T1_CT {
		ctequipment = Team1.RoundData[r.RoundNumber-1].FTE_Eq_value
		tequipment = Team2.RoundData[r.RoundNumber-1].FTE_Eq_value
	} else {
		ctequipment = Team2.RoundData[r.RoundNumber-1].FTE_Eq_value
		tequipment = Team1.RoundData[r.RoundNumber-1].FTE_Eq_value
	}

	// Get comprehensive round outcome from ABM distributions (uses CT win probability)
	outcome := DetermineRoundOutcome(ctequipment, tequipment)

	// Determine which team won
	r.is_T1_WinnerTeam = !outcome.CTWins
	if r.is_T1_CT {
		// If Team1 is CT, they win when CT wins
		r.is_T1_WinnerTeam = outcome.CTWins
	}

}
