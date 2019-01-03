package tournaments

import (
	"sort"
	"time"

	"github.com/m4rw3r/uuid"
)

type PlayerResult struct {
	Place      int       `json:"place"`
	When       time.Time `json:"when"`
	NumPlayers int       `json:"numPlayers"`
}

type PlayerResults []PlayerResult

func (s PlayerResults) Len() int      { return len(s) }
func (s PlayerResults) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s PlayerResults) Less(i, j int) bool {
	if s[i].Place < s[j].Place {
		return true
	} else if s[i].Place > s[j].Place {
		return false
	} else {
		return s[i].NumPlayers > s[j].NumPlayers
	}
}

func (s PlayerResults) BetterThan(t PlayerResults) bool {

	if len(t) == 0 {
		return true
	}

	sort.Stable(s)
	sort.Stable(t)

	maxResults := len(t)
	if len(t) > len(s) {
		maxResults = len(s)
	}

	for i := 0; i < maxResults; i++ {
		if s[i].Place < t[i].Place {
			return true
		}
		if s[i].Place > t[i].Place {
			return false
		}
		if s[i].NumPlayers > t[i].NumPlayers {
			return true
		}
		if s[i].NumPlayers < t[i].NumPlayers {
			return false
		}

	}

	// If place and numPlayers still tie, it's "better" to have played more tournaments
	if len(s) > len(t) {
		return true
	}

	// t is best, or they are completely equal
	return false
}

func (s PlayerResults) Equals(t PlayerResults) bool {

	if len(t) != len(s) {
		return false
	}

	sort.Stable(s)
	sort.Stable(t)

	for i := 0; i < len(t); i++ {
		if s[i].Place != t[i].Place {
			return false
		}
		if s[i].NumPlayers > t[i].NumPlayers {
			return false
		}
	}

	return true
}

type PlayerStanding struct {
	Player     uuid.UUID     `json:"uuid"`
	Results    PlayerResults `json:"results"`
	Winnings   int           `json:"winnings"`
	AvgPlace   float64       `json:"avgPlace"`
	Points     int           `json:"points"`
	NumHeadsUp int           `json:"headsUp"`
	NumWins    int           `json:"wins"`
	NumPlayed  int           `json:"played"`
	Enough     bool          `json:"playedEnough"`
	NumTotal   int           `json:"numTotal"`
	Knockouts  int           `json:"knockouts"`
}

func (s *PlayerStanding) Equals(t *PlayerStanding) bool {
	if s == t {
		return true
	}
	if !s.Results.Equals(t.Results) {
		return false
	}
	if s.Winnings != t.Winnings || s.AvgPlace != t.AvgPlace || s.Points != t.Points || s.NumHeadsUp != t.NumHeadsUp || s.NumWins != t.NumWins || s.NumPlayed != t.NumPlayed || s.NumTotal != t.NumTotal || s.Knockouts != t.Knockouts {
		return false
	}
	return true
}

type PlayerStandings []*PlayerStanding

func (s PlayerStandings) Duplicate() PlayerStandings {
	return s.Combine(PlayerStandings{})
}

func (existingPS PlayerStandings) Combine(newPS PlayerStandings) PlayerStandings {
	var combined PlayerStandings
	for _, nps := range newPS {
		foundInOld := false
		for _, ops := range existingPS {
			if nps.Player == ops.Player {
				foundInOld = true
				cps := new(PlayerStanding)
				cps.Player = ops.Player
				cps.Results = append(ops.Results, nps.Results...)
				cps.Winnings = ops.Winnings + nps.Winnings
				cps.Points = ops.Points + nps.Points
				cps.NumHeadsUp = ops.NumHeadsUp + nps.NumHeadsUp
				cps.NumWins = ops.NumWins + nps.NumWins
				cps.NumPlayed = ops.NumPlayed + nps.NumPlayed
				cps.Enough = cps.NumPlayed > 10
				cps.NumTotal = ops.NumTotal + nps.NumTotal
				cps.Knockouts = ops.Knockouts + nps.Knockouts
				cps.AvgPlace = ((ops.AvgPlace * float64(ops.NumPlayed)) + (nps.AvgPlace * float64(nps.NumPlayed))) / float64(cps.NumPlayed)

				combined = append(combined, cps)
			}
		}
		if !foundInOld {
			combined = append(combined, nps)
		}
	}

	for _, ops := range existingPS {
		foundInCombined := false
		for _, cps := range combined {
			if ops.Player == cps.Player {
				foundInCombined = true
			}
		}
		if !foundInCombined {
			combined = append(combined, ops)
		}
	}

	return combined
}

func (s PlayerStandings) Len() int      { return len(s) }
func (s PlayerStandings) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type SortedStandings struct {
	ByWinnings      PlayerStandings `json:"byWinnings"`
	ByAvgPlace      PlayerStandings `json:"byAvgPlace"`
	ByPoints        PlayerStandings `json:"byPoints"`
	ByHeadsUp       PlayerStandings `json:"byHeadsUp"`
	ByWinRatio      PlayerStandings `json:"byWinRatio"`
	ByWinRatioTotal PlayerStandings `json:"byWinRatioTotal"`
	ByNumPlayed     PlayerStandings `json:"byNumPlayed"`
	ByKnockouts     PlayerStandings `json:"byKnockouts"`
}

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
		if s.PlayerStandings[i].AvgPlace > s.PlayerStandings[j].AvgPlace {
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

type ByBestPlayer struct{ PlayerStandings }

func (s ByBestPlayer) Less(i, j int) bool {
	a, b := s.PlayerStandings[i], s.PlayerStandings[j]
	if a.AvgPlace < b.AvgPlace {
		return true
	}

	if a.AvgPlace == b.AvgPlace {
		if a.Results.BetterThan(b.Results) {
			return true
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

type ByWinRatio struct{ PlayerStandings }

func (s ByWinRatio) Less(i, j int) bool {
	a := float64(s.PlayerStandings[i].NumWins) / float64(s.PlayerStandings[i].NumPlayed)
	b := float64(s.PlayerStandings[j].NumWins) / float64(s.PlayerStandings[j].NumPlayed)

	if a > b {
		return true
	}

	if a == b {
		if s.PlayerStandings[i].NumWins > s.PlayerStandings[j].NumWins {
			return true
		}
	}
	return false
}

type ByWinRatioTotal struct{ PlayerStandings }

func (s ByWinRatioTotal) Less(i, j int) bool {
	a := float64(s.PlayerStandings[i].NumWins) / float64(s.PlayerStandings[i].NumTotal)
	b := float64(s.PlayerStandings[j].NumWins) / float64(s.PlayerStandings[j].NumTotal)

	if a > b {
		return true
	}

	if a == b {
		if s.PlayerStandings[i].NumWins > s.PlayerStandings[j].NumWins {
			return true
		}
	}
	return false
}

type ByNumPlayed struct{ PlayerStandings }

func (s ByNumPlayed) Less(i, j int) bool {
	aP := float64(s.PlayerStandings[i].NumPlayed)
	bP := float64(s.PlayerStandings[j].NumPlayed)
	aT := float64(s.PlayerStandings[i].NumTotal)
	bT := float64(s.PlayerStandings[j].NumTotal)

	aR := (aT - aP) / aT
	bR := (bT - bP) / bT

	if aR < bR {
		return true
	}

	if aR == bR {
		if s.PlayerStandings[i].NumWins > s.PlayerStandings[j].NumWins {
			return true
		}
	}
	return false
}

type ByKnockouts struct{ PlayerStandings }

func (s ByKnockouts) Less(i, j int) bool {
	if s.PlayerStandings[i].Knockouts > s.PlayerStandings[j].Knockouts {
		return true
	}

	if s.PlayerStandings[i].Knockouts == s.PlayerStandings[j].Knockouts {
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

// TODO: Split into smaller functions
func NewStandings(tournaments Tournaments) PlayerStandings {

	// First, find all active players for these tournaments
	// Also, get the max number of players for a given set of tournaments
	// This gives us the basis for low point scoring
	activePlayers, maxPlayers := getActivePlayers(tournaments)

	// Then, loop through tournaments again to keep track of relevant stats
	winnings := make(map[uuid.UUID]int)
	sumPlace := make(map[uuid.UUID]int)
	points := make(map[uuid.UUID][]int)
	numHeadsUp := make(map[uuid.UUID]int)
	numWins := make(map[uuid.UUID]int)
	numPlayed := make(map[uuid.UUID]int)
	knockouts := make(map[uuid.UUID]int)
	results := make(map[uuid.UUID][]PlayerResult)
	numTotal := 0

	for _, t := range tournaments {
		if !t.Played || len(t.Result) == 0 {
			continue
		}
		numTotal += 1
		seenPlayer := make(map[uuid.UUID]bool)
		for i, player := range t.Result {
			place := i + 1
			results[player] = append(results[player], PlayerResult{
				Place:      place,
				When:       t.Info.Scheduled,
				NumPlayers: len(t.Result),
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

		for player := range t.BountyHunters {
			knockouts[player] += len(t.BountyHunters[player])
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
		if numPlayed[player] >= 11 {
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
			NumTotal:   numTotal,
			Knockouts:  knockouts[player],
		})
	}

	return standings
}

func TotalStandings(seasons []int) *SortedStandings {

	sortedStandings := new(SortedStandings)

	var totalStandings PlayerStandings
	for _, season := range seasons {
		tList, err := TournamentsBySeason(season)
		if err != nil {
			// TODO
		}
		standings := NewStandings(tList)
		totalStandings = totalStandings.Combine(standings)
	}

	// Got to use new rules to tie break here..
	totalStandings.ByWinnings(false)
	sortedStandings.ByWinnings = totalStandings

	totalStandings = totalStandings.Duplicate()
	totalStandings.ByAvgPlace()
	sortedStandings.ByAvgPlace = totalStandings

	totalStandings = totalStandings.Duplicate()
	totalStandings.ByPoints()
	sortedStandings.ByPoints = totalStandings

	totalStandings = totalStandings.Duplicate()
	totalStandings.ByHeadsUp()
	sortedStandings.ByHeadsUp = totalStandings

	totalStandings = totalStandings.Duplicate()
	totalStandings.ByWinRatio()
	sortedStandings.ByWinRatio = totalStandings

	totalStandings = totalStandings.Duplicate()
	totalStandings.ByWinRatioTotal()
	sortedStandings.ByWinRatioTotal = totalStandings

	totalStandings = totalStandings.Duplicate()
	totalStandings.ByNumPlayed()
	sortedStandings.ByNumPlayed = totalStandings

	totalStandings = totalStandings.Duplicate()
	totalStandings.ByKnockouts()
	sortedStandings.ByKnockouts = totalStandings

	return sortedStandings

}

func SeasonStandings(season int) *SortedStandings {

	tList, err := TournamentsBySeason(season)
	if err != nil {
		// TODO
	}

	// Compute standings
	sortedStandings := new(SortedStandings)

	standings := NewStandings(tList)
	standings.ByWinnings(season < 2013)
	sortedStandings.ByWinnings = standings

	standings = standings.Duplicate()
	standings.ByAvgPlace()
	sortedStandings.ByAvgPlace = standings

	standings = standings.Duplicate()
	standings.ByPoints()
	sortedStandings.ByPoints = standings

	standings = standings.Duplicate()
	standings.ByHeadsUp()
	sortedStandings.ByHeadsUp = standings

	standings = standings.Duplicate()
	standings.ByWinRatio()
	sortedStandings.ByWinRatio = standings

	standings = standings.Duplicate()
	standings.ByWinRatioTotal()
	sortedStandings.ByWinRatioTotal = standings

	standings = standings.Duplicate()
	standings.ByNumPlayed()
	sortedStandings.ByNumPlayed = standings

	standings = standings.Duplicate()
	standings.ByKnockouts()
	sortedStandings.ByKnockouts = standings

	return sortedStandings
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
	sort.Stable(ByAvgPlace{s})
}

func (s PlayerStandings) ByPoints() {
	sort.Stable(ByPoints{s})
}

func (s PlayerStandings) ByHeadsUp() {
	sort.Stable(sort.Reverse(ByHeadsUp{s}))
}

func (s PlayerStandings) ByBestPlayer() {
	sort.Stable(ByBestPlayer{s})
}

func (s PlayerStandings) ByWinRatio() {
	sort.Stable(ByWinRatio{s})
}

func (s PlayerStandings) ByWinRatioTotal() {
	sort.Stable(ByWinRatioTotal{s})
}

func (s PlayerStandings) ByNumPlayed() {
	sort.Stable(ByNumPlayed{s})
}

func (s PlayerStandings) ByWorstPlayer() {
	sort.Stable(sort.Reverse(ByBestPlayer{s}))
}

func (s PlayerStandings) ByKnockouts() {
	sort.Stable(ByKnockouts{s})
}
