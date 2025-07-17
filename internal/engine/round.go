package engine

import (
	"math/rand"
)

type Round struct {
	Team1       *Team
	Team2       *Team
	CTTeam      bool // false if Team1 is CT
	RoundNumber int
	WinnerTeam  bool //False if Team1 wins, true if Team2 wins
	WinnerSide  bool // false if CT wins, true if T wins
	BombPlanted bool
	SurvivingCT int
	SurvivingT  int
	Sideswitch  bool // True if sideswitch has occurred
	Bombplanted bool
	GameRules   *GameRules
}

func NewRound(T1 *Team, T2 *Team, roundNumber int, ctteam bool, sidewitch bool, gamerules *GameRules) *Round {

	T1.NewRound()
	T2.NewRound()

	return &Round{
		Team1:       T1,
		Team2:       T2,
		CTTeam:      ctteam,
		RoundNumber: roundNumber,
		WinnerTeam:  false,
		WinnerSide:  false,
		BombPlanted: false,
		SurvivingCT: 0,
		SurvivingT:  0,
		Sideswitch:  sidewitch,
		Bombplanted: false,
		GameRules:   gamerules,
	}
}

func (r *Round) BuyPhase() {

}

func (r *Round) RoundStart() {
	// Simulate the round logic here
	// For now, we will just randomly determine a winner
	r.determineWinner()

	r.determineBombplant()

	//round end logic
	r.RoundEnd()
}

func (r *Round) RoundEnd() {

	//determine surviving players

	//determine remaining equipment

	//determine funds earned

}

func (r *Round) determineWSurvivors() {
	// Placeholder for CSF logic
	ctteam_equipment := r.Team1.Equipment
	tteam_equipment := r.Team2.Equipment

	if r.CTTeam {
		ctteam_equipment = r.Team1.Equipment
		tteam_equipment = r.Team2.Equipment
	}

}

func (r *Round) determineBombplant() {
	// Placeholder for CSF logic

	ctteam_equipment := r.Team1.Equipment
	tteam_equipment := r.Team2.Equipment

	if r.CTTeam {
		ctteam_equipment = r.Team1.Equipment
		tteam_equipment = r.Team2.Equipment
	}

	probability := ContestSuccessFunction_simple(tteam_equipment, ctteam_equipment, r.GameRules.CSF_r)

	// Determine if the bomb was planted
	if rand.Float64() < probability {
		r.BombPlanted = true
	}
}

func (r *Round) determineWinner() {
	// Placeholder for CSF logic
	team1equipment := r.Team1.Equipment
	team2equipment := r.Team2.Equipment

	probability := ContestSuccessFunction_simple(team1equipment, team2equipment, r.GameRules.CSF_r)

	if rand.Float64() < probability {
		r.WinnerTeam = false // Team1 wins
	} else {
		r.WinnerTeam = true // Team2 wins
	}

	if r.WinnerTeam == r.CTTeam {
		r.WinnerSide = false // CT wins
	} else {
		r.WinnerSide = true // T wins
	}

}
