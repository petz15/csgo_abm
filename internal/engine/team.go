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
	Spent            int
	Strategy         models.Strategy
}

func NewTeam(name string, startingfunds int, side bool, defaultequipment int, strategy models.Strategy) *Team {
	return &Team{
		Name:             name,
		Funds:            startingfunds, // Starting funds
		Playersalive:     5,             // Starting with 5 players
		Side:             side,
		DefaultEquipment: defaultequipment, // Default equipment cost
		Equipment:        0,
		Score:            0,
		Spent:            0,
		Strategy:         strategy,
	}
}

func (t *Team) RoundEnd(winner bool, fundsearned int, playersalive int, remainingequipement int) {
	if winner {
		t.Score += 1 // Increment score for winning team
	}
	t.EarnFunds(fundsearned)          // Earn funds based on round outcome
	t.Playersalive = playersalive     // Update players alive after the round
	t.Equipment = remainingequipement // Update remaining equipment
}

func (t *Team) NewOT(OTfunds int, OTEquipment int, switchsides bool) {
	if switchsides {
		t.Side = !t.Side // Switch side if needed
	}
	t.Spent = 0
	t.Earned = 0
	t.Equipment = 5 * OTEquipment // Reset equipment for overtime
	t.Funds = OTfunds             // Reset funds for overtime
	t.Playersalive = 5            // Reset players alive for overtime
}

func (t *Team) NewRound(switchside bool) {
	t.Spent = 0                                                             // Reset spent funds at the start of the round
	t.Earned = 0                                                            // Reset earned funds at the start of the round
	t.Equipment = t.Equipment + ((5 - t.Playersalive) * t.DefaultEquipment) // Add default equipment cost for dead players
	t.Playersalive = 5                                                      // Reset players alive for the new round
	if switchside {
		t.Side = !t.Side                     // Switch side if needed
		t.Equipment = 5 * t.DefaultEquipment // Reset equipment for Side switch
	}
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
