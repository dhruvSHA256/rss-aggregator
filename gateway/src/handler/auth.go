package handler

import (
	"fmt"
	"gateway/helper"
	"gateway/internal/auth"
	"gateway/internal/database"
	"net/http"
	"strings"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *ApiConfig) MiddlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get("Authorization")
		if val == "" {
			helper.RespondWithError(w, 403, fmt.Sprint("auth header missing"))
			return
		}
		vals := strings.Split(val, " ")
		if len(vals) != 2 {
			helper.RespondWithError(w, 403, fmt.Sprint("invalid auth header"))
		}

		token := vals[1]
		userId, err := auth.DecodeJWT(token)
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
