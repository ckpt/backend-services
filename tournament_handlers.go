package main

import (
	"encoding/json"
	"errors"
	"github.com/ckpt/backend-services/tournaments"
	"github.com/m4rw3r/uuid"
	"github.com/zenazn/goji/web"
	"net/http"
	"sort"
	"strconv"
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

func setTournamentPlayed(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	tournament, err := tournaments.TournamentByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find tournament", 404}
	}

	tempInfo := make(map[string]bool)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&tempInfo); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := tournament.SetPlayed(tempInfo["played"]); err != nil {
		return &appError{err, "Failed to update tournament played status", 500}
	}
	w.WriteHeader(204)
	return nil
}

func setTournamentResult(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	tID, err := uuid.FromString(c.URLParams["uuid"])
	tournament, err := tournaments.TournamentByUUID(tID)
	if err != nil {
		return &appError{err, "Cant find tournament", 404}
	}

	type Result struct {
		Result []uuid.UUID
	}

	resultData := new(Result)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(resultData); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := tournament.SetResult(resultData.Result); err != nil {
		return &appError{err, "Failed to update tournament result", 500}
	}
	w.WriteHeader(204)
	return nil
}

func addTournamentNoShow(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	tID, err := uuid.FromString(c.URLParams["uuid"])
	tournament, err := tournaments.TournamentByUUID(tID)
	if err != nil {
		return &appError{err, "Cant find tournament", 404}
	}

	absenteeData := new(tournaments.Absentee)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(absenteeData); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if !c.Env["authIsAdmin"].(bool) || absenteeData.Player.IsZero() {
		absenteeData.Player = c.Env["authPlayer"].(uuid.UUID)
	}

	if err := tournament.AddNoShow(absenteeData.Player, absenteeData.Reason); err != nil {
		return &appError{err, "Failed to set absentee for tournament", 500}
	}
	w.WriteHeader(204)
	return nil
}

func removeTournamentNoShow(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	tID, err := uuid.FromString(c.URLParams["uuid"])
	tournament, err := tournaments.TournamentByUUID(tID)
	if err != nil {
		return &appError{err, "Cant find tournament", 404}
	}

	pID, err := uuid.FromString(c.URLParams["playeruuid"])

	if !c.Env["authIsAdmin"].(bool) && pID != c.Env["authPlayer"].(uuid.UUID) {
		return &appError{errors.New("Unauthorized"), "Must be given player or admin to remove absentee", 403}
	}

	if err := tournament.RemoveNoShow(pID); err != nil {
		return &appError{err, "Failed to remove absentee for tournament", 500}
	}
	w.WriteHeader(204)
	return nil
}

func getTournamentResult(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	tID, err := uuid.FromString(c.URLParams["uuid"])
	tournament, err := tournaments.TournamentByUUID(tID)
	if err != nil {
		return &appError{err, "Cant find tournament", 404}
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(tournament.Result)
	return nil
}

func listAllSeasons(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	allTournaments, err := tournaments.AllTournaments()
	if err != nil {
		return &appError{err, "Cant load tournaments", 404}
	}

	seasonList := allTournaments.Seasons()
	sort.Ints(seasonList)
	encoder := json.NewEncoder(w)
	encoder.Encode(map[string][]int{"seasons": seasonList})
	return nil
}

func listTournamentsBySeason(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	season, err := strconv.Atoi(c.URLParams["year"])
	tList, err := tournaments.TournamentsBySeason(season)
	if err != nil {
		return &appError{err, "Cant find tournaments", 404}
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(map[string][]*tournaments.Tournament{"tournaments": tList})
	return nil
}

func getSeasonStandings(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	season, _ := strconv.Atoi(c.URLParams["year"])
	sortedStandings := tournaments.SeasonStandings(season)

	encoder := json.NewEncoder(w)
	encoder.Encode(sortedStandings)
	return nil
}

func getSeasonStats(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	season, _ := strconv.Atoi(c.URLParams["year"])

	seasonStats := tournaments.SeasonStats([]int{season})

	encoder := json.NewEncoder(w)
	encoder.Encode(seasonStats)
	return nil
}

func getSeasonTitles(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	season, _ := strconv.Atoi(c.URLParams["year"])

	seasonTitles := tournaments.Titles([]int{season})

	encoder := json.NewEncoder(w)
	encoder.Encode(seasonTitles)
	return nil
}

func getTotalStandings(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	tList, err := tournaments.AllTournaments()
	if err != nil {
		return &appError{err, "Cant find tournaments", 404}
	}

	seasons := tList.Seasons()

	totalStandings := tournaments.TotalStandings(seasons)

	encoder := json.NewEncoder(w)
	encoder.Encode(totalStandings)
	return nil
}

func getTotalStats(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	tList, err := tournaments.AllTournaments()
	if err != nil {
		return &appError{err, "Cant find tournaments", 404}
	}

	seasons := tList.Seasons()

	fullStats := tournaments.SeasonStats(seasons)

	encoder := json.NewEncoder(w)
	encoder.Encode(fullStats)
	return nil
}

func getTotalTitles(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	tList, err := tournaments.AllTournaments()
	if err != nil {
		return &appError{err, "Cant find tournaments", 404}
	}

	seasons := tList.Seasons()

	allTitles := tournaments.Titles(seasons)

	encoder := json.NewEncoder(w)
	encoder.Encode(allTitles)
	return nil
}
