package players

// The user associated with a player
type User struct {
	Username string `json:"username"`
	password string
	Apikey   string       `json:"apikey"`
	Admin    bool         `json:"admin"`
	Locked   bool         `json:"locked"`
	Settings UserSettings `json:"settings"`
}

// The user preferences/settings of the user
type UserSettings struct {
	Notifications map[string]bool `json:"notifications"`
}

func UserByName(username string) (*User, error) {
	return storage.LoadUser(username)
}
