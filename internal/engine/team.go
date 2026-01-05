package engine

// Team represents a team in the simulation with its properties and methods.
type Team struct {
	Name      string
	Strategy  string // Strategy name for the team
	RoundData []Team_RoundData
}

type Team_RoundData struct {
	is_Side_CT            bool    // true for CT, false for T
	Funds                 float64 // Total funds available for the team at the start of the round
	Funds_start           float64 // Funds at the start of the round
	Earned                float64 // Total funds earned by the team at the conclusion of the round
	RS_Eq_value           float64 // Round start equipment value
	FTE_Eq_value          float64 // Freeze time end equipment value (after buy phase)
	RE_Eq_value           float64 // Round end equipment value
	Survivors             int     // Number of surviving players at the end of the round
	Score_Start           int
	Score_End             int     //updated at the end of the round
	Consecutiveloss       int     // Consecutive losses for the team, updated at the end of the round
	Consecutivewins       int     // Consecutive wins for the team, updated at the end of the round
	Consecutivewins_start int     // Consecutive wins at the start of the round
	Consecutiveloss_start int     // Consecutive losses at the start of the round
	LossBonusLevel        int     // Level of loss bonus calculated at the end of the round
	Spent                 float64 // Total funds spent by the team during buy time
}

func NewTeam(name string, startingfunds float64, side bool, defaultequipment float64, strategy string) *Team {

	new_RD := Team_RoundData{
		is_Side_CT:            side,
		Funds:                 5 * startingfunds, // Starting funds
		Funds_start:           5 * startingfunds, // Funds at the start of the round
		Earned:                0,
		RS_Eq_value:           5 * defaultequipment,
		FTE_Eq_value:          0,
		RE_Eq_value:           0,
		Survivors:             0,
		Score_Start:           0,
		Score_End:             0,
		Consecutiveloss:       0,
		Consecutivewins:       0,
		Consecutivewins_start: 0,
		Consecutiveloss_start: 0,
		LossBonusLevel:        1,
		Spent:                 0,
	}

	return &Team{
		Name:      name,
		Strategy:  strategy,
		RoundData: []Team_RoundData{new_RD},
	}
}

func (t *Team) NewRound(defaultequipment float64) {
	// Create a new round data entry for the team
	previousRound := t.RoundData[len(t.RoundData)-1]
	newRoundData := Team_RoundData{
		is_Side_CT:            previousRound.is_Side_CT,
		Funds:                 previousRound.Funds,
		Funds_start:           previousRound.Funds,
		Earned:                0,
		RS_Eq_value:           previousRound.RE_Eq_value + (float64(5-previousRound.Survivors) * defaultequipment),
		FTE_Eq_value:          0,
		RE_Eq_value:           0,
		Survivors:             0,
		Score_Start:           previousRound.Score_End,
		Score_End:             previousRound.Score_End,
		Consecutiveloss:       previousRound.Consecutiveloss,
		Consecutivewins:       previousRound.Consecutivewins,
		Consecutivewins_start: previousRound.Consecutivewins,
		Consecutiveloss_start: previousRound.Consecutiveloss,
		LossBonusLevel:        previousRound.LossBonusLevel,
		Spent:                 0,
	}
	t.RoundData = append(t.RoundData, newRoundData)
}

func (t *Team) RoundEnd(winner bool) {
	RD := &t.RoundData[len(t.RoundData)-1]
	if winner {
		RD.Score_End = RD.Score_Start + 1 // Increment score for winning team
		RD.Consecutiveloss = 0
		RD.Consecutivewins++
	} else {
		RD.Consecutiveloss++
		RD.Consecutivewins = 0
	}

}

func (t *Team) NewOT(OTfunds float64, OTEquipment float64) {
	t.RoundData[len(t.RoundData)-1].RS_Eq_value = 5 * OTEquipment // Reset equipment for overtime
	t.RoundData[len(t.RoundData)-1].Funds = 5 * OTfunds           // Reset funds for overtime
	t.RoundData[len(t.RoundData)-1].Funds_start = 5 * OTfunds

	t.RoundData[len(t.RoundData)-1].Consecutiveloss = 0
	t.RoundData[len(t.RoundData)-1].Consecutivewins = 0
	t.RoundData[len(t.RoundData)-1].Consecutivewins_start = 0
	t.RoundData[len(t.RoundData)-1].Consecutiveloss_start = 0
	t.RoundData[len(t.RoundData)-1].LossBonusLevel = 1

}

func (t *Team) Sideswitch(OT bool, startingfunds float64, defaultEquipment float64, OTfunds float64, OTEquipment float64) {
	t.RoundData[len(t.RoundData)-1].is_Side_CT = !t.RoundData[len(t.RoundData)-1].is_Side_CT // Switch side if needed
	if OT {
		t.NewOT(OTfunds, OTEquipment)
	} else {
		t.RoundData[len(t.RoundData)-1].Funds = 5 * startingfunds // Reset funds for new half
		t.RoundData[len(t.RoundData)-1].Funds_start = 5 * startingfunds
		t.RoundData[len(t.RoundData)-1].RS_Eq_value = 5 * defaultEquipment // Reset equipment for new half

		t.RoundData[len(t.RoundData)-1].Consecutiveloss = 0
		t.RoundData[len(t.RoundData)-1].Consecutivewins = 0
		t.RoundData[len(t.RoundData)-1].Consecutivewins_start = 0
		t.RoundData[len(t.RoundData)-1].Consecutiveloss_start = 0
		t.RoundData[len(t.RoundData)-1].LossBonusLevel = 1
	}
}

// for now it is a set of variables, in the future it could be a json file with information
// the state of the game, team, round etc.
func (t *Team) BuyPhase(Round int, ot bool, t2 *Team, gameR GameRules, g *Game) {

	investment := CallStrategy(t, t2, Round, ot, gameR, g) // Call the strategy manager to get investment amount

	t.SpendFunds(investment) // Spend investment amount

}

func (t *Team) SpendFunds(amount float64) {
	RD := &t.RoundData[len(t.RoundData)-1]
	if amount < 0 {
		amount = 0
	}

	if amount <= RD.Funds {
		RD.Funds -= amount
		RD.Spent += amount
		RD.FTE_Eq_value += amount + t.GetRSEquipment()
	} else if amount > RD.Funds {
		//limit to max
		max_amount := RD.Funds
		RD.Funds -= max_amount
		RD.Spent += max_amount
		RD.FTE_Eq_value += max_amount + t.GetRSEquipment()
	}
}

func (t *Team) GetCurrentFunds() float64 {
	RD := &t.RoundData[len(t.RoundData)-1]
	return RD.Funds
}

func (t *Team) GetCurrentLossBonusLevel() int {
	RD := &t.RoundData[len(t.RoundData)-1]
	return RD.LossBonusLevel
}

func (t *Team) GetConsecutiveloss() int {
	RD := &t.RoundData[len(t.RoundData)-1]
	return RD.Consecutiveloss
}

func (t *Team) GetConsecutivewins() int {
	RD := &t.RoundData[len(t.RoundData)-1]
	return RD.Consecutivewins
}

func (t *Team) SetCurrentLossBonusLevel(value int, maxValue int) {
	RD := &t.RoundData[len(t.RoundData)-1]
	if value < 0 {
		value = 0
	}
	if value > maxValue {
		value = maxValue
	}
	RD.LossBonusLevel = value
}

func (t *Team) GetScore() int {
	RD := &t.RoundData[len(t.RoundData)-1]
	return RD.Score_End
}

func (t *Team) SetEarned(amount float64) {
	RD := &t.RoundData[len(t.RoundData)-1]
	RD.Earned = amount
	RD.Funds += amount
}

func (t *Team) SetFunds(amount float64) {
	RD := &t.RoundData[len(t.RoundData)-1]
	RD.Funds = amount
}

func (t *Team) SetREEqValue(value float64) {
	RD := &t.RoundData[len(t.RoundData)-1]
	RD.RE_Eq_value = value
}

func (t *Team) GetSide() bool {
	RD := &t.RoundData[len(t.RoundData)-1]
	return RD.is_Side_CT
}

func (t *Team) GetRSEquipment() float64 {
	RD := &t.RoundData[len(t.RoundData)-1]
	return float64(RD.RS_Eq_value)
}

// Cleanup clears historical round data to free memory
// Should only be called after game is complete and results are extracted
func (t *Team) Cleanup() {
	// Keep only the last round's data for final state
	if len(t.RoundData) > 0 {
		lastRound := t.RoundData[len(t.RoundData)-1]
		t.RoundData = []Team_RoundData{lastRound}
	}
}

func (t *Team) SetSurvivors(survivors int) {
	RD := &t.RoundData[len(t.RoundData)-1]
	RD.Survivors = survivors
}

func (t *Team) GetpreviousSurvivors() int {
	if len(t.RoundData) < 2 {
		return 0
	}
	RD := &t.RoundData[len(t.RoundData)-2]
	return RD.Survivors
}
