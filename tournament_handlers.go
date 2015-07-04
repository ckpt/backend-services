package main

import (
	"encoding/json"
	"github.com/ckpt/backend-services/tournaments"
	"github.com/m4rw3r/uuid"
	"github.com/zenazn/goji/web"
	"net/http"
)

func createNewTournament(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	info := new(tournaments.Info)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(info); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}
	nTournament, err := tournaments.NewTournament(*info)
	if err != nil {
		return &appError{err, "Failed to create new tournament", 500}
	}
	w.Header().Set("Location", "/tournaments/"+nTournament.UUID.String())
	w.WriteHeader(201)
	return nil
}

func listAllTournaments(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	tournamentList, err := tournaments.AllTournaments()
	if err != nil {
		return &appError{err, "Cant load tournaments", 500}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(tournamentList)
	return nil
}

func getTournament(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	tournament, err := tournaments.TournamentByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find tournament", 404}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(tournament)
	return nil
}

func updateTournamentInfo(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	tournament, err := tournaments.TournamentByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find tournament", 404}
	}
	tempInfo := new(tournaments.Info)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tempInfo); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := tournament.UpdateInfo(*tempInfo); err != nil {
		return &appError{err, "Failed to update tournament info", 500}
	}
	w.WriteHeader(204)
	return nil
}
