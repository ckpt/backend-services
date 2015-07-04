package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/m4rw3r/uuid"
	"net/http"

	"github.com/rs/cors"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"github.com/ckpt/backend-services/players"
	"github.com/ckpt/backend-services/locations"
	"github.com/ckpt/backend-services/middleware"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(web.C, http.ResponseWriter, *http.Request) *appError

var currentuser string

func (fn appHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	if e := fn(c, w, r); e != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(e.Code)
		encoder := json.NewEncoder(w)
		encoder.Encode(map[string]string{"error": e.Error.Error() +
			" (" + e.Message + ")"})
	}
}

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

func login(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	type LoginRequest struct {
		Username string
		Password string
	}

	loginReq := new(LoginRequest)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(loginReq); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}
	// Hard code for now
	fmt.Printf("%+v\n", loginReq)
	if loginReq.Username == "mortenk" &&
		loginReq.Password == "testing123" {
		authUser, err := players.UserByName(loginReq.Username)
		if err != nil {
			return &appError{err, "Failed to fetch user data", 500}
		}
		if authUser.Locked {
			return &appError{errors.New("Locked"), "User locked", 403}
		}
		encoder := json.NewEncoder(w)
		encoder.Encode(authUser)
		return nil
	}

	// Else, forbidden
	w.WriteHeader(403)
	return nil
}


func createNewLocation(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	nLocation := new(locations.Location)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(nLocation); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}
	nLocation, err := locations.NewLocation(nLocation.Host, nLocation.Profile)
	if err != nil {
		return &appError{err, "Failed to create new location", 500}
	}
	w.Header().Set("Location", "/locations/"+nLocation.UUID.String())
	w.WriteHeader(201)
	return nil
}

func listAllLocations(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	loclist, err := locations.AllLocations()
	if err != nil {
		return &appError{err, "Cant load locations", 500}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(loclist)
	return nil
}

func getLocation(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	location, err := locations.LocationByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find location", 404}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(location)
	return nil
}

func updateLocationProfile(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	location, err := locations.LocationByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find location", 404}
	}
	tempProfile := new(locations.Profile)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tempProfile); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := location.UpdateProfile(*tempProfile); err != nil {
		return &appError{err, "Failed to update location profile", 500}
	}
	w.WriteHeader(204)
	return nil
}

func addLocationPicture(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	location, err := locations.LocationByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find location", 404}
	}

	type Message struct {
		Picture []byte
	}
	pic := new(Message)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(pic); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}
	if err != nil {
		return &appError{err, "Picture is not base64 encoded", 400}
	}

	if err := location.AddPicture(pic.Picture); err != nil {
		return &appError{err, "Failed to add location picture", 500}
	}
	w.WriteHeader(201)
	return nil
}

func main() {
	//fmt.Printf("%+v\n", getMembers())
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	})
	goji.Use(c.Handler)
	goji.Use(middleware.TokenHandler)

	goji.Post("/login", appHandler(login))

	goji.Get("/players", appHandler(listAllPlayers))
	goji.Post("/players", appHandler(createNewPlayer))
	goji.Get("/players/quotes", appHandler(getAllPlayerQuotes))
	goji.Get("/players/:uuid", appHandler(getPlayer))
	goji.Put("/players/:uuid", appHandler(updatePlayer))
	goji.Get("/players/:uuid/profile", appHandler(getPlayerProfile))
	goji.Put("/players/:uuid/profile", appHandler(updatePlayerProfile))
	goji.Get("/players/:uuid/user", appHandler(getUserForPlayer))
	goji.Put("/players/:uuid/user", appHandler(setUserForPlayer))

	goji.Post("/users", appHandler(createNewUser))

	goji.Get("/locations", appHandler(listAllLocations))
	goji.Post("/locations", appHandler(createNewLocation))
	goji.Get("/locations/:uuid", appHandler(getLocation))
	goji.Put("/locations/:uuid", appHandler(updateLocationProfile))
	goji.Patch("/locations/:uuid", appHandler(updateLocationProfile))
	goji.Post("/locations/:uuid/pictures", appHandler(addLocationPicture))

	goji.Serve()
}
