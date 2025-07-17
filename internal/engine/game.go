package engine

import (
	"csgo-economy-sim/internal/models"
	"math/rand"
)

type Game struct {
	Team1          *Team
	Team2          *Team
	CurrentCT      bool // True if Team1 is CT
	OT             bool
	OTcounter      int
	firsthalf      bool
	sideswitch     bool
	Rounds         []*Round
	Score          [2]int
	CurrentRound   int
	GameinProgress bool
	GameRules      GameRules
	WinnerTeam     bool // false if Team1 wins, true if Team2 wins
}

func NewGame(Team1Name string, Team1Strategy string, Team2Name string, Team2Strategy string, gamerules string) *Game {

	gameRules := NewGameRules(gamerules)

	currentCT := rand.Intn(2) == 0

	Team_1 := NewTeam(Team1Name, gameRules.startingFunds, currentCT, gameRules.defaultEquipment, models.Strategy{Name: Team1Strategy})
	Team_2 := NewTeam(Team2Name, gameRules.startingFunds, !currentCT, gameRules.defaultEquipment, models.Strategy{Name: Team2Strategy})

	return &Game{
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

		round.RoundStart()
		g.Rounds = append(g.Rounds, round)
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
		g.Team1.NewOT(g.GameRules.otFunds)
		g.Team2.NewOT(g.GameRules.otFunds)
	}

	if g.OT && g.CurrentRound == (15+(g.OTcounter*3)+1) {
		g.switchSide()
	}

}

func (g *Game) switchSide() {
	g.sideswitch = true
	g.firsthalf = !g.firsthalf
	g.CurrentCT = !g.CurrentCT
	g.Team1.Sideswitch(g.OT, g.GameRules.startingFunds, g.GameRules.otFunds)
	g.Team2.Sideswitch(g.OT, g.GameRules.startingFunds, g.GameRules.otFunds)
}

func (g *Game) GameFinished() {
	if g.Score[0] == 16 && g.Score[1] < 15 {
		g.GameinProgress = false
		g.WinnerTeam = false // Team1 wins
	} else if g.Score[1] == 16 && g.Score[0] < 15 {
		g.GameinProgress = false
		g.WinnerTeam = true // Team2 wins
	} else if g.CurrentRound >= 30 && ((g.Score[0]-16)/g.OTcounter) > 4 && ((g.Score[1]-16)/g.OTcounter) > 4 && (g.Score[0]-g.Score[1]) > 2 {
		g.GameinProgress = false
		if g.Score[0] > g.Score[1] {
			g.WinnerTeam = false // Team1 wins
		} else {
			g.WinnerTeam = true // Team2 wins
		}
	}
}

func (g *Game) UpdateScore(winner bool) {
	if winner == true {
		g.Score[0]++
	} else {
		g.Score[1]++
	}
}
