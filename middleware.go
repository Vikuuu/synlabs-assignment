package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Vikuuu/synlabs-assignment/internal/auth"
)

func (cfg *apiConfig) middlewareIsAdmin(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Fatalf("error getting jwtToken: %s", err)
			respondWithError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateJWT(jwt, cfg.secret)
		if err != nil {
			log.Fatalf("error validating token: %s", err)
			respondWithError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userType, err := cfg.db.GetUserFromID(context.Background(), int32(userID))
		if err != nil {
			log.Fatalf("error getting user: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if userType == "applicant" {
			respondWithError(
				w,
				"You are not authorized to access this endpoint",
				http.StatusUnauthorized,
			)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (cfg *apiConfig) middlewareIsApplicant(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Fatalf("error getting jwtToken: %s", err)
			respondWithError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateJWT(jwt, cfg.secret)
		if err != nil {
			log.Fatalf("error validating token: %s", err)
			respondWithError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userType, err := cfg.db.GetUserFromID(context.Background(), int32(userID))
		if err != nil {
			log.Fatalf("error getting user: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if userType == "admin" {
			respondWithError(
				w,
				"You are not authorized to access this endpoint",
				http.StatusUnauthorized,
			)
			return
		}

		log.Println(userID)
		ctx := context.WithValue(r.Context(), "userID", userID)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}
