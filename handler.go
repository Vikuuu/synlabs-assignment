package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)

	if userType == "applicant" {
		appID, err := cfg.db.CreateApplicantProfile(context.Background(), dat.ID)
		if err != nil {
			log.Fatalf("error creating applicant profile: %s", err)
			return
		}

		err = cfg.db.AddProfileIDInUser(
			context.Background(),
			sql.NullInt32{Int32: appID, Valid: true},
		)
		if err != nil {
			log.Fatalf("error adding profile_id: %s", err)
			return
		}
	}
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string            `json:"access_token"`
	UserType    database.UserType `json:"user_type"`
}

func (cfg *apiConfig) handlerLogIn(w http.ResponseWriter, r *http.Request) {
	payload := loginPayload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)
	if err != nil {
		log.Fatalf("error decoding JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get user data from the table
	user, err := cfg.db.GetUser(context.Background(), payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := "Invalid credentials"
			respondWithError(w, errMsg, http.StatusUnauthorized)
			return
		} else {
			log.Fatalf("error getting user: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	err = auth.CheckPassword(payload.Password, user.PasswordHash)
	if err != nil {
		errMsg := "Invalid credentials"
		respondWithError(w, errMsg, http.StatusUnauthorized)
		return
	}

	// create jwt and return the response
	expiresIn := time.Hour
	tokenSecret := cfg.secret

	jwtToken, err := auth.MakeJWT(user.ID, tokenSecret, expiresIn)
	if err != nil {
		log.Fatalf("error creating JWT: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(loginResponse{
		AccessToken: jwtToken,
		UserType:    user.UserType,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

type uploadResumePayload struct {
	FileAddress string `json:"file_address"`
}

type apiPayload struct {
	Name      string   `json:"name"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	Skills    []string `json:"skills"`
	Education []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"education"`
	Experience []struct {
		Dates []string `json:"dates"`
		Name  string   `json:"name"`
		URL   string   `json:"url"`
	} `json:"experience"`
}

type uploadResumeResponse struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Skills    string `json:"skills"`
	Education string `json:"education"`
}

func (cfg *apiConfig) handlerUploadResume(w http.ResponseWriter, r *http.Request) {
	// TODO: Authenticate API: (GetBearerToken)
	jwt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Fatalf("error getting jwtToken: %s", err)
		respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateJWT(jwt, cfg.secret)
	if err != nil {
		log.Printf("error validating token: %s", err)
		respondWithError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// TODO: Authorization: Applicant only
	userType, err := cfg.db.GetUserFromID(context.Background(), int32(userID))
	if err != nil {
		log.Fatalf("error getting user: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userType == "admin" {
		respondWithError(w, "You are not allowed to access this endpoint", http.StatusUnauthorized)
		return
	}

	// TODO: Decode the payload
	payload := uploadResumePayload{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&payload)
	if err != nil {
		log.Fatalf("error decoding JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("FileAddress: %s", payload.FileAddress)
	// TODO: Check for suffix to be .pdf or .docx
	if !strings.HasSuffix(payload.FileAddress, "pdf") &&
		!strings.HasSuffix(payload.FileAddress, "docx") {
		respondWithError(w, "Only pdf and docx file supported", http.StatusBadRequest)
		return
	}

	// TODO: Open the file
	file, err := os.Open(payload.FileAddress)
	if err != nil {
		log.Fatalf("error opening file: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// TODO: Make call to the 3rd party API
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.apilayer.com/resume_parser/upload", file)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("apiKey", "0bWeisRWoLj3UdXt3MXMSMWptYFIpQfS")

	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("error making request to API: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("Response: %v", res)
	// TODO: Decode the data from 3rd party API
	apiPl := apiPayload{}
	decoder = json.NewDecoder(res.Body)
	err = decoder.Decode(&apiPl)
	if err != nil {
		log.Fatalf("error marshaling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Concatenate the skills into comma separated.
	skills := strings.Join(apiPl.Skills, ",")
	e := []string{}
	for _, edu := range apiPl.Education {
		e = append(e, edu.Name)
	}
	educations := strings.Join(e, ",")
	// TODO: Save all the details into Applicant profile
	dat, err := cfg.db.UpdateProfile(context.Background(), database.UpdateProfileParams{
		Name:      sql.NullString{String: apiPl.Name, Valid: true},
		Email:     sql.NullString{String: apiPl.Email, Valid: true},
		Phone:     sql.NullString{String: apiPl.Phone, Valid: true},
		Skills:    sql.NullString{String: skills, Valid: true},
		Education: sql.NullString{String: educations, Valid: true},
		Applicant: int32(userID),
	})
	if err != nil {
		log.Fatalf("error updating profile: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: If all goes well return 200
	resp, err := json.Marshal(uploadResumeResponse{
		Name:      dat.Name.String,
		Email:     dat.Email.String,
		Phone:     dat.Phone.String,
		Skills:    dat.Skills.String,
		Education: dat.Education.String,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
