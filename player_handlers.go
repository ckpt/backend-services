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

	if !c.Env["authIsAdmin"].(bool) && c.Env["authPlayer"].(uuid.UUID) != pUUID {
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

func setUserSettings(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])

	if !c.Env["authIsAdmin"].(bool) && c.Env["authPlayer"].(uuid.UUID) != pUUID {
		return &appError{errors.New("Unauthorized"), "Must be correct user or admin to chenge settings", 403}
	}

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}

	sUpdate := new(players.UserSettings)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(sUpdate); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := player.SetUserSettings(*sUpdate); err != nil {
		return &appError{err, "Failed to change settings for user", 500}
	}
	w.WriteHeader(204)
	return nil
}

func setUserAdmin(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])

	if !c.Env["authIsAdmin"].(bool) {
		return &appError{errors.New("Unauthorized"), "Must be admin to set admin status", 403}
	}

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}

	var adminState bool
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&adminState); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := player.SetUserAdmin(adminState); err != nil {
		return &appError{err, "Failed to change settings for user", 500}
	}
	w.WriteHeader(204)
	return nil
}


func showPlayerDebt(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])

	if !c.Env["authIsAdmin"].(bool) && c.Env["authPlayer"].(uuid.UUID) != pUUID {
		return &appError{errors.New("Unauthorized"), "Must be player or admin to show debt", 403}
	}

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(player.Debts)
	return nil
}

func showPlayerCredits(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])

	if !c.Env["authIsAdmin"].(bool) && c.Env["authPlayer"].(uuid.UUID) != pUUID {
		return &appError{errors.New("Unauthorized"), "Must be player or admin to show credits", 403}
	}

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}

	credits, _ := player.Credits()

	encoder := json.NewEncoder(w)
	encoder.Encode(credits)
	return nil
}

func addPlayerDebt(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}

	nDebt := new(players.Debt)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(nDebt); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if !c.Env["authIsAdmin"].(bool) &&
		(c.Env["authPlayer"].(uuid.UUID) == pUUID ||
			c.Env["authPlayer"].(uuid.UUID) != nDebt.Creditor) {
		return &appError{errors.New("Unauthorized"), "Must be creditor or admin to add debt", 403}
	}

	err = player.AddDebt(*nDebt)
	if err != nil {
		return &appError{err, "Failed to add debt", 500}
	}
	w.Header().Set("Location", "/players/"+pUUID.String()+"/debts")
	w.WriteHeader(201)
	return nil
}

func settlePlayerDebt(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])
	dUUID, err := uuid.FromString(c.URLParams["debtuuid"])

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}

	debt, err := player.DebtByUUID(dUUID)
	if err != nil {
		return &appError{err, "Cant find debt for player", 404}
	}

	if !c.Env["authIsAdmin"].(bool) && c.Env["authPlayer"].(uuid.UUID) != debt.Creditor {
		return &appError{errors.New("Unauthorized"), "Must be creditor or admin to settle debt", 403}
	}

	err = player.SettleDebt(dUUID)
	if err != nil {
		return &appError{err, "Failed to settle debt", 500}
	}
	w.Header().Set("Location", "/players/"+pUUID.String()+"/debts")
	w.WriteHeader(204)
	return nil
}

func resetPlayerDebts(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}

	if !c.Env["authIsAdmin"].(bool) {
		return &appError{errors.New("Unauthorized"), "Must be admin to reset debts", 403}
	}

	err = player.ResetDebt()
	if err != nil {
		return &appError{err, "Failed to reset debts", 500}
	}
	w.WriteHeader(204)
	return nil
}

func setPlayerVotes(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}

	nVotes := new(players.Votes)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(nVotes); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if !c.Env["authIsAdmin"].(bool) && c.Env["authPlayer"].(uuid.UUID) != pUUID {
		return &appError{errors.New("Unauthorized"), "Must be player or admin to set votes", 403}
	}

	err = player.SetVotes(*nVotes)
	if err != nil {
		return &appError{err, "Failed to set votes", 500}
	}
	w.WriteHeader(204)
	return nil
}

func addPlayerQuote(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}
	
	var q string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&q); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}


	if !c.Env["authIsAdmin"].(bool) && c.Env["authPlayer"].(uuid.UUID) == pUUID {
		return &appError{errors.New("Unauthorized"), "Must be other player or admin to add quote", 403}
	}

	err = player.AddQuote(q)
	if err != nil {
		return &appError{err, "Failed to add quote", 500}
	}
	w.Header().Set("Location", "/players/"+pUUID.String())
	w.WriteHeader(201)
	return nil
}

func setPlayerGossip(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}

	nGossip := make(map[string]string)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&nGossip); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if !c.Env["authIsAdmin"].(bool) && c.Env["authPlayer"].(uuid.UUID) != pUUID {
		return &appError{errors.New("Unauthorized"), "Must be player or admin to set gossip", 403}
	}

	err = player.SetGossip(nGossip)
	if err != nil {
		return &appError{err, "Failed to set gossip", 500}
	}
	w.WriteHeader(204)
	return nil
}

func resetPlayerGossip(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	pUUID, err := uuid.FromString(c.URLParams["uuid"])

	player, err := players.PlayerByUUID(pUUID)
	if err != nil {
		return &appError{err, "Cant find player", 404}
	}

	if !c.Env["authIsAdmin"].(bool) && c.Env["authPlayer"].(uuid.UUID) != pUUID {
		return &appError{errors.New("Unauthorized"), "Must be admin or player to reset gossip", 403}
	}

	err = player.ResetGossip()
	if err != nil {
		return &appError{err, "Failed to reset gossip", 500}
	}
	w.WriteHeader(204)
	return nil
}



func testPlayerNotify(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if !c.Env["authIsAdmin"].(bool) {
		return &appError{errors.New("Unauthorized"), "Must be admin to test notifications", 403}
	}

	players.NotifyUser("Panzer", "knumor@gmail.com", "Test", "Åpenbar test")
	w.WriteHeader(204)
	return nil
}
