package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errResponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, msg string, code int) {
	er := errResponse{
		Error: msg,
	}
	errResp, err := json.Marshal(er)
	if err != nil {
		log.Fatalf("error encoding JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(errResp)
}
