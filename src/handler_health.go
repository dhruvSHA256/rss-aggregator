package main

import "net/http"

func handleCheckHealth(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, struct{}{})
}
