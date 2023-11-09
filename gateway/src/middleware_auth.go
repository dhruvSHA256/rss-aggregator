package main

import (
	"fmt"
	"net/http"
	"rss-aggregator/internal/auth"
	"rss-aggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := auth.DecodeJWT(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Invalid token: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByID(r.Context(), userId)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("User not found: %v", err))
			return
		}
		handler(w, r, user)
	}
}
