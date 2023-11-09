package handler

import (
	"net/http"
	"gateway/helper"
)

func HandleCheckHealth(w http.ResponseWriter, r *http.Request) {
	helper.RespondWithJson(w, 200, struct{}{})
}
