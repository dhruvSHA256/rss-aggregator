package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"rss-aggregator/internal/database"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func respondWithError(w http.ResponseWriter, code int, payload string) {
	if code > 499 {
		log.Println("Responding with 5xx error:", payload)
	}
	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errResponse{Error: payload})

}

func main() {
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not set in env")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set in env")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant connect to DB")
	}

	apiCfg := apiConfig{
		DB: database.New(db),
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handleCheckHealth)
	v1Router.Post("/user", apiCfg.handleCreateUser)
	v1Router.Get("/user", apiCfg.middlewareAuth(apiCfg.handleGetUser))
	v1Router.Post("/feed", apiCfg.middlewareAuth(apiCfg.handleCreateFeed))
	v1Router.Get("/feeds", apiCfg.handleGetFeeds)
	v1Router.Post("/feed_follow", apiCfg.middlewareAuth(apiCfg.handleFollowFeed))
	v1Router.Get("/feed_follow", apiCfg.middlewareAuth(apiCfg.handleGetFollowFeed))
	v1Router.Delete("/feed_follow", apiCfg.middlewareAuth(apiCfg.handleDeleteFollowFeed))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Running on: %v", portString)
	err1 := srv.ListenAndServe()
	if err1 != nil {
		log.Fatal(err)
	}

}
