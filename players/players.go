package players

import (
	"errors"
	"github.com/m4rw3r/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// We use dummy in memory storage for now
var storage PlayerStorage = NewRedisPlayerStorage()

// Constants
type VoteType int

const (
	VOTE_FAV VoteType = iota
	VOTE_LOSER
)

// A Player is a player in CKPT, current or former.o// It also contains a User.
type Player struct {
	UUID    uuid.UUID `json:"uuid"`
	Profile Profile   `json:"profile"`
	Nick    string    `json:"nick"`
	User    User      `json:"user"`
	Active  bool      `json:"active"`
	// Quotes are just strings set by other players
	Quotes []string `json:"quotes"`
	// Gossip is one string from each
	// of the other players, indexed by uuid as string
	Gossip     map[string]string `json:"gossip"`
	Complaints []Complaint       `json:"complaints"`
	Votes      []Vote            `json:"votes"`
}

// The basic profile of the player
type Profile struct {
	Name        string    `json:"name"`
	Picture     []byte    `json:"picture"`
	Birthday    time.Time `json:"birthday"`
	Email       string    `json:"email"`
	Description string    `json:"description"`
}

// A debt from one player to another
type Debt struct {
	UUID     uuid.UUID `json:"uuid"`
	Debitor  *Player   `json:"debitor"`
	Creditor *Player   `json:"creditor"`
	Due      time.Time `json:"due"`
	Amount   int       `json:"amount"`
}

// A complaint from another player
// for harassment or similar
type Complaint struct {
	From    *Player `json:"from"`
	Content string  `json:"content"`
}

type Vote struct {
	From *Player  `json:"from"`
	Type VoteType `json:"type"`
}

// A storage interface for Players
type PlayerStorage interface {
	Store(*Player) error
	Delete(uuid.UUID) error
	Load(uuid.UUID) (*Player, error)
	LoadAll() ([]*Player, error)
	LoadUser(username string) (*User, error)
}

//
// Player related functions and methods
//

// Create a player
func NewPlayer(nick string, profile Profile) (*Player, error) {
	p := new(Player)
	newUUID, err := uuid.V4()
	if err != nil {
		// FIXME: Handle error
	}
	p.UUID = newUUID
	p.Nick = nick
	p.Profile = profile
	err = storage.Store(p)
	if err != nil {
		return nil, errors.New(err.Error() + " - Could not write player to storage")
	}
	return p, nil
}

func AllPlayers() ([]*Player, error) {
	return storage.LoadAll()
}

func DeleteByUUID(uuid uuid.UUID) bool {
	err := storage.Delete(uuid)
	if err != nil {
		return false
	}
	return true
}

func PlayerByUUID(uuid uuid.UUID) (*Player, error) {
	return storage.Load(uuid)
}

func PlayerByUserToken(token string) (*Player, error) {
	players, err := storage.LoadAll()
	if err != nil {
		return nil, errors.New(err.Error() + " - Could not load player by token")
	}
	for _, p := range players {
		if p.User.Apikey == token {
			return p, nil
		}
	}
	return nil, errors.New("Could not find player with given token")
}

func (p *Player) SetUser(user User) error {
	p.User = user
	err := storage.Store(p)
	if err != nil {
		return errors.New("Could not change player user")
	}
	return nil
}

func (p *Player) SetUserPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    if err != nil {
		return errors.New(err.Error() + " - Could not change player user password")
	}
	p.User.password = string(hashedPassword)
	if err := storage.Store(p); err != nil {
		return errors.New(err.Error() + " - Could not change player user password")
	}
	return nil
}

func (p *Player) SetProfile(profile Profile) error {
	p.Profile = profile
	err := storage.Store(p)
	if err != nil {
		return errors.New("Could not change profile")
	}
	return nil
}
func (p *Player) SetNick(nick string) error {
	p.Nick = nick
	err := storage.Store(p)
	if err != nil {
		return errors.New("Could not change nick")
	}
	return nil
}
func (p *Player) SetActive(active bool) error {
	p.Active = active
	err := storage.Store(p)
	if err != nil {
		return errors.New("Could not change active status")
	}
	return nil
}
