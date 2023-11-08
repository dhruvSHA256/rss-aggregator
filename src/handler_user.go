package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rss-aggregator/internal/database"
	"time"

	"github.com/google/uuid"
)

func (apiCfg apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name  string `name`
		Email string `email`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprint("error parsing JSON:", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Email:     params.Email,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprint("Couldn't create user", err))
		return
	}
	respondWithJson(w, 200, dataBaseUsertoUser(user))
}

func (apiCfg apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJson(w, 200, dataBaseUsertoUser(user))
}
