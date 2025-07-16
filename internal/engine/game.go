package engine

import (
	"csgo-economy-sim/internal/models"
)

type Game struct {
	Team1        *Team
	Team2        *Team
	CurrentCT    bool // True if Team1 is CT
	OT           bool
	firsthalf    bool
	Rounds       []Round
	Score        [2]int
	CurrentRound int
	GameRules    GameRules
}

func SetupGame(Team1Name string, Team1Strategy string, Team2Name string, Team2Strategy string, gamerules string) *Game {

	ctTeam := NewTeam(Team1Name, 800, true, 800, models.Strategy{Name: Team1Strategy})
	tTeam := NewTeam(Team2Name, 800, false, 800, models.Strategy{Name: Team2Strategy})

	gameRules := initGameRules(gamerules)

	return &Game{
		Team1:     ctTeam,
		Team2:     tTeam,
		CurrentCT: true,
		OT:        false,
		firsthalf: true,
		Rounds:    []Round{},
		Score:     [2]int{0, 0},
		GameRules: gameRules,
	}

}

func NewGame() *Game {

	return &Game{
		CTTeam:       ctTeam,
		TTeam:        tTeam,
		Rounds:       []Round{},
		Score:        [2]int{0, 0},
		CurrentRound: 0,
	}
}

func (g *Game) Start() {
	for g.CurrentRound < 30 {
		round := NewRound(g.CTTeam, g.TTeam)
		round.Play()
		g.Rounds = append(g.Rounds, round)
		g.UpdateScore(round.Winner)
		g.CurrentRound++

		if g.CurrentRound == 15 {
			g.CTTeam, g.TTeam = g.TTeam, g.CTTeam // Swap sides at halftime
		}
	}

	if g.Score[0] == g.Score[1] {
		g.StartOvertime()
	}
}

func (g *Game) UpdateScore(winner *Team) {
	if winner == g.CTTeam {
		g.Score[0]++
	} else {
		g.Score[1]++
	}
}

func (g *Game) StartOvertime() {
	for g.CurrentRound < 35 {
		round := NewRound(g.CTTeam, g.TTeam)
		round.Play()
		g.Rounds = append(g.Rounds, round)
		g.UpdateScore(round.Winner)
		g.CurrentRound++
	}
}

func (g *Game) GetResults() (int, int) {
	return g.Score[0], g.Score[1]
}
