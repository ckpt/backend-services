package main

import (
	"encoding/json"
	"errors"
	"github.com/ckpt/backend-services/news"
	"github.com/m4rw3r/uuid"
	"github.com/zenazn/goji/web"
	"net/http"
)

func createNewNewsItem(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	nNewsItem := new(news.NewsItem)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(nNewsItem); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}
	nNewsItem, err := news.NewNewsItem(*nNewsItem, c.Env["authPlayer"].(uuid.UUID))
	if err != nil {
		return &appError{err, "Failed to create new NewsItem", 500}
	}
	w.Header().Set("Location", "/news/"+nNewsItem.UUID.String())
	w.WriteHeader(201)
	return nil
}

func listAllNews(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	list, err := news.AllNewsItems()
	if err != nil {
		return &appError{err, "Cant load NewsItems", 500}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(list)
	return nil
}

func getNewsItem(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	newsItem, err := news.NewsItemByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find the NewsItem", 404}
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(newsItem)
	return nil
}

func updateNewsItem(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	newsUUID, err := uuid.FromString(c.URLParams["uuid"])
	newsItem, err := news.NewsItemByUUID(newsUUID)
	if err != nil {
		return &appError{err, "Cant find NewsItem", 404}
	}
	if !c.Env["authIsAdmin"].(bool) && c.Env["authPlayer"].(uuid.UUID) != newsItem.Author {
		return &appError{errors.New("Unauthorized"), "Must be author or admin to update news item", 403}
	}
	tempNewsItem := new(news.NewsItem)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tempNewsItem); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if err := newsItem.UpdateNewsItem(*tempNewsItem); err != nil {
		return &appError{err, "Failed to update news item", 500}
	}
	w.WriteHeader(204)
	return nil
}

func addNewsComment(c web.C, w http.ResponseWriter, r *http.Request) *appError {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	newsUUID, err := uuid.FromString(c.URLParams["uuid"])
	newsItem, err := news.NewsItemByUUID(newsUUID)
	if err != nil {
		return &appError{err, "Cant find NewsItem", 404}
	}
	tempInfo := new(news.Comment)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(tempInfo); err != nil {
		return &appError{err, "Invalid JSON", 400}
	}

	if !c.Env["authIsAdmin"].(bool) || tempInfo.Player.IsZero() {
		tempInfo.Player = c.Env["authPlayer"].(uuid.UUID)
	}
	if err := newsItem.AddComment(tempInfo.Player, tempInfo.Content); err != nil {
		return &appError{err, "Failed to add news comment", 500}
	}
	w.WriteHeader(204)
	return nil
}
