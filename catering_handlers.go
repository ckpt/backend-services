package main

import (
	"encoding/json"
	"github.com/ckpt/backend-services/caterings"
	"github.com/m4rw3r/uuid"
	"github.com/zenazn/goji/web"
	"net/http"
)

func createNewCatering(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	nCatering := new(caterings.Catering)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(nCatering); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}
	nCatering, err := caterings.NewCatering(nCatering.Tournament, nCatering.Info)
	if err != nil {
		return &appError{err, "Failed to create new catering", 500}
	}
	w.Header().Set("Location", "/caterings/"+nCatering.UUID.String())
	w.WriteHeader(201)
	return nil
}

func listAllCaterings(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	list, err := caterings.AllCaterings()
	if err != nil {
		return &appError{err, "Cant load caterings", 500}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(list)
	return nil
}

func getCatering(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	catering, err := caterings.CateringByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find catering", 404}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(catering)
	return nil
}

func updateCateringInfo(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	catering, err := caterings.CateringByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find catering", 404}
	}
	tempInfo := new(caterings.Info)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tempInfo); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := catering.UpdateInfo(*tempInfo); err != nil {
		return &appError{err, "Failed to update catering info", 500}
	}
	w.WriteHeader(204)
	return nil
}

func addCateringVote(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	catering, err := caterings.CateringByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find catering", 404}
	}
	tempInfo := new(caterings.Vote)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tempInfo); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := catering.AddVote(tempInfo.Player, tempInfo.Score); err != nil {
		return &appError{err, "Failed to add catering vote", 500}
	}
	w.WriteHeader(204)
	return nil
}

func updateCateringVote(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	cateringuuid, err := uuid.FromString(c.URLParams["uuid"])
	playeruuid, err := uuid.FromString(c.URLParams["playeruuid"])
	catering, err := caterings.CateringByUUID(cateringuuid)
	if err != nil {
		return &appError{err, "Cant find catering", 404}
	}

	tempInfo := new(caterings.Vote)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tempInfo); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := catering.RemoveVote(playeruuid); err != nil {
		return &appError{err, "Failed to remove old catering vote when updating", 500}
	}

	if err := catering.AddVote(playeruuid, tempInfo.Score); err != nil {
		return &appError{err, "Failed to add updated catering vote", 500}
	}
	w.WriteHeader(204)
	return nil
}
