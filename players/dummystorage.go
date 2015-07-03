package players

import (
	"errors"
	"github.com/m4rw3r/uuid"
	"time"
)

type DummyPlayerStorage struct {
	players []*Player
	users   []*User
	debts   []*Debt
}

func (dps *DummyPlayerStorage) init() {
	dummyUUIDs := createUUIDs(2)
	dps.users = []*User{
		&User{
			Username: "mortenk",
			password: "admin123",
			Apikey:   "secretsupersecret",
			Admin:    true,
			Locked:   false,
		},
	}
	dps.players = []*Player{
		&Player{
			UUID: dummyUUIDs[0],
			Profile: Profile{
				Birthday: time.Date(1979, time.April, 14, 0, 0, 0, 0, time.Local),
				Name:     "Morten Knutsen",
				Email:    "knumor@gmail.com",
			},
			Nick:   "Panzer",
			Quotes: []string{"Blinde høner kan også finne korn!"},
			User:   *dps.users[0],
			Active: true,
		},
		&Player{
			UUID: dummyUUIDs[1],
			Profile: Profile{
				Birthday: time.Date(1979, time.October, 20, 0, 0, 0, 0, time.Local),
				Name:     "Bjørn Helge Kopperud",
				Email:    "bjorn@kjekkegutter.no",
			},
			Nick:   "Bjøro",
			Quotes: []string{"Horespill!"},
			Active: true,
		},
	}
}

func (dps *DummyPlayerStorage) Store(p *Player) error {
	for i := range dps.players {
		if dps.players[i].UUID == p.UUID {
			dps.players[i] = p
			return nil
		}
	}
	dps.players = append(dps.players, p)
	return nil
}

func (dps *DummyPlayerStorage) Load(uuid uuid.UUID) (*Player, error) {
	for _, player := range dps.players {
		if player.UUID == uuid {
			return player, nil
		}
	}
	return nil, errors.New("Not found")
}

func (dps *DummyPlayerStorage) Delete(uuid uuid.UUID) error {
	// FIXME: Not implemented yet
	return errors.New("Not implemented yet")
}

func (dps *DummyPlayerStorage) LoadAll() ([]*Player, error) {
	return dps.players, nil
}

func (dps *DummyPlayerStorage) LoadUser(username string) (*User, error) {
	for _, user := range dps.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("Not found")
}

func NewDummyPlayerStorage() *DummyPlayerStorage {
	dps := new(DummyPlayerStorage)
	dps.init()
	return dps
}

func createUUIDs(number int) []uuid.UUID {
	var uuids []uuid.UUID
	for number > 0 {
		uuid, _ := uuid.V4()
		uuids = append(uuids, uuid)
		number--
	}
	return uuids
}
