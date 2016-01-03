package players

import (
	"errors"
	"github.com/ckpt/backend-services/utils"
	"github.com/imdario/mergo"
	"github.com/m4rw3r/uuid"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

// We use dummy in memory storage for now
var storage PlayerStorage = NewRedisPlayerStorage()

// Init a message queue
var eventqueue utils.AMQPQueue = utils.NewRMQ(os.Getenv("CKPT_AMQP_URL"), "ckpt.events")

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
	Votes      Votes             `json:"votes"`
	// Debts where this player is the debitor
	Debts []Debt `json:"debts"`
}

// The basic profile of the player
type Profile struct {
	Name        string    `json:"name"`
	Picture     []byte    `json:"picture"`
	Birthday    time.Time `json:"birthday"`
	Email       string    `json:"email"`
	Description string    `json:"description"`
}

// A debt to another player (hte creditor)
type Debt struct {
	UUID        uuid.UUID `json:"uuid"`
	Debitor     uuid.UUID `json:"debitor"`
	Creditor    uuid.UUID `json:"creditor"`
	Description string    `json:"description"`
	Amount      int       `json:"amount"`
	Created     time.Time `json:"created"`
	Settled     time.Time `json:"settled"`
}

// A complaint from another player
// for harassment or similar
type Complaint struct {
	From    *Player `json:"from"`
	Content string  `json:"content"`
}

type Votes struct {
	Winner uuid.UUID `json:"winner"`
	Loser  uuid.UUID `json:"loser"`
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

func (p *Player) SetUserSettings(settings UserSettings) error {
	p.User.Settings = settings
	if err := storage.Store(p); err != nil {
		return errors.New(err.Error() + " - Could not change player user settings")
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
func (p *Player) AddDebt(d Debt) error {
	newDebt := new(Debt)
	if err := mergo.MergeWithOverwrite(newDebt, d); err != nil {
		return errors.New(err.Error() + " - Could not set Debt data")
	}
	newDebt.UUID, _ = uuid.V4()
	if d.Created.IsZero() {
		newDebt.Created = time.Now()
	} else {
		newDebt.Created = d.Created
	}
	if !d.Settled.IsZero() {
		newDebt.Settled = d.Settled
	}
	newDebt.Debitor = p.UUID
	p.Debts = append(p.Debts, *newDebt)
	err := storage.Store(p)
	if err != nil {
		return errors.New("Could not add debt")
	}
	eventqueue.Publish(utils.CKPTEvent{
		Type:         utils.PLAYER_EVENT,
		RestrictedTo: []uuid.UUID{p.UUID},
		Subject:      "Gjeld registrert",
		Message:      "Det er registrert et nytt gjeldskrav mot deg på ckpt.no!"})
	return nil
}
func (p *Player) SettleDebt(debtuuid uuid.UUID) error {
	for i, debt := range p.Debts {
		if debt.UUID == debtuuid {
			p.Debts[i].Settled = time.Now()
		}
	}
	err := storage.Store(p)
	if err != nil {
		return errors.New("Could not settle debt")
	}
	eventqueue.Publish(utils.CKPTEvent{
		Type:         utils.PLAYER_EVENT,
		RestrictedTo: []uuid.UUID{p.UUID},
		Subject:      "Gjeld tilbakebetalt",
		Message:      "Et av dine utestående krav er markert som innfridd på ckpt.no!"})
	return nil
}
func (p *Player) ResetDebt() error {
	p.Debts = []Debt{}
	err := storage.Store(p)
	if err != nil {
		return errors.New("Could not reset debt")
	}
	return nil
}

func (p *Player) DebtByUUID(uuid uuid.UUID) (*Debt, error) {
	for _, debt := range p.Debts {
		if debt.UUID == uuid {
			return &debt, nil
		}
	}
	return nil, errors.New("Debt not found")
}
func (creditor *Player) Credits() ([]Debt, error) {
	var credits []Debt
	all, _ := AllPlayers()
	for _, p := range all {
		if p.UUID == creditor.UUID {
			continue
		}
		for _, d := range p.Debts {
			if d.Creditor == creditor.UUID {
				credits = append(credits, d)
			}
		}
	}
	return credits, nil
}
func (p *Player) SetVotes(v Votes) error {
	if err := mergo.MergeWithOverwrite(&p.Votes, v); err != nil {
		return errors.New(err.Error() + " - Could not set Votes data")
	}
	err := storage.Store(p)
	if err != nil {
		return errors.New("Could not set votes")
	}
	return nil
}
