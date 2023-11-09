package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"gateway/helper"
	"gateway/internal/database"
	"gateway/models"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg ApiConfig) HandleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `name`
		Url  string `url`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helper.RespondWithError(w, 400, fmt.Sprint("error parsing JSON:", err))
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
		helper.RespondWithError(w, 400, fmt.Sprint("Couldn't create user", err))
		return
	}
	helper.RespondWithJson(w, 200, models.DataBaseFeedtoFeed(feed))
}

func (apiCfg ApiConfig) HandleGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())
	if err != nil {
		helper.RespondWithError(w, 404, fmt.Sprintf("Coundlt get feeds: %v\n", err))
		return
	}
	helper.RespondWithJson(w, 200, models.DataBaseFeedstoFeeds(feeds))
}

func (apiCfg ApiConfig) HandleFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helper.RespondWithError(w, 400, fmt.Sprint("error parsing JSON:", err))
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
		helper.RespondWithError(w, 400, fmt.Sprint("Couldn't follow feed", err))
		return
	}
	helper.RespondWithJson(w, 200, models.DatabaseFeedFollowtoFeedFollow(feedFollow))
}

func (apiCfg ApiConfig) HandleGetFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feeds, err := apiCfg.DB.GetFeedFollow(r.Context(), user.ID)
	if err != nil {
		helper.RespondWithError(w, 404, fmt.Sprintf("Couldnt get feeds: %v\n", err))
		return
	}
	helper.RespondWithJson(w, 200, models.DataBaseFeedstoFeeds(feeds))
}

func (apiCfg ApiConfig) HandleDeleteFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	feedId, err := uuid.Parse(chi.URLParam(r, "feedFollowId"))
	if err != nil {
		helper.RespondWithError(w, 400, fmt.Sprint("error parsing id:", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedId,
	})
	if err != nil {
		helper.RespondWithError(w, 400, fmt.Sprint("couldnt delete feed follow", err))
		return
	}
	helper.RespondWithJson(w, 200, struct{}{})
}

func (apiCfg ApiConfig) HandleGetPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPosts(r.Context(), user.ID)
	if err != nil {
		helper.RespondWithError(w, 404, fmt.Sprintf("Couldnt get posts: %v\n", err))
		return
	}
	helper.RespondWithJson(w, 200, posts)
}
