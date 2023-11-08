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

// func (apiCfg apiConfig) handleGetFeed(w http.ResponseWriter, r *http.Request) {
// 	apiKey, err := auth.GetAPIKey(r.Header)
// 	if err != nil {
// 		respondWithError(w, 400, fmt.Sprintf("Couldnt get api key: %v", err))
// 		return
// 	}

// 	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
// 	if err != nil {
// 		respondWithError(w, 404, fmt.Sprintf("User not found: %v", err))
// 		return
// 	}
// 	respondWithJson(w, 200, dataBaseUsertoUser(user))
// }
