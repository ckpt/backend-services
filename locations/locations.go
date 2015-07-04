package locations

import (
	"errors"
	"github.com/imdario/mergo"
	"github.com/m4rw3r/uuid"
)

// We use dummy in memory storage for now
var storage LocationStorage = NewRedisLocationStorage()

type Coord struct {
	Lat  float64 `json:lat`
	Long float64 `json:long`
}

type Location struct {
	UUID     uuid.UUID `json:"uuid"`
	Host     uuid.UUID `json:"host"`
	Profile  Profile   `json:"profile"`
	Pictures [][]byte  `json:"pictures"`
	Active   bool      `json:"active"`
}

type Profile struct {
	URL         string   `json:"url"`
	Coordinates Coord    `json:"coordinates"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Facilities  []string `json:"facilities"`
}

// A storage interface for Locations
type LocationStorage interface {
	Store(*Location) error
	Delete(uuid.UUID) error
	Load(uuid.UUID) (*Location, error)
	LoadAll() ([]*Location, error)
	LoadByPlayer(uuid.UUID) (*Location, error)
}

//
// Location related functions and methods
//

// Create a Location
func NewLocation(host uuid.UUID, lp Profile) (*Location, error) {
	l := new(Location)
	l.UUID, _ = uuid.V4()
	l.Active = true
	l.Host = host
	if err := mergo.MergeWithOverwrite(&l.Profile, lp); err != nil {
		return nil, errors.New(err.Error() + " - Could not set initial location profile")
	}
	if err := storage.Store(l); err != nil {
		return nil, errors.New(err.Error() + " - Could not write location to storage")
	}
	return l, nil
}

func AllLocations() ([]*Location, error) {
	return storage.LoadAll()
}

func DeleteByUUID(uuid uuid.UUID) bool {
	err := storage.Delete(uuid)
	if err != nil {
		return false
	}
	return true
}

func LocationByUUID(uuid uuid.UUID) (*Location, error) {
	return storage.Load(uuid)
}

func (l *Location) AddPicture(picture []byte) error {
	l.Pictures = append(l.Pictures, picture)
	err := storage.Store(l)
	if err != nil {
		return errors.New(err.Error() + " - Could not add picture to location")
	}
	return nil
}
func (l *Location) RemovePicture(picIndex int) error {
	l.Pictures = append(l.Pictures[:picIndex], l.Pictures[picIndex+1:]...)
	err := storage.Store(l)
	if err != nil {
		return errors.New(err.Error() + " - Could not delete picture at index " + string(picIndex))
	}
	return nil
}
func (l *Location) UpdateProfile(lp Profile) error {
	if err := mergo.MergeWithOverwrite(&l.Profile, lp); err != nil {
		return errors.New(err.Error() + " - Could not update location profile")
	}
	err := storage.Store(l)
	if err != nil {
		return errors.New(err.Error() + " - Could not store updated location profile")
	}
	return nil
}
