package players

import "errors"
import "github.com/m4rw3r/uuid"
import "golang.org/x/crypto/bcrypt"

//import "fmt"

// The user associated with a player
type User struct {
	Username string `json:"username"`
	password string
	Apikey   string       `json:"apikey"`
	Admin    bool         `json:"admin"`
	Locked   bool         `json:"locked"`
	Settings UserSettings `json:"settings"`
}

// Create a user
func NewUser(player uuid.UUID, userdata *User) (*User, error) {
	p, err := storage.Load(player)
	if err != nil {
		return nil, err
	}
	p.User = *userdata
	err = storage.Store(p)
	if err != nil {
		return nil,
			errors.New(err.Error() + " - Could not write user to storage")
	}
	return &p.User, nil
}

// The user preferences/settings of the user
type UserSettings struct {
	Notifications map[string]bool `json:"notifications"`
}

func UserByName(username string) (*User, error) {
	return storage.LoadUser(username)
}

func AuthUser(username string, password string) bool {
	user, err := storage.LoadUser(username)
	if err != nil {
		return false
	}
	// No error from comparison means the hashes match
	if err := bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password)); err == nil {
		return true
	}
	return false
}
