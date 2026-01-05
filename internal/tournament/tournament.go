package tournament

import (
	"context"
	"csgo_abm/internal/engine"
	"time"
)

type Format string

const (
	FormatRoundRobin Format = "roundrobin"
)

type SeriesSpec struct {
	NumGames       int
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

// RoundRobinSchedule pairs each strategy with every other,
// with both orderings for balance (currently not necessary as starting side is randomized)
func RoundRobinSchedule(strategies []string) []MatchSpec {
	var matches []MatchSpec
	for i := 0; i < len(strategies); i++ {
		for j := i + 1; j < len(strategies); j++ {
			// A vs B
			matches = append(matches, MatchSpec{
				Team1Name:     strategies[i],
				Team1Strategy: strategies[i],
				Team2Name:     strategies[j],
				Team2Strategy: strategies[j],
			})
			// B vs A
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

// RunMatchup executes many independent ABM games for a matchup to estimate performance
func RunMatchup(ctx context.Context, m MatchSpec, rules engine.GameRules, spec SeriesSpec) (SeriesResult, error) {
	res := SeriesResult{Match: m}
	if spec.NumGames <= 0 {
		spec.NumGames = 1000
	}
	if spec.MaxConcurrent <= 0 {
		spec.MaxConcurrent = 1
	}
	type item struct{ idx int }
	jobs := make(chan item)
	results := make(chan GameOutcome)
	// Workers
	for w := 0; w < spec.MaxConcurrent; w++ {
		go func() {
			for it := range jobs {
				// Allow cancel
				select {
				case <-ctx.Done():
					return
				default:
				}
				game := engine.NewGame("", m.Team1Name, m.Team1Strategy, m.Team2Name, m.Team2Strategy, rules)
				game.SetSeed(spec.Seed + int64(it.idx))
				done := make(chan struct{})
				go func() { game.Start(); close(done) }()
				if spec.TimeoutPerGame > 0 {
					timer := time.NewTimer(spec.TimeoutPerGame)
					select {
					case <-done:
						timer.Stop()
					case <-timer.C:
						// timeout -> treat as no result; skip
					}
				} else {
					<-done
				}
				results <- GameOutcome{T1Wins: game.Is_T1_Winner, Score: game.Score}
			}
		}()
	}
	// Feed jobs
	go func() {
		for i := 0; i < spec.NumGames; i++ {
			jobs <- item{idx: i}
		}
		close(jobs)
	}()
	// Collect
	for i := 0; i < spec.NumGames; i++ {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case out := <-results:
			res.GameResults = append(res.GameResults, out)
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
		for _, g := range sr.GameResults {
			if g.T1Wins {
				rows[i1].Wins++
				rows[i2].Losses++
				rows[i1].MapWins++
				rows[i2].MapLoss++
			} else {
				rows[i2].Wins++
				rows[i1].Losses++
				rows[i2].MapWins++
				rows[i1].MapLoss++
			}
		}
	}
	return Standings{Rows: rows}
}

// GetRows exposes rows for simple CSV exporting without introducing new types in analysis
func (s Standings) GetRows() []StandingsRow { return s.Rows }
