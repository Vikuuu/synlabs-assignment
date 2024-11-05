package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/Vikuuu/synlabs-assignment/internal/database"
)

type apiConfig struct {
	db     *database.Queries
	secret string
}

func (cfg *apiConfig) WithAuthAdmin(handler http.HandlerFunc) http.Handler {
	return cfg.middlewareIsAdmin(handler)
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("connection cannot be made to db: %s", err)
	}
	defer db.Close()

	config := apiConfig{
		db:     database.New(db),
		secret: os.Getenv("SECRET"),
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("GET /", handlerLandingPage)
	mux.HandleFunc("POST /signup", config.handlerSignUp)
	mux.HandleFunc("POST /login", config.handlerLogIn)
	mux.HandleFunc("POST /uploadResume", config.handlerUploadResume)
	mux.Handle("POST /admin/job", config.WithAuthAdmin(config.handlerAddJob))
	mux.Handle("GET /admin/job/{job_id}", config.WithAuthAdmin(config.handlerJob))

	log.Printf("Serving on Port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
