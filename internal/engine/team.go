package engine

// Team represents a team in the simulation with its properties and methods.
type Team struct {
	Name                string
	Side                bool    // true for CT, false for T
	Funds               float64 // Total funds available for the team
	Earned              float64 // Total funds earned by the team
	Equipment           float64 // Total equipment value for the team
	Remainingequipement float64 // Remaining equipment value after spending
	defaultEquipment    float64 // Default equipment if no other equipment
	Playersalive        int
	Score               int
	Consecutiveloss     int // Consecutive losses for the team
	Spent               float64
	Strategy            string // Strategy name for the team
}

func NewTeam(name string, startingfunds float64, side bool, defaultequipment float64, strategy string) *Team {
	return &Team{
		Name:                name,
		Funds:               5 * startingfunds, // Starting funds
		Playersalive:        5,                 // Starting with 5 players
		Side:                side,
		Equipment:           5 * defaultequipment, // Starting equipment
		defaultEquipment:    defaultequipment,     // Default equipment cost
		Remainingequipement: 0,                    // Initialize remaining equipment
		Score:               0,
		Consecutiveloss:     0, // Initialize consecutive losses
		Spent:               0,
		Strategy:            strategy,
	}
}

func (t *Team) RoundEnd(winner bool, fundsearned float64, playersalive int, remainingequipement float64) {
	if winner {
		t.Score += 1 // Increment score for winning team
		t.Consecutiveloss = 0
	} else {
		t.Consecutiveloss++
	}
	t.EarnFunds(fundsearned)                    // Earn funds based on round outcome
	t.Playersalive = playersalive               // Update players alive after the round
	t.Remainingequipement = remainingequipement // Update remaining equipment
}

func (t *Team) NewOT(OTfunds float64) {
	t.Equipment = 5 * t.defaultEquipment // Reset equipment for overtime
	t.Funds = 5 * OTfunds                // Reset funds for overtime
	t.Playersalive = 5                   // Reset players alive for overtime
}

func (t *Team) Sideswitch(OT bool, startingfunds float64, OTfunds float64) {
	t.Side = !t.Side // Switch side if needed
	if OT {
		t.Funds = 5 * OTfunds // Reset funds for overtime
	} else {
		t.Funds = 5 * startingfunds // Reset funds for regular rounds
	}
	t.Equipment = 5 * t.defaultEquipment // Reset equipment for new side
	t.Playersalive = 5                   // Reset players alive for new side
}

func (t *Team) NewRound() {
	t.Spent = 0                                                                            // Reset spent funds at the start of the round
	t.Earned = 0                                                                           // Reset earned funds at the start of the round
	t.Equipment = t.Remainingequipement + (float64(5-t.Playersalive) * t.defaultEquipment) // Add default equipment cost for dead players
	t.Playersalive = 5                                                                     // Set all players alive at the start of the round
}

// for now it is a set of variables, in the future it could be a json file with information
// the state of the game, team, round etc.
func (t *Team) BuyPhase(Round int, ScoreOpo int) {

	investment := CallStrategy(t, Round, ScoreOpo) // Call the strategy manager to get investment amount

	t.SpendFunds(investment) // Spend investment amount

}

func (t *Team) EarnFunds(amount float64) {
	t.Funds += amount
	t.Earned += amount
}

func (t *Team) SpendFunds(amount float64) {
	if amount <= t.Funds {
		t.Funds -= amount
		t.Spent += amount
		t.Equipment += amount
	} else {
		// Handle insufficient funds
	}
}
