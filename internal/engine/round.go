package engine

import "math"

type Round struct {
	RoundNumber    int
	IsT1CT         bool `json:"is_t1_ct"` // true if Team1 is CT, false if Team1 is T
	Calc_Outcome   RoundOutcome
	IsT1WinnerTeam bool       `json:"is_t1_winner_team"` //true if Team1 wins, false if Team2 wins
	Sideswitch     bool       // True if sideswitch has occurred
	gameRules      *GameRules `json:"-"` // Don't export game rules
	OT             bool       // True if this is an overtime round
	game           *Game      `json:"-"` // Don't export game rules
}

func NewRound(T1 *Team, T2 *Team, roundNumber int, ctteam bool, gamerules *GameRules, ot bool, g *Game) *Round {

	if roundNumber != 1 { //avoid calling NewRound on first round twice
		T1.NewRound(gamerules.DefaultEquipment)
		T2.NewRound(gamerules.DefaultEquipment)
	}

	return &Round{
		RoundNumber: roundNumber,
		IsT1CT:      ctteam,
		Sideswitch:  false,
		gameRules:   gamerules,
		OT:          ot,
		game:        g,
	}
}

// HandleSideSwitch manages side switching at halftime
func (r *Round) HandleSideSwitch(Team1 *Team, Team2 *Team) {
	r.Sideswitch = true
	r.IsT1CT = !r.IsT1CT
	Team1.Sideswitch(r.OT, r.gameRules.StartingFunds, r.gameRules.DefaultEquipment, r.gameRules.OTFunds, r.gameRules.OTEquipment)
	Team2.Sideswitch(r.OT, r.gameRules.StartingFunds, r.gameRules.DefaultEquipment, r.gameRules.OTFunds, r.gameRules.OTEquipment)
}

// HandleOTStart manages the start of overtime
func (r *Round) HandleOTStart(Team1 *Team, Team2 *Team) {
	r.OT = true
	Team1.NewOT(r.gameRules.OTFunds, r.gameRules.OTEquipment)
	Team2.NewOT(r.gameRules.OTFunds, r.gameRules.OTEquipment)
}

// HandleOTSideSwitch manages side switching during overtime
func (r *Round) HandleOTSideSwitch(Team1 *Team, Team2 *Team) {
	r.Sideswitch = true
	r.IsT1CT = !r.IsT1CT
	Team1.Sideswitch(true, r.gameRules.StartingFunds, r.gameRules.DefaultEquipment, r.gameRules.OTFunds, r.gameRules.OTEquipment)
	Team2.Sideswitch(true, r.gameRules.StartingFunds, r.gameRules.DefaultEquipment, r.gameRules.OTFunds, r.gameRules.OTEquipment)
}

// solution for what information gets passed to the teams, could be a json file in the gamerules which specifies
// which variables are passed to the teams. No clue how this will be done yet, but it is an idea.
func (r *Round) BuyPhase(Team1 *Team, Team2 *Team) {
	Team1.BuyPhase(r.RoundNumber, r.OT, Team2, *r.gameRules, r.game) // Call the team's buy phase logic
	Team2.BuyPhase(r.RoundNumber, r.OT, Team1, *r.gameRules, r.game)

}

func (r *Round) RoundStart(Team1 *Team, Team2 *Team) {
	// Simulate the round using the comprehensive CalculateRoundOutcome function
	r.CalculateRoundOutcome(Team1, Team2)

}

func (r *Round) RoundEnd(Team1 *Team, Team2 *Team) {
	//determine funds earned
	r.determineFundsEarned(Team1, Team2)

	//update Teams
	Team1.RoundEnd(r.IsT1WinnerTeam)
	Team2.RoundEnd(!r.IsT1WinnerTeam)
}

func (r *Round) determineFundsEarned(Team1 *Team, Team2 *Team) {
	// Determine funds earned based on round outcome
	// https://counterstrike.fandom.com/wiki/Money -> be sure to check CSGO values
	// https://www.rockpapershotgun.com/csgo-economy-guide
	winnerFunds := 0.0
	loserFunds := 0.0

	switch r.Calc_Outcome.ReasonCode {
	case 1: //T win by bomb explosion
		winnerFunds += r.gameRules.RoundOutcomeReward[0] * 5
		winnerFunds += r.gameRules.BombplantReward //plant bonus
	case 2: //T win by elimination
		winnerFunds += r.gameRules.RoundOutcomeReward[1] * 5
		if r.Calc_Outcome.BombPlanted { //plant bonus
			winnerFunds += r.gameRules.BombplantReward
		}
	case 3: //CT win by defuse
		winnerFunds += r.gameRules.RoundOutcomeReward[2] * 5
		winnerFunds += r.gameRules.BombdefuseReward      //defuse bonus
		loserFunds += r.gameRules.BombplantRewardall * 5 //plant bonus for all T players
		loserFunds += r.gameRules.BombplantReward        //defuse bonus for all T players
	case 4: //CT win by elimination
		winnerFunds += r.gameRules.RoundOutcomeReward[3] * 5
	}

	var ctteam *Team
	var tteam *Team
	if r.IsT1CT {
		ctteam = Team1
		tteam = Team2
		//while we're here, set RE EQ values i.e. saved equipment values and survivors
		Team1.SetREEqValue(math.Floor(sumArray(r.Calc_Outcome.CTEquipmentPerPlayer))) // rounding down to make numbers cleaner
		Team2.SetREEqValue(math.Floor(sumArray(r.Calc_Outcome.TEquipmentPerPlayer)))
		Team1.SetSurvivors(r.Calc_Outcome.CTSurvivors)
		Team2.SetSurvivors(r.Calc_Outcome.TSurvivors)

	} else {
		ctteam = Team2
		tteam = Team1
		//while we're here, set RE EQ values i.e. saved equipment values
		Team1.SetREEqValue(math.Floor(sumArray(r.Calc_Outcome.TEquipmentPerPlayer)))
		Team2.SetREEqValue(math.Floor(sumArray(r.Calc_Outcome.CTEquipmentPerPlayer)))
		Team1.SetSurvivors(r.Calc_Outcome.TSurvivors)
		Team2.SetSurvivors(r.Calc_Outcome.CTSurvivors)
	}

	var lossbonus int
	//Loser bonus evaluation
	if r.Calc_Outcome.CTWins {
		lossbonus = r.LossBonusCalculation(tteam)
	} else {
		lossbonus = r.LossBonusCalculation(ctteam)
	}

	maxLossLevel := len(r.gameRules.LossBonus) - 1
	// Update loss bonus levels

	if r.gameRules.LossBonusCalc {

		if r.Calc_Outcome.CTWins {
			tteam.SetCurrentLossBonusLevel(tteam.GetCurrentLossBonusLevel()+1, maxLossLevel)
			ctteam.SetCurrentLossBonusLevel(ctteam.GetCurrentLossBonusLevel()-1, maxLossLevel)
		} else {
			ctteam.SetCurrentLossBonusLevel(ctteam.GetCurrentLossBonusLevel()+1, maxLossLevel)
			tteam.SetCurrentLossBonusLevel(tteam.GetCurrentLossBonusLevel()-1, maxLossLevel)
		}
	} else {
		if r.Calc_Outcome.CTWins {
			tteam.SetCurrentLossBonusLevel(tteam.GetCurrentLossBonusLevel()+1, maxLossLevel)
			ctteam.SetCurrentLossBonusLevel(0, maxLossLevel)
		} else {
			ctteam.SetCurrentLossBonusLevel(ctteam.GetCurrentLossBonusLevel()+1, maxLossLevel)
			tteam.SetCurrentLossBonusLevel(0, maxLossLevel)
		}
	}

	// Kills and loss bonus

	if r.Calc_Outcome.CTWins {
		loserFunds += float64((5 - (r.Calc_Outcome.CTSurvivors))) * (r.gameRules.EliminationReward + (r.gameRules.AdditionalReward_T_Elimination * 5))
		winnerFunds += float64((5 - (r.Calc_Outcome.TSurvivors))) * (r.gameRules.EliminationReward + (r.gameRules.AdditionalReward_CT_Elimination * 5))

		// Add loss bonus to losing team
		// Reduction for surviving T players if round end reason is 4
		reduction := 0
		if r.Calc_Outcome.ReasonCode == 4 {
			reduction = r.Calc_Outcome.TSurvivors
		}
		loserFunds += float64(lossbonus) * float64(5-reduction)
	} else {
		loserFunds += float64((5 - (r.Calc_Outcome.TSurvivors))) * (r.gameRules.EliminationReward + (r.gameRules.AdditionalReward_CT_Elimination * 5))
		winnerFunds += float64((5 - (r.Calc_Outcome.CTSurvivors))) * (r.gameRules.EliminationReward + (r.gameRules.AdditionalReward_T_Elimination * 5))

		// Add loss bonus to losing team
		loserFunds += float64(lossbonus) * 5
	}

	FundsearnedT1 := 0.0
	FundsearnedT2 := 0.0

	if r.IsT1WinnerTeam {
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
	if lossbonuslevel < len(r.gameRules.LossBonus) {
		lossBonus = int(r.gameRules.LossBonus[lossbonuslevel])
	} else {
		lossBonus = int(r.gameRules.LossBonus[len(r.gameRules.LossBonus)-1])
	}
	return lossBonus
}

// CalculateRoundOutcome uses the comprehensive ABM-based probability function
// to determine all aspects of the round outcome in one call
func (r *Round) CalculateRoundOutcome(Team1 *Team, Team2 *Team) {
	// Get equipment values for CSF calculation

	//ensure that equipment values are at least 1.0 to avoid zero multipliers
	ctequipment := 1.0
	tequipment := 1.0

	if r.IsT1CT {
		ctequipment += Team1.RoundData[r.RoundNumber-1].FTE_Eq_value
		tequipment += Team2.RoundData[r.RoundNumber-1].FTE_Eq_value
	} else {
		ctequipment += Team2.RoundData[r.RoundNumber-1].FTE_Eq_value
		tequipment += Team1.RoundData[r.RoundNumber-1].FTE_Eq_value
	}

	// Get comprehensive round outcome from ABM distributions (uses CT win probability)
	r.Calc_Outcome = DetermineRoundOutcome(ctequipment, tequipment, r.game.rng, *r.gameRules)

	// Determine which team won
	r.IsT1WinnerTeam = !r.Calc_Outcome.CTWins
	if r.IsT1CT {
		// If Team1 is CT, they win when CT wins
		r.IsT1WinnerTeam = r.Calc_Outcome.CTWins
	}

}
