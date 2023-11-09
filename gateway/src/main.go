package main

import (
	"database/sql"
	"gateway/handler"
	"gateway/internal/database"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant connect to DB")
	}

	db := database.New(conn)
	apiCfg := handler.ApiConfig{
		DB: db,
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
	v1Router.Post("/user", apiCfg.HandleCreateUser)
	v1Router.Get("/user", apiCfg.MiddlewareAuth(apiCfg.HandleGetUser))
	v1Router.Post("/auth", apiCfg.HandleAuthUser)
	v1Router.Post("/feed", apiCfg.MiddlewareAuth(apiCfg.HandleCreateFeed))
	v1Router.Get("/feeds", apiCfg.HandleGetFeeds)
	// v1Router.Post("/feed_follow", apiCfg.MiddlewareAuth(apiCfg.HandleFollowFeed))
	// v1Router.Get("/feed_follow", apiCfg.MiddlewareAuth(apiCfg.HandleGetFollowFeed))
	// v1Router.Delete("/feed_follow/{feedFollowID}", apiCfg.MiddlewareAuth(apiCfg.HandleDeleteFollowFeed))
	// v1Router.Get("/posts", apiCfg.MiddlewareAuth(apiCfg.HandleGetPosts))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Running on: %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
