package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Vikuuu/synlabs-assignment/internal/auth"
	"github.com/Vikuuu/synlabs-assignment/internal/database"
)

func handlerLandingPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Working"))
}

type signupPayload struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	UserType        string `json:"user_type"`
	ProfileHeadline string `json:"profile_headline"`
	Address         string `json:"address"`
}

type signupResponse struct {
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	UserType database.UserType `json:"user_type"`
}

func (cfg *apiConfig) handlerSignUp(w http.ResponseWriter, r *http.Request) {
	payload := signupPayload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Fatalf("error decoding JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userType := database.UserTypeApplicant
	if payload.UserType == "admin" {
		userType = database.UserTypeAdmin
	}

	hashPassword, err := auth.HashPassword(payload.Password)

	// Add the in sql
	dat, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		Name:            payload.Name,
		Email:           payload.Email,
		Address:         payload.Address,
		UserType:        userType,
		PasswordHash:    hashPassword,
		ProfileHeadline: payload.ProfileHeadline,
	})
	if err != nil {
		log.Fatalf("error saving to db: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := signupResponse{
		Name:     dat.Name,
		Email:    dat.Email,
		UserType: dat.UserType,
	}

	resp, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error encoding JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
