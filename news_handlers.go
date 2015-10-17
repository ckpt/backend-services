package main

import (
	"encoding/json"
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
	uuid, err := uuid.FromString(c.URLParams["uuid"])
	newsItem, err := news.NewsItemByUUID(uuid)
	if err != nil {
		return &appError{err, "Cant find NewsItem", 404}
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

// Vote -> Comment
// func addnewsVote(c web.C, w http.ResponseWriter, r *http.Request) *appError {
// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	uuid, err := uuid.FromString(c.URLParams["uuid"])
// 	news, err := newss.newsByUUID(uuid)
// 	if err != nil {
// 		return &appError{err, "Cant find news", 404}
// 	}
// 	tempInfo := new(newss.Vote)
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(tempInfo); err != nil {
// 		return &appError{err, "Invalid JSON", 400}
// 	}

// 	if err := news.AddVote(tempInfo.Player, tempInfo.Score); err != nil {
// 		return &appError{err, "Failed to add news vote", 500}
// 	}
// 	w.WriteHeader(204)
// 	return nil
// }

// func updatenewsVote(c web.C, w http.ResponseWriter, r *http.Request) *appError {
// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	newsuuid, err := uuid.FromString(c.URLParams["uuid"])
// 	playeruuid, err := uuid.FromString(c.URLParams["playeruuid"])
// 	news, err := newss.newsByUUID(newsuuid)
// 	if err != nil {
// 		return &appError{err, "Cant find news", 404}
// 	}

// 	tempInfo := new(newss.Vote)
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(tempInfo); err != nil {
// 		return &appError{err, "Invalid JSON", 400}
// 	}

// 	if err := news.RemoveVote(playeruuid); err != nil {
// 		return &appError{err, "Failed to remove old news vote when updating", 500}
// 	}

// 	if err := news.AddVote(playeruuid, tempInfo.Score); err != nil {
// 		return &appError{err, "Failed to add updated news vote", 500}
// 	}
// 	w.WriteHeader(204)
// 	return nil
// }
