package main

import (
	"github.com/m4rw3r/uuid"
	"time"
)

// -------------------------------------
// Data models
// -------------------------------------

// Member
//
// A Member is a player in CKPT, current or former.
// It also functions as a user.
type Member struct {
	UUID    uuid.UUID     `json:"uuid"`
	Profile MemberProfile `json:"profile"`
	Nick    string        `json:"nick"`
	User    MemberUser    `json:"-"`
	Active  bool          `json:"active"`
}

// The basic profile of the member
type MemberProfile struct {
	Name        string    `json:"name"`
	Picture     []byte    `json:"picture"`
	Birthday    time.Time `json:"birthday"`
	Email       string    `json:"email"`
	Description string    `json:"description"`
}

// The user associated with the member
type MemberUser struct {
	Username string `json:"username"`
	password string
	locked   bool
	Settings MemberUserSettings `json:"settings"`
}

// The user preferences/settings of the user
type MemberUserSettings struct {
	Notifications map[string]bool `json:"notifications"`
}

// 	* (api_token)
// * other accounts?
// * debts
// * quotes
// * votes
// * gossip
// * complaints
// * notification settings
