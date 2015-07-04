package main

import (
	"encoding/json"
	"github.com/ckpt/backend-services/locations"
	"github.com/m4rw3r/uuid"
	"github.com/zenazn/goji/web"
	"net/http"
)

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
