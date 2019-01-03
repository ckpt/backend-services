package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

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

	if !players.AuthUser(loginReq.Username, loginReq.Password) {
		return &appError{errors.New("Forbidden"), "Invalid username/password", 403}
	}

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

func main() {
	//
	// Event queue hadling
	//
	err := players.StartEventProcessor()
	if err != nil {
		fmt.Printf("%+v", err.Error())
		println("Could not initialize event queue. Exiting")
		os.Exit(1)
	}

	//
	// HTTP Serving
	//
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"GET", "PUT", "PATCH", "POST", "OPTIONS", "DELETE"},
	})
	goji.Use(c.Handler)
	goji.Use(middleware.TokenHandler)

	goji.Post("/login", appHandler(login))

	goji.Get("/players", appHandler(listAllPlayers))
	goji.Post("/players", appHandler(createNewPlayer))
	goji.Get("/players/quotes", appHandler(getAllPlayerQuotes))
	goji.Get("/players/:uuid", appHandler(getPlayer))
	goji.Put("/players/:uuid", appHandler(updatePlayer))
	goji.Post("/players/:uuid/quotes", appHandler(addPlayerQuote))
	goji.Get("/players/:uuid/profile", appHandler(getPlayerProfile))
	goji.Put("/players/:uuid/profile", appHandler(updatePlayerProfile))
	goji.Get("/players/:uuid/user", appHandler(getUserForPlayer))
	goji.Put("/players/:uuid/user", appHandler(setUserForPlayer))
	goji.Put("/players/:uuid/user/password", appHandler(setUserPassword))
	goji.Put("/players/:uuid/user/settings", appHandler(setUserSettings))
	goji.Put("/players/:uuid/user/admin", appHandler(setUserAdmin))
	goji.Put("/players/:uuid/gossip", appHandler(setPlayerGossip))
	goji.Patch("/players/:uuid/gossip", appHandler(setPlayerGossip))
	goji.Delete("/players/:uuid/gossip", appHandler(resetPlayerGossip))
	goji.Get("/players/:uuid/debts", appHandler(showPlayerDebt))
	goji.Delete("/players/:uuid/debts", appHandler(resetPlayerDebts))
	goji.Get("/players/:uuid/credits", appHandler(showPlayerCredits))
	goji.Post("/players/:uuid/debts", appHandler(addPlayerDebt))
	goji.Delete("/players/:uuid/debts/:debtuuid", appHandler(settlePlayerDebt))
	goji.Put("/players/:uuid/votes", appHandler(setPlayerVotes))
	goji.Patch("/players/:uuid/votes", appHandler(setPlayerVotes))
	goji.Post("/players/notification_test", appHandler(testPlayerNotify))

	goji.Post("/users", appHandler(createNewUser))

	goji.Get("/locations", appHandler(listAllLocations))
	goji.Post("/locations", appHandler(createNewLocation))
	goji.Get("/locations/:uuid", appHandler(getLocation))
	goji.Put("/locations/:uuid", appHandler(updateLocationProfile))
	goji.Patch("/locations/:uuid", appHandler(updateLocationProfile))
	goji.Post("/locations/:uuid/pictures", appHandler(addLocationPicture))

	goji.Get("/tournaments", appHandler(listAllTournaments))
	goji.Post("/tournaments", appHandler(createNewTournament))
	goji.Get("/tournaments/:uuid", appHandler(getTournament))
	goji.Put("/tournaments/:uuid", appHandler(updateTournamentInfo))
	goji.Patch("/tournaments/:uuid", appHandler(updateTournamentInfo))
	goji.Put("/tournaments/:uuid/played", appHandler(setTournamentPlayed))
	goji.Get("/tournaments/:uuid/result", appHandler(getTournamentResult))
	goji.Put("/tournaments/:uuid/result", appHandler(setTournamentResult))
	goji.Put("/tournaments/:uuid/bountyhunters", appHandler(setTournamentBountyHunters))
	goji.Post("/tournaments/:uuid/noshows", appHandler(addTournamentNoShow))
	goji.Delete("/tournaments/:uuid/noshows/:playeruuid", appHandler(removeTournamentNoShow))

	goji.Get("/seasons", appHandler(listAllSeasons))
	goji.Get("/seasons/stats", appHandler(getTotalStats))
	goji.Get("/seasons/standings", appHandler(getTotalStandings))
	goji.Get("/seasons/titles", appHandler(getTotalTitles))
	goji.Get("/seasons/:year/tournaments", appHandler(listTournamentsBySeason))
	goji.Get("/seasons/:year/standings", appHandler(getSeasonStandings))
	goji.Get("/seasons/:year/titles", appHandler(getSeasonTitles))
	goji.Get("/seasons/:year/stats", appHandler(getSeasonStats))

	goji.Get("/caterings", appHandler(listAllCaterings))
	goji.Post("/caterings", appHandler(createNewCatering))
	goji.Get("/caterings/:uuid", appHandler(getCatering))
	goji.Put("/caterings/:uuid", appHandler(updateCateringInfo))
	goji.Patch("/caterings/:uuid", appHandler(updateCateringInfo))
	goji.Post("/caterings/:uuid/votes", appHandler(addCateringVote))
	goji.Put("/caterings/:uuid/votes/:playeruuid", appHandler(updateCateringVote))

	goji.Get("/news", appHandler(listAllNews))
	goji.Get("/news/:uuid", appHandler(getNewsItem))
	goji.Patch("/news/:uuid", appHandler(updateNewsItem))
	goji.Post("/news", appHandler(createNewNewsItem))
	goji.Post("/news/:uuid/comments", appHandler(addNewsComment))
	// TODO: Comment updates/deletion

	goji.Serve()
}
