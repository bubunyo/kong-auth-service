package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type D map[string]interface{}

// DataResponse wraps all response objects in a json `data` object
type DataResponse struct {
	Data interface{} `json:"data"`
}

type UserData struct {
	Data struct {
		EmailAddress string `json:"email_address"`
		Password     string `json:"password"`
	} `json:"data"`
}

type RequestHandler func(w http.ResponseWriter, req *http.Request)

func resHandler(w http.ResponseWriter, r interface{}) error {
	return json.NewEncoder(w).Encode(DataResponse{Data: r})
}

func AccountRoutes(r *mux.Router, db *sql.DB) {
	k := NewKong()
	r.HandleFunc("/", CreateAccount(db, k)).Methods("POST")
	r.HandleFunc("/create", Authenticate(db, k)).Methods("POST")
}

func (u UserData) Validate() error {
	return nil
}

func CreateAccount(db *sql.DB, kong *Kong) RequestHandler {
	return func(w http.ResponseWriter, req *http.Request) {
		// Validate Data
		// Create User
		// Register details in kong and get Credentials
		// Generate JWT

		_ = resHandler(w, D{"msg": "Account Creation"})
	}
}

func Authenticate(db *sql.DB, kong *Kong) RequestHandler {
	return func(w http.ResponseWriter, req *http.Request) {

		// Validate Data
		// Find User or return 404
		// Generate JWT

		_ = resHandler(w, D{"msg": "Authenticate"})
	}
}
