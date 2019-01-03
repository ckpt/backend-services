package caterings

import (
	"errors"

	"github.com/imdario/mergo"
	"github.com/m4rw3r/uuid"
)

// We use dummy in memory storage for now
var storage CateringStorage = NewRedisCateringStorage()

type Catering struct {
	UUID       uuid.UUID `json:"uuid"`
	Info       Info      `json:"info"`
	Tournament uuid.UUID `json:"tournament"`
	Votes      []Vote    `json:"votes"`
}

type Info struct {
	Caterer uuid.UUID `json:"caterer"`
	Meal    string    `json:"meal"`
}

type Vote struct {
	Player uuid.UUID `json:"player"`
	Score  int       `json:"score"`
}

// A storage interface for Caterings
type CateringStorage interface {
	Store(*Catering) error
	Delete(uuid.UUID) error
	Load(uuid.UUID) (*Catering, error)
	LoadAll() ([]*Catering, error)
	LoadByTournament(uuid.UUID) (*Catering, error)
	//LoadByPlayer(uuid.UUID) ([]*Catering, error)
}

//
// Catering related functions and methods
//

// Create a Catering
func NewCatering(tournament uuid.UUID, ci Info) (*Catering, error) {
	c := new(Catering)
	c.UUID, _ = uuid.V4()
	c.Tournament = tournament
	if err := mergo.MergeWithOverwrite(&c.Info, ci); err != nil {
		return nil, errors.New(err.Error() + " - Could not set initial catering info")
	}
	if err := storage.Store(c); err != nil {
		return nil, errors.New(err.Error() + " - Could not write catering to storage")
	}
	return c, nil
}

func AllCaterings() ([]*Catering, error) {
	return storage.LoadAll()
}

func DeleteByUUID(uuid uuid.UUID) bool {
	err := storage.Delete(uuid)
	if err != nil {
		return false
	}
	return true
}

func CateringByUUID(uuid uuid.UUID) (*Catering, error) {
	return storage.Load(uuid)
}

func (c *Catering) UpdateInfo(ci Info) error {
	if err := mergo.MergeWithOverwrite(&c.Info, ci); err != nil {
		return errors.New(err.Error() + " - Could not update catering info")
	}
	err := storage.Store(c)
	if err != nil {
		return errors.New(err.Error() + " - Could not store updated catering info")
	}
	return nil
}

func (c *Catering) AddVote(player uuid.UUID, score int) error {
	vote := Vote{Player: player, Score: score}
	c.Votes = append(c.Votes, vote)
	err := storage.Store(c)
	if err != nil {
		return errors.New(err.Error() + " - Could not store updated catering info with added vote")
	}
	return nil
}

func (c *Catering) RemoveVote(player uuid.UUID) error {
	for i, v := range c.Votes {
		if v.Player == player {
			c.Votes = append(c.Votes[:i], c.Votes[i+1:]...)
		}
	}
	err := storage.Store(c)
	if err != nil {
		return errors.New(err.Error() + " - Could not store updated catering info with removed vote")
	}
	return nil
}
