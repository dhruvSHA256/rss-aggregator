package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rss-aggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

func (apiCfg apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `name`
		Url  string `url`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint("error parsing JSON:", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprint("Couldn't create user", err))
		return
	}
	respondWithJson(w, 200, dataBaseFeedtoFeed(feed))
}

func (apiCfg apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Coundlt get feeds: %v\n", err))
		return
	}
	respondWithJson(w, 200, dataBaseFeedstoFeeds(feeds))
}

func (apiCfg apiConfig) handleFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint("error parsing JSON:", err))
		return
	}
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprint("Couldn't follow feed", err))
		return
	}
	respondWithJson(w, 200, databaseFeedFollowtoFeedFollow(feedFollow))
}

func (apiCfg apiConfig) handleGetFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfg.DB.GetFeedFollow(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("Couldnt get feeds: %v\n", err))
		return
	}
	respondWithJson(w, 200, dataBaseFeedstoFeeds(feeds))
}
