package engine

import (
	"math/rand"
)

type Game struct {
	ID             string
	GameinProgress bool
	Score          [2]int
	Rounds         []Round
	CurrentCT      bool // True if Team1 is CT
	OT             bool
	OTcounter      int
	firsthalf      bool
	sideswitch     bool
	CurrentRound   int
	GameRules      GameRules
	WinnerTeam     bool // false if Team1 wins, true if Team2 wins
	Team1          *Team
	Team2          *Team
}

func NewGame(id string, Team1Name string, Team1Strategy string, Team2Name string, Team2Strategy string, gamerules_path string) *Game {

	gameRules := NewGameRules(gamerules_path)

	currentCT := rand.Intn(2) == 0

	Team_1 := NewTeam(Team1Name, gameRules.StartingFunds, currentCT, gameRules.DefaultEquipment, Team1Strategy)
	Team_2 := NewTeam(Team2Name, gameRules.StartingFunds, !currentCT, gameRules.DefaultEquipment, Team2Strategy)

	return &Game{
		ID:             id,
		Team1:          Team_1,
		Team2:          Team_2,
		CurrentCT:      currentCT,
		OT:             false,
		OTcounter:      0,
		firsthalf:      true,
		sideswitch:     false,
		Score:          [2]int{0, 0},
		CurrentRound:   0,
		GameRules:      gameRules,
		GameinProgress: false,
	}

}

func (g *Game) Start() {
	g.GameinProgress = true

	for g.GameinProgress {
		g.CurrentRound++
		g.GameState()

		round := NewRound(g.Team1, g.Team2, g.CurrentRound, g.CurrentCT, g.sideswitch, &g.GameRules)

		round.BuyPhase()

		round.RoundStart()

		// Create a value-copy of the round for storage
		roundCopy := g.createRoundValueCopy(*round)
		g.Rounds = append(g.Rounds, roundCopy)

		g.UpdateScore(round.WinnerTeam)

		g.sideswitch = false
		g.GameFinished()

	}

}

func (g *Game) GameState() {
	if g.CurrentRound == 15 {
		g.switchSide()
	}

	if g.CurrentRound == 31 || (g.OT && g.CurrentRound == (15+(g.OTcounter*6)+1)) {
		g.OT = true
		g.OTcounter++
		g.Team1.NewOT(g.GameRules.OTFunds)
		g.Team2.NewOT(g.GameRules.OTFunds)
	}

	if g.OT && g.CurrentRound == (15+(g.OTcounter*3)+1) {
		g.switchSide()
	}

}

func (g *Game) switchSide() {
	g.sideswitch = true
	g.firsthalf = !g.firsthalf
	g.CurrentCT = !g.CurrentCT
	g.Team1.Sideswitch(g.OT, g.GameRules.StartingFunds, g.GameRules.OTFunds)
	g.Team2.Sideswitch(g.OT, g.GameRules.StartingFunds, g.GameRules.OTFunds)
}

func (g *Game) GameFinished() {
	if !g.OT {
		if g.Score[0] == 16 && g.Score[1] < 15 {
			g.GameinProgress = false
			g.WinnerTeam = false // Team1 wins
		} else if g.Score[1] == 16 && g.Score[0] < 15 {
			g.GameinProgress = false
			g.WinnerTeam = true // Team2 wins
		}
	} else if ((g.Score[0]-16)/g.OTcounter) > 4 && ((g.Score[1]-16)/g.OTcounter) > 4 && (g.Score[0]-g.Score[1]) > 2 {
		g.GameinProgress = false
		if g.Score[0] > g.Score[1] {
			g.WinnerTeam = false // Team1 wins
		} else {
			g.WinnerTeam = true // Team2 wins
		}
	}
}

func (g *Game) UpdateScore(winner bool) {
	g.WinnerTeam = winner
	if winner {
		g.Score[0]++
	} else {
		g.Score[1]++
	}
}

// createRoundValueCopy creates a deep copy of a Round with value-only fields (no pointers)
func (g *Game) createRoundValueCopy(r Round) Round {
	// Create a new round with all the primitive values copied
	roundCopy := Round{
		RoundNumber:   r.RoundNumber,
		CTTeam:        r.CTTeam,
		WinnerTeam:    r.WinnerTeam,
		WinnerSide:    r.WinnerSide,
		BombPlanted:   r.BombPlanted,
		Sideswitch:    r.Sideswitch,
		Bombplanted:   r.Bombplanted,
		SurvivingT1:   r.SurvivingT1,
		SurvivingT2:   r.SurvivingT2,
		EquipmentT1:   r.EquipmentT1,
		EquipmentT2:   r.EquipmentT2,
		FundsearnedT1: r.FundsearnedT1,
		FundsearnedT2: r.FundsearnedT2,
	}

	// Create deep copies of Team objects if they exist
	if r.Team1 != nil {
		team1Copy := Team{
			Name:                r.Team1.Name,
			Side:                r.Team1.Side,
			Funds:               r.Team1.Funds,
			Earned:              r.Team1.Earned,
			Equipment:           r.Team1.Equipment,
			Remainingequipement: r.Team1.Remainingequipement,
			Playersalive:        r.Team1.Playersalive,
			Score:               r.Team1.Score,
			Consecutiveloss:     r.Team1.Consecutiveloss,
			Spent:               r.Team1.Spent,
			Strategy:            r.Team1.Strategy,
		}
		roundCopy.Team1 = &team1Copy
	}

	if r.Team2 != nil {
		team2Copy := Team{
			Name:                r.Team2.Name,
			Side:                r.Team2.Side,
			Funds:               r.Team2.Funds,
			Earned:              r.Team2.Earned,
			Equipment:           r.Team2.Equipment,
			Remainingequipement: r.Team2.Remainingequipement,
			Playersalive:        r.Team2.Playersalive,
			Score:               r.Team2.Score,
			Consecutiveloss:     r.Team2.Consecutiveloss,
			Spent:               r.Team2.Spent,
			Strategy:            r.Team2.Strategy,
		}
		roundCopy.Team2 = &team2Copy
	}

	return roundCopy
}
