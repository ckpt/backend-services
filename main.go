package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"github.com/ckpt/backend-services/middleware"
	"github.com/ckpt/backend-services/players"
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
