package handler

import (
	"fmt"
	"net/http"
	"gateway/helper"
	"gateway/internal/auth"
	"gateway/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *ApiConfig) MiddlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := auth.DecodeJWT(r.Header)
		if err != nil {
			helper.RespondWithError(w, 403, fmt.Sprintf("Invalid token: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUserByID(r.Context(), userId)
		if err != nil {
			helper.RespondWithError(w, 404, fmt.Sprintf("User not found: %v", err))
			return
		}
		handler(w, r, user)
	}
}
