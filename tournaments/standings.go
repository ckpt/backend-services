package tournaments

import (
	"github.com/m4rw3r/uuid"
	"sort"
	"time"
)

type YellowPeriod struct {
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
	Player uuid.UUID `json:"uuid"`
	Active bool      `json:"active"`
}

type PlayerResult struct {
	Place int       `json:"place"`
	When  time.Time `json:"when"`
}

type PlayerStanding struct {
	Player     uuid.UUID      `json:"uuid"`
	Results    []PlayerResult `json:"results"`
	Winnings   int            `json:"winnings"`
	AvgPlace   float64        `json:"avgPlace"`
	Points     int            `json:"points"`
	NumHeadsUp int            `json:"headsUp"`
	NumWins    int            `json:"wins"`
	NumPlayed  int            `json:"played"`
	Enough     bool           `json:"playedEnough"`
	NumTotal   int            `json:"numTotal"`
}

type PlayerStandings []*PlayerStanding

func (s PlayerStandings) Len() int      { return len(s) }
func (s PlayerStandings) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByWinnings struct{ PlayerStandings }

func (s ByWinnings) Less(i, j int) bool {
	if s.PlayerStandings[i].Winnings < s.PlayerStandings[j].Winnings {
		return true
	}

	if s.PlayerStandings[i].Winnings == s.PlayerStandings[j].Winnings {
		if s.PlayerStandings[i].Points > s.PlayerStandings[j].Points {
			return true
		}

		if s.PlayerStandings[i].Points == s.PlayerStandings[j].Points {
			if s.PlayerStandings[i].NumWins < s.PlayerStandings[j].NumWins {
				return true
			}
		}
	}
	return false
}

// Implements old tie break used in earlier seasons
type ByWinningsOld struct{ PlayerStandings }

func (s ByWinningsOld) Less(i, j int) bool {
	if s.PlayerStandings[i].Winnings < s.PlayerStandings[j].Winnings {
		return true
	}

	if s.PlayerStandings[i].Winnings == s.PlayerStandings[j].Winnings {
		if s.PlayerStandings[i].AvgPlace < s.PlayerStandings[j].AvgPlace {
			return true
		}

		if s.PlayerStandings[i].Points == s.PlayerStandings[j].Points {
			if s.PlayerStandings[i].NumWins < s.PlayerStandings[j].NumWins {
				return true
			}
		}
	}
	return false
}

type ByAvgPlace struct{ PlayerStandings }

func (s ByAvgPlace) Less(i, j int) bool {
	if s.PlayerStandings[i].AvgPlace < s.PlayerStandings[j].AvgPlace {
		return true
	}

	if s.PlayerStandings[i].AvgPlace == s.PlayerStandings[j].AvgPlace {
		if s.PlayerStandings[i].Winnings > s.PlayerStandings[j].Winnings {
			return true
		}

		if s.PlayerStandings[i].Winnings == s.PlayerStandings[j].Winnings {
			if s.PlayerStandings[i].NumWins > s.PlayerStandings[j].NumWins {
				return true
			}
		}
	}
	return false
}

type ByPoints struct{ PlayerStandings }

func (s ByPoints) Less(i, j int) bool {
	if s.PlayerStandings[i].Points < s.PlayerStandings[j].Points {
		return true
	}

	if s.PlayerStandings[i].Points == s.PlayerStandings[j].Points {
		if s.PlayerStandings[i].Winnings > s.PlayerStandings[j].Winnings {
			return true
		}

		if s.PlayerStandings[i].Winnings == s.PlayerStandings[j].Winnings {
			if s.PlayerStandings[i].NumWins > s.PlayerStandings[j].NumWins {
				return true
			}
		}
	}
	return false
}

type ByHeadsUp struct{ PlayerStandings }

func (s ByHeadsUp) Less(i, j int) bool {
	hur1 := float64(s.PlayerStandings[i].NumHeadsUp) / float64(s.PlayerStandings[i].NumPlayed)
	hur2 := float64(s.PlayerStandings[j].NumHeadsUp) / float64(s.PlayerStandings[j].NumPlayed)

	wr1 := float64(s.PlayerStandings[i].NumWins) / float64(s.PlayerStandings[i].NumHeadsUp)
	wr2 := float64(s.PlayerStandings[j].NumWins) / float64(s.PlayerStandings[j].NumHeadsUp)

	if hur1 < hur2 {
		return true
	}

	if hur1 == hur2 {
		if wr1 < wr2 {
			return true
		}

		if wr1 == wr2 {
			if s.PlayerStandings[i].Winnings < s.PlayerStandings[j].Winnings {
				return true
			}
		}
	}
	return false
}

func getActivePlayers(tournaments Tournaments) ([]uuid.UUID, int) {
	var activePlayers []uuid.UUID

	maxPlayers := 0
	seenPlayer := make(map[uuid.UUID]bool)
	for _, t := range tournaments {
		if !t.Played || len(t.Result) == 0 {
			continue
		}

		for _, player := range t.Result {
			seenPlayer[player] = true
		}

		if len(t.Result) > maxPlayers {
			maxPlayers = len(t.Result)
		}
	}

	for k, _ := range seenPlayer {
		activePlayers = append(activePlayers, k)
	}

	return activePlayers, maxPlayers
}

func YellowPeriods(tournaments Tournaments) []YellowPeriod {
	var periods []YellowPeriod
	var currentPeriod *YellowPeriod
	var season, seasonIndex int

	sort.Sort(tournaments)
	for i := range tournaments {
		if !tournaments[i].Played {
			continue
		}
		// Leader is based on results for the season, so start from "scratch" on new seasons
		if tournaments[i].Info.Season != season {
			season = tournaments[i].Info.Season
			seasonIndex = i
		}
		standings := NewStandings(tournaments[seasonIndex : i+1])
		standings.ByWinnings(season < 2013)
		if currentPeriod == nil {
			currentPeriod = &YellowPeriod{
				From:   tournaments[i].Info.Scheduled,
				To:     tournaments[i].Info.Scheduled,
				Player: standings[0].Player,
				Active: true,
			}
		} else if currentPeriod.Player == standings[0].Player {
			currentPeriod.To = tournaments[i].Info.Scheduled
		} else {
			currentPeriod.Active = false
			currentPeriod.To = tournaments[i].Info.Scheduled
			periods = append(periods, *currentPeriod)
			currentPeriod = &YellowPeriod{
				From:   tournaments[i].Info.Scheduled,
				To:     tournaments[i].Info.Scheduled,
				Player: standings[0].Player,
				Active: true,
			}
		}
	}
	periods = append(periods, *currentPeriod)
	return periods
}

func NewStandings(tournaments Tournaments) PlayerStandings {

	// First, find all active players for these tournaments
	// Also, get the max number of players for a given tournament
	// This gives us the basis for low point scoring
	activePlayers, maxPlayers := getActivePlayers(tournaments)

	// Then, loop through tournaments again to keep track of relevant stats
	winnings := make(map[uuid.UUID]int)
	sumPlace := make(map[uuid.UUID]int)
	points := make(map[uuid.UUID][]int)
	numHeadsUp := make(map[uuid.UUID]int)
	numWins := make(map[uuid.UUID]int)
	numPlayed := make(map[uuid.UUID]int)
	results := make(map[uuid.UUID][]PlayerResult)

	for _, t := range tournaments {
		if !t.Played || len(t.Result) == 0 {
			continue
		}

		seenPlayer := make(map[uuid.UUID]bool)
		for i, player := range t.Result {
			place := i + 1
			results[player] = append(results[player], PlayerResult{
				Place: place,
				When:  t.Info.Scheduled,
			})

			sumPlace[player] += place
			numPlayed[player] += 1
			seenPlayer[player] = true
			points[player] = append(points[player], place)

			switch place {
			case 1:
				numWins[player] += 1
				winnings[player] += (len(t.Result) - 2) * t.Info.Stake
				numHeadsUp[player] += 1
			case 2:
				numHeadsUp[player] += 1
			default:
				winnings[player] -= t.Info.Stake
			}
		}

		for _, player := range activePlayers {
			if _, seen := seenPlayer[player]; !seen {
				points[player] = append(points[player], maxPlayers+1)
			}
		}
	}

	// Finally, loop through active players and set totals, returning standings
	var standings PlayerStandings

	for _, player := range activePlayers {

		// Remove worst point score at 10 and 20 played tournaments
		pp := points[player]
		sort.Ints(pp)

		if numPlayed[player] >= 10 {
			i := len(pp) - 1
			pp = append(pp[:i], pp[i+1:]...) // delete element at index i (last)
		}
		if numPlayed[player] >= 20 {
			i := len(pp) - 1
			pp = append(pp[:i], pp[i+1:]...) // delete element at index i (last)
		}

		// Now, sum up the points
		sumPoints := 0
		for _, p := range pp {
			sumPoints += p
		}

		// Check if player has enough tournaments
		// TODO: Is this a property of season?
		enough := false
		if numPlayed[player] >= 8 {
			enough = true
		}

		standings = append(standings, &PlayerStanding{
			Player:     player,
			Results:    results[player],
			Winnings:   winnings[player],
			AvgPlace:   float64(sumPlace[player]) / float64(numPlayed[player]),
			Points:     sumPoints,
			NumHeadsUp: numHeadsUp[player],
			NumWins:    numWins[player],
			NumPlayed:  numPlayed[player],
			Enough:     enough,
			NumTotal:   len(tournaments),
		})
	}

	return standings
}

// Various ways to sort the player standings using helper structs that
// implement different comparison methods.

func (s PlayerStandings) ByWinnings(oldTieBreak bool) {
	if oldTieBreak {
		sort.Sort(sort.Reverse(ByWinningsOld{s}))
	} else {
		sort.Sort(sort.Reverse(ByWinnings{s}))
	}
}

func (s PlayerStandings) ByAvgPlace() {
	sort.Sort(ByAvgPlace{s})
}

func (s PlayerStandings) ByPoints() {
	sort.Sort(ByPoints{s})
}

func (s PlayerStandings) ByHeadsUp() {
	sort.Sort(sort.Reverse(ByHeadsUp{s}))
}
