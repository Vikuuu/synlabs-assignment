package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Vikuuu/synlabs-assignment/internal/database"
)

type addJobPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CompanyName string `json:"company_name"`
}

type addJobResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CompanyName string `json:"company_name"`
}

func (cfg *apiConfig) handlerAddJob(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		respondWithError(w, "Failed to retrieve user ID", http.StatusInternalServerError)
		return
	}

	payload := addJobPayload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Printf("error decoding JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: create job openings
	data, err := cfg.db.CreateJob(context.Background(), database.CreateJobParams{
		Title:       payload.Title,
		Description: payload.Description,
		PostedOn:    time.Now(),
		CompanyName: payload.CompanyName,
		PostedBy:    int32(userID),
	})
	if err != nil {
		log.Printf("error creating job: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(addJobResponse{
		Title:       data.Title,
		Description: data.Description,
		CompanyName: data.CompanyName,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

type jobResponse struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PostedOn    time.Time `json:"posted_on"`
	CompanyName string    `json:"company_name"`
	PostedBy    int32     `json:"posted_by"`
}

func (cfg *apiConfig) handlerJob(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.Atoi(r.PathValue("job_id"))
	if err != nil {
		log.Printf("error in type changing: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := cfg.db.GetJob(context.Background(), int32(jobID))
	if err != nil {
		respondWithError(w, "error getting job", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(jobResponse{
		Title:       data.Title,
		Description: data.Description,
		PostedOn:    data.PostedOn,
		CompanyName: data.CompanyName,
		PostedBy:    data.PostedBy,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
