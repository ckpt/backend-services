package main

import (
	"encoding/json"
	"errors"
	"github.com/ckpt/backend-services/players"
	"github.com/m4rw3r/uuid"
	"github.com/zenazn/goji/web"
	"net/http"
)

func listAllPlayers(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	playerlist, err := players.AllPlayers()
	if err != nil {
		return &appError{err, "Cant load players", 500}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(playerlist)
	return nil
}

func getAllPlayerQuotes(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	playerlist, err := players.AllPlayers()
	if err != nil {
		return &appError{err, "Cant load players", 500}
	}
	quotes := make(map[string][]string)
	for _, player := range playerlist {
		pq := player.Quotes
		if len(pq) > 0 {
			quotes[player.UUID.String()] = append(quotes[player.UUID.String()], pq...)
		}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(quotes)
	return nil
}

func getPlayer(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	player, err := players.PlayerByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(player)
	return nil
}

func getPlayerProfile(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	player, err := players.PlayerByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(player.Profile)
	return nil
}

func updatePlayer(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	player, err := players.PlayerByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}
	tempPlayer := new(players.Player)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tempPlayer); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := player.SetActive(tempPlayer.Active); err != nil {
		return &appError{err, "Failed to set active status", 500}
	}
	if err := player.SetNick(tempPlayer.Nick); err != nil {
		return &appError{err, "Failed to set nick", 500}
	}
	if err := player.SetProfile(tempPlayer.Profile); err != nil {
		return &appError{err, "Failed to set player profile", 500}
	}
	w.WriteHeader(204)
	return nil
}

func updatePlayerProfile(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	player, err := players.PlayerByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}
	tempProfile := new(players.Profile)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tempProfile); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := player.SetProfile(*tempProfile); err != nil {
		return &appError{err, "Failed to set player profile", 500}
	}
	w.WriteHeader(204)
	return nil
}

func createNewPlayer(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	nPlayer := new(players.Player)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(nPlayer); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}
	nPlayer, err := players.NewPlayer(nPlayer.Nick, nPlayer.Profile)
	if err != nil {
		return &appError{err, "Failed to create new player", 500}
	}
	w.Header().Set("Location", "/players/"+nPlayer.UUID.String())
	w.WriteHeader(201)
	return nil
}

func createNewUser(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	type newUser struct {
		Player uuid.UUID `json:"player"`
		players.User
	}
	nUser := new(newUser)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(nUser); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}
	_, err := players.NewUser(nUser.Player, &nUser.User)
	if err != nil {
		return &appError{err, "Failed to create new user", 500}
	}
	w.Header().Set("Location", "/players/"+nUser.Player.String()+"/user")
	w.WriteHeader(201)
	return nil
}

func getUserForPlayer(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	player, err := players.PlayerByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(player.User)
	return nil
}

func setUserForPlayer(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if !c.Env["authIsAdmin"].(bool) {
		return &appError{errors.New("Unauthorized"), "Admins only", 403}
	}
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	player, err := players.PlayerByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}
	tempUser := new(players.User)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tempUser); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}
	user, err := players.UserByName(tempUser.Username)
	if err != nil {
		return &appError{err, "Cant find user", 400}
	}

	if err := player.SetUser(*user); err != nil {
		return &appError{err, "Failed to set user for player", 500}
	}
	w.WriteHeader(204)
	return nil
}

func setUserPassword(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])

	if !c.Env["authIsAdmin"].(bool) || c.Env["authPlayer"].(uuid.UUID) != pUUID {
		return &appError{errors.New("Unauthorized"), "Must be correct user or admin to set password", 403}
	}

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}

	type PWUpdate struct {
		Password string
	}

	pwupdate := new(PWUpdate)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(pwupdate); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := player.SetUserPassword(pwupdate.Password); err != nil {
		return &appError{err, "Failed to set password for player", 500}
	}
	w.WriteHeader(204)
	return nil
}
