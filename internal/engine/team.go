package engine

import "csgo-economy-sim/internal/models"

// Team represents a team in the simulation with its properties and methods.
type Team struct {
	Name             string
	Side             bool // true for CT, false for T
	Funds            int
	Earned           int
	Equipment        int
	DefaultEquipment int // Default equipment if no other equipment
	Playersalive     int
	Score            int
	Consecutiveloss  int // Consecutive losses for the team
	Spent            int
	Strategy         models.Strategy
}

func NewTeam(name string, startingfunds int, side bool, defaultequipment int, strategy models.Strategy) *Team {
	return &Team{
		Name:             name,
		Funds:            5 * startingfunds, // Starting funds
		Playersalive:     5,                 // Starting with 5 players
		Side:             side,
		DefaultEquipment: defaultequipment, // Default equipment cost
		Equipment:        0,
		Score:            0,
		Consecutiveloss:  0, // Initialize consecutive losses
		Spent:            0,
		Strategy:         strategy,
	}
}

func (t *Team) RoundEnd(winner bool, fundsearned int, playersalive int, remainingequipement int) {
	if winner {
		t.Score += 1 // Increment score for winning team
		t.Consecutiveloss = 0
	} else {
		t.Consecutiveloss++
	}
	t.EarnFunds(fundsearned)          // Earn funds based on round outcome
	t.Playersalive = playersalive     // Update players alive after the round
	t.Equipment = remainingequipement // Update remaining equipment
}

func (t *Team) NewOT(OTfunds int) {
	t.Equipment = 5 * t.DefaultEquipment // Reset equipment for overtime
	t.Funds = 5 * OTfunds                // Reset funds for overtime
	t.Playersalive = 5                   // Reset players alive for overtime
}

func (t *Team) Sideswitch(OT bool, startingfunds int, OTfunds int) {
	t.Side = !t.Side // Switch side if needed
	if OT {
		t.Funds = 5 * OTfunds // Reset funds for overtime
	} else {
		t.Funds = 5 * startingfunds // Reset funds for regular rounds
	}
	t.Equipment = 5 * t.DefaultEquipment // Reset equipment for new side
	t.Playersalive = 5                   // Reset players alive for new side
}

func (t *Team) NewRound() {
	t.Spent = 0                                                             // Reset spent funds at the start of the round
	t.Earned = 0                                                            // Reset earned funds at the start of the round
	t.Equipment = t.Equipment + ((5 - t.Playersalive) * t.DefaultEquipment) // Add default equipment cost for dead players
	t.Playersalive = 5                                                      // Set all players alive at the start of the round
}

// this is unclear for now, as strategy is not defined yet, especially what information a team receives
func (t *Team) BuyPhase(Round int, ScoreOpo int) {
	// Logic for buying equipment at the start of the round
	// This could include checking funds, buying weapons, armor, etc.
	if t.Funds >= t.DefaultEquipment {
		t.SpendFunds(t.DefaultEquipment)  // Spend default equipment cost
		t.Equipment += t.DefaultEquipment // Add to equipment
	}
}

func (t *Team) EarnFunds(amount int) {
	t.Funds += amount
	t.Earned += amount
}

func (t *Team) SpendFunds(amount int) {
	if amount <= t.Funds {
		t.Funds -= amount
		t.Spent += amount
		t.Equipment += amount // Assuming equipment cost is deducted from funds
	} else {
		// Handle insufficient funds
	}
}
