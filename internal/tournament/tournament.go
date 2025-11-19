package tournament

import (
	"context"
	"csgo_abm/internal/engine"
	"fmt"
	"math/rand"
	"time"
)

type Format string

const (
	FormatRoundRobin Format = "roundrobin"
)

type SeriesSpec struct {
	BestOf         int
	Seed           int64
	MaxConcurrent  int
	TimeoutPerGame time.Duration
}

type MatchSpec struct {
	Team1Name     string
	Team1Strategy string
	Team2Name     string
	Team2Strategy string
}

type TournamentConfig struct {
	Format       Format
	Series       SeriesSpec
	Participants []string // strategy names
}

type GameOutcome struct {
	T1Wins bool
	Score  [2]int
}

type SeriesResult struct {
	Match       MatchSpec
	SeriesWins  [2]int
	GameResults []GameOutcome
}

type StandingsRow struct {
	Strategy string
	Wins     int
	Losses   int
	MapWins  int
	MapLoss  int
}

type Standings struct {
	Rows []StandingsRow
}

// RoundRobinSchedule pairs each strategy with every other, mirrored for side fairness
func RoundRobinSchedule(strategies []string) []MatchSpec {
	var matches []MatchSpec
	for i := 0; i < len(strategies); i++ {
		for j := i + 1; j < len(strategies); j++ {
			matches = append(matches, MatchSpec{
				Team1Name:     strategies[i],
				Team1Strategy: strategies[i],
				Team2Name:     strategies[j],
				Team2Strategy: strategies[j],
			})
			matches = append(matches, MatchSpec{
				Team1Name:     strategies[j],
				Team1Strategy: strategies[j],
				Team2Name:     strategies[i],
				Team2Strategy: strategies[i],
			})
		}
	}
	return matches
}

// RunSeries executes a best-of-N series between two strategies
func RunSeries(ctx context.Context, m MatchSpec, rules engine.GameRules, spec SeriesSpec) (SeriesResult, error) {
	res := SeriesResult{Match: m}
	if spec.BestOf <= 0 {
		spec.BestOf = 3
	}
	needed := spec.BestOf/2 + 1
	r := rand.New(rand.NewSource(spec.Seed))
	for g := 0; g < spec.BestOf; g++ {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		default:
		}
		game := engine.NewGame("", m.Team1Name, m.Team1Strategy, m.Team2Name, m.Team2Strategy, rules)
		// Seed each game deterministically from series seed
		game.SetSeed(spec.Seed + int64(g) + r.Int63())
		done := make(chan struct{})
		go func() {
			game.Start()
			close(done)
		}()
		if spec.TimeoutPerGame > 0 {
			timer := time.NewTimer(spec.TimeoutPerGame)
			select {
			case <-done:
				timer.Stop()
			case <-timer.C:
				return res, fmt.Errorf("game timeout in series %s vs %s", m.Team1Name, m.Team2Name)
			}
		} else {
			<-done
		}
		out := GameOutcome{T1Wins: game.Is_T1_Winner, Score: game.Score}
		res.GameResults = append(res.GameResults, out)
		if out.T1Wins {
			res.SeriesWins[0]++
		} else {
			res.SeriesWins[1]++
		}
		if res.SeriesWins[0] == needed || res.SeriesWins[1] == needed {
			break
		}
	}
	return res, nil
}

// ComputeStandings aggregates series results into a table
func ComputeStandings(strategies []string, series []SeriesResult) Standings {
	idx := map[string]int{}
	rows := make([]StandingsRow, 0, len(strategies))
	for _, s := range strategies {
		idx[s] = len(rows)
		rows = append(rows, StandingsRow{Strategy: s})
	}
	for _, sr := range series {
		i1 := idx[sr.Match.Team1Strategy]
		i2 := idx[sr.Match.Team2Strategy]
		if sr.SeriesWins[0] > sr.SeriesWins[1] {
			rows[i1].Wins++
			rows[i2].Losses++
		} else {
			rows[i2].Wins++
			rows[i1].Losses++
		}
		for _, g := range sr.GameResults {
			if g.T1Wins {
				rows[i1].MapWins++
				rows[i2].MapLoss++
			} else {
				rows[i2].MapWins++
				rows[i1].MapLoss++
			}
		}
	}
	return Standings{Rows: rows}
}

// GetRows exposes rows for simple CSV exporting without introducing new types in analysis
func (s Standings) GetRows() []StandingsRow { return s.Rows }
