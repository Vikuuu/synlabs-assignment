package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Vikuuu/synlabs-assignment/internal/database"
)

type jobListResponse struct {
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	PostedOn          time.Time `json:"posted_on"`
	TotalApplications int32     `json:"total_application"`
	CompanyName       string    `json:"company_name"`
	PostedBy          int32     `json:"posted_by"`
}

func (cfg *apiConfig) handlerViewJobs(w http.ResponseWriter, r *http.Request) {
	data, err := cfg.db.GetJobsApplicant(context.Background())
	if err != nil {
		log.Printf("error getting data: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := []jobListResponse{}
	for _, val := range data {
		i := jobListResponse{
			Title:             val.Title,
			Description:       val.Description,
			PostedOn:          val.PostedOn,
			TotalApplications: val.TotalApplications.Int32,
			CompanyName:       val.CompanyName,
			PostedBy:          val.PostedBy,
		}
		res = append(res, i)
	}
	resp, err := json.Marshal(res)
	if err != nil {
		log.Printf("error marshaling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (cfg *apiConfig) handlerApplyJob(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.Atoi(r.URL.Query().Get("job_id"))
	if err != nil {
		respondWithError(w, "Failed to get JobID", http.StatusInternalServerError)
		return
	}
	userID := r.Context().Value("userID").(int)

	err = cfg.db.ApplyJob(context.Background(), database.ApplyJobParams{
		ApplicantID: sql.NullInt32{Int32: int32(userID), Valid: true},
		JobID:       sql.NullInt32{Int32: int32(jobID), Valid: true},
	})

	err = cfg.db.UpdateTotalApplications(context.Background(), int32(jobID))
	if err != nil {
		respondWithError(
			w,
			"Failed to increase the count of total applicants",
			http.StatusInternalServerError,
		)
		return
	}

	type applyResponse struct {
		Success bool `json:"success"`
	}
	res := applyResponse{Success: true}

	resp, err := json.Marshal(res)
	if err != nil {
		log.Printf("error marshaling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
