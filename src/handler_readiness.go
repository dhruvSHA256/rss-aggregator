package main

import "net/http"

func checkHealth(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, 200, struct{}{})
}
