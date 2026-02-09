package engine

import (
	"math"
	"math/rand"
)

//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//CAUTION There is currently a max number of overtime set to 50. Meaning an overtime can have max 15 * 2 + 50 * 6= 330 rounds.
//This is a temporary solution to prevent infinite loops in the game simulation.

type Game struct {
	ID             string
	GameinProgress bool
	CurrentRound   int
	is_T1_CT       bool //true if Team1 is CT, false if T1 is T
	Score          [2]int
	Rounds         []Round
	OT             bool
	OTcounter      int
	firsthalf      bool
	sideswitch     bool
	GameRules      GameRules
	Is_T1_Winner   bool // true if T1 wins, false if T2 wins
	Team1          *Team
	Team2          *Team
	rng            *rand.Rand // Thread-safe RNG for this game instance
}

// NewGame creates a new game with pre-validated GameRules object (optimized for batch simulations)
func NewGame(id string, Team1Name string, Team1Strategy string, Team2Name string, Team2Strategy string, gameRules GameRules) *Game {

	// Create a thread-safe RNG for this game instance
	rng := rand.New(rand.NewSource(rand.Int63()))

	currentCT := rng.Intn(2) == 0

	Team_1 := NewTeam(Team1Name, gameRules.StartingFunds, currentCT, gameRules.DefaultEquipment, Team1Strategy)
	Team_2 := NewTeam(Team2Name, gameRules.StartingFunds, !currentCT, gameRules.DefaultEquipment, Team2Strategy)

	return &Game{
		ID:             id,
		Team1:          Team_1,
		Team2:          Team_2,
		is_T1_CT:       currentCT,
		CurrentRound:   1,
		OT:             false,
		OTcounter:      0,
		firsthalf:      true,
		sideswitch:     false,
		Score:          [2]int{0, 0},
		GameRules:      gameRules,
		GameinProgress: false,
		rng:            rng,
	}

}

// SetSeed sets the RNG seed for this game to ensure reproducible outcomes per game/series
func (g *Game) SetSeed(seed int64) {
	g.rng = rand.New(rand.NewSource(seed))
}

func (g *Game) Start() {
	g.GameinProgress = true

	for g.GameinProgress {

		round := NewRound(g.Team1, g.Team2, g.CurrentRound, g.is_T1_CT, &g.GameRules, g.OT, g)

		// Handle side switches and OT transitions
		if g.CurrentRound == g.GameRules.HalfLength+1 {
			// Regular halftime side switch
			round.HandleSideSwitch(g.Team1, g.Team2)
			g.is_T1_CT = !g.is_T1_CT
			g.firsthalf = !g.firsthalf
		} else if g.CurrentRound == (g.GameRules.HalfLength*2)+1 || (g.OT && g.CurrentRound == ((g.GameRules.HalfLength*2)+(g.OTcounter*g.GameRules.OTHalfLength*2)+1)) {
			// Start of overtime
			g.OT = true
			g.OTcounter++
			round.HandleOTStart(g.Team1, g.Team2)
		} else if g.OT && g.CurrentRound == (g.GameRules.HalfLength*2)+(g.OTcounter*g.GameRules.OTHalfLength)+1 {
			// OT halftime side switch
			round.HandleOTSideSwitch(g.Team1, g.Team2)
			g.is_T1_CT = !g.is_T1_CT
		}

		round.BuyPhase(g.Team1, g.Team2)

		round.RoundStart(g.Team1, g.Team2)

		round.RoundEnd(g.Team1, g.Team2)

		g.Rounds = append(g.Rounds, *round)

		g.UpdateScore(round.IsT1WinnerTeam)

		// Clear round pointer to help GC
		round = nil

		g.CurrentRound++
		g.GameFinished()

	}

}

func (g *Game) GameFinished() {
	if !g.OT {
		if g.Score[0] >= (g.GameRules.HalfLength+1) && g.Score[1] < (g.GameRules.HalfLength) {
			g.GameinProgress = false
			g.Is_T1_Winner = true // Team1 wins
		} else if g.Score[1] >= (g.GameRules.HalfLength+1) && g.Score[0] < (g.GameRules.HalfLength) {
			g.GameinProgress = false
			g.Is_T1_Winner = false // Team2 wins
		}
	} else if ((g.Score[0]-g.GameRules.HalfLength-(g.OTcounter*g.GameRules.OTHalfLength)) >= 1 || (g.Score[1]-g.GameRules.HalfLength-(g.OTcounter*g.GameRules.OTHalfLength)) >= 1) && math.Abs(float64(g.Score[0]-g.Score[1])) >= 2 {
		g.GameinProgress = false
		if g.Score[0] > g.Score[1] {
			g.Is_T1_Winner = true // Team1 wins
		} else {
			g.Is_T1_Winner = false // Team2 wins
		}
		//CAUTION WITH THE NEXT PART, THIS DEFINES THE MAXIMUM NUMBER OF ROUNDS IN OVERTIME
	} else if g.OTcounter > 50 {
		// If the game has gone on for too long, end it
		g.GameinProgress = false
		g.Is_T1_Winner = g.rng.Intn(2) == 0 // Randomly decide a winner
	}
}

func (g *Game) UpdateScore(winner bool) {
	g.Is_T1_Winner = winner
	if winner {
		g.Score[0]++
	} else {
		g.Score[1]++
	}
}

// Cleanup clears memory-intensive data structures after game completion
func (g *Game) Cleanup() {
	if g.Team1 != nil {
		g.Team1.Cleanup()
	}
	if g.Team2 != nil {
		g.Team2.Cleanup()
	}
	// Don't clear Rounds if exporting, but can be cleared after export
	// g.Rounds = nil
}

func (g *Game) GetPreviousRoundEndReason() int {
	if len(g.Rounds) == 0 {
		return 0 // or some default value indicating no previous round
	}
	return g.Rounds[len(g.Rounds)-1].Calc_Outcome.ReasonCode
}

func (g *Game) GetPreviousBombPlant() bool {
	if len(g.Rounds) == 0 {
		return false // or some default value indicating no previous round
	}
	return g.Rounds[len(g.Rounds)-1].Calc_Outcome.BombPlanted
}
