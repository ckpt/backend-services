package tournaments

import (
	"errors"
	"github.com/imdario/mergo"
	"github.com/m4rw3r/uuid"
	"time"
)

// We use dummy in memory storage for now
var storage TournamentStorage = NewRedisTournamentStorage()

type Absentee struct {
	Player uuid.UUID `json:"player"`
	Reason string    `json:"reason"`
}

type Result []uuid.UUID

type Bet struct {
	Player     uuid.UUID `json:"player"`
	Prediction Result    `json:"prediction"`
}

type Info struct {
	Scheduled time.Time `json:"scheduled"`
	Stake     int       `json:"stake"`
	Location  uuid.UUID `json:"location"`
	Catering  uuid.UUID `json:"catering"`
	Season    int       `json:"season"`
}

type Tournament struct {
	UUID    uuid.UUID  `json:"uuid"`
	Info    Info       `json:"info"`
	Noshows []Absentee `json:"noshows"`
	Result  Result     `json:"result"`
	Played  bool       `json:"played"`
	Bets    []Bet      `json:"bets"`
}

type Tournaments []*Tournament

func (t Tournaments) Len() int           { return len(t) }
func (t Tournaments) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Tournaments) Less(i, j int) bool { return t[i].Info.Scheduled.Before(t[j].Info.Scheduled) }


// A storage interface for Tournaments
type TournamentStorage interface {
	Store(*Tournament) error
	Delete(uuid.UUID) error
	Load(uuid.UUID) (*Tournament, error)
	LoadAll() (Tournaments, error)
	LoadBySeason(int) (Tournaments, error)
}

//
// Tournaments related functions and methods
//

func (t Tournaments) GroupByMonths(season int) map[time.Month]Tournaments {
	byMonth := make(map[time.Month]Tournaments)

	for _,entry := range(t) {
		y, m := entry.Info.Scheduled.Year(), entry.Info.Scheduled.Month()
		if y == season {
			byMonth[m] = append(byMonth[m], entry)
		}
	}

	return byMonth
}

func (t Tournaments) Seasons() []int  {
	seenSeasons := make(map[int]bool)

	for _,entry := range(t) {
		seenSeasons[entry.Info.Season] = true
	}

	var seasons []int
	for k := range seenSeasons {
		seasons = append(seasons, k)
	}
	return seasons
}

func (t Tournaments) Played() Tournaments {
	var played Tournaments

	for _,entry := range(t) {
		if entry.Played {
			played = append(played, entry)
		}
	}

	return played
}


//
// Tournament related functions and methods
//

// Helpers
func validateTournamentInfo(info Info) error {
	if info.Scheduled.IsZero() {
		return errors.New("Tournament needs scheduled date")
	}
	if info.Stake == 0 {
		return errors.New("Tournament needs a stake")
	}
	if info.Season == 0 {
		return errors.New("Tournament needs a season")
	}
	return nil
}

func fixupTournamentInfo(oldinfo *Info, newinfo Info) {
	if oldinfo.Scheduled != newinfo.Scheduled {
		oldinfo.Scheduled = newinfo.Scheduled
	}
}

// Create a Tournament
func NewTournament(tdata Info) (*Tournament, error) {
	if err := validateTournamentInfo(tdata); err != nil {
		return nil, errors.New(err.Error() + " - Could not create tournament")
	}

	t := new(Tournament)
	t.UUID, _ = uuid.V4()
	if err := mergo.MergeWithOverwrite(&t.Info, tdata); err != nil {
		return nil, errors.New(err.Error() + " - Could not set initial tournament data")
	}
	// Merge seems to not handle time.Time for some reason, thus fixup
	fixupTournamentInfo(&t.Info, tdata)
	if err := storage.Store(t); err != nil {
		return nil, errors.New(err.Error() + " - Could not write tournament to storage")
	}
	return t, nil
}

func AllTournaments() (Tournaments, error) {
	return storage.LoadAll()
}

func DeleteByUUID(uuid uuid.UUID) bool {
	err := storage.Delete(uuid)
	if err != nil {
		return false
	}
	return true
}

func TournamentByUUID(uuid uuid.UUID) (*Tournament, error) {
	return storage.Load(uuid)
}

func TournamentsBySeason(season int) (Tournaments, error) {
	return storage.LoadBySeason(season)
}

func TournamentsByPeriod(from time.Time, to time.Time) (Tournaments, error) {
	tournaments, err := storage.LoadAll()
	if err != nil {
		return nil, errors.New(err.Error() + " - Could not load tournaments from storage")
	}

	var inRange Tournaments
	for _, t := range tournaments {
		if t.Info.Scheduled.After(from) && t.Info.Scheduled.Before(to) {
			inRange = append(inRange, t)
		}
	}
	return inRange, nil
}

func (t *Tournament) UpdateInfo(tdata Info) error {
	if err := mergo.MergeWithOverwrite(&t.Info, tdata); err != nil {
		return errors.New(err.Error() + " - Could not update tournament info")
	}
	// Merge seems to not handle time.Time for some reason, thus fixup
	fixupTournamentInfo(&t.Info, tdata)
	err := storage.Store(t)
	if err != nil {
		return errors.New(err.Error() + " - Could not store updated tournament info")
	}
	return nil
}

func (t *Tournament) SetPlayed(isPlayed bool) error {
	t.Played = isPlayed
	err := storage.Store(t)
	if err != nil {
		return errors.New(err.Error() + " - Could not store updated tournament state")
	}
	return nil
}

func (t *Tournament) SetResult(result Result) error {
	t.Played = true
	t.Result = result
	err := storage.Store(t)
	if err != nil {
		return errors.New(err.Error() + " - Could not store tournament result")
	}
	return nil
}
