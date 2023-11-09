package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"gateway/helper"
	"gateway/internal/auth"
	"gateway/internal/database"
	"gateway/models"
	"time"

	"github.com/google/uuid"
)

func (apiCfg ApiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name  string `name`
		Email string `email`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		helper.RespondWithError(w, 400, fmt.Sprint("error parsing JSON:", err))
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
		helper.RespondWithError(w, 400, fmt.Sprint("Couldn't create user", err))
		return
	}
	jwt_token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		log.Fatalln("Unable to generate jwt", err)
	} else {
		log.Println(jwt_token)
	}
	helper.RespondWithJson(w, 200, models.DataBaseUsertoUser(user))
}

func (apiCfg ApiConfig) HandleGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	helper.RespondWithJson(w, 200, models.DataBaseUsertoUser(user))
}
