package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type D map[string]interface{}

// ResponseData wraps all response objects in a json `data` object
type ResponseData struct {
	Data interface{} `json:"data"`
}

type Data struct {
	ID           string `db:"id"`
	EmailAddress string `json:"email_address",validate:"required,email",db:"email"`
	Password     string `json:"password"validate:"required",db:"password"`
}

type ReqBody struct {
	Data Data `json:"data",validate:"required"`
}
type JSONErrs []error

func (je JSONErrs) MarshalJSON() ([]byte, error) {
	res := make([]interface{}, len(je))
	for i, e := range je {
		if _, ok := e.(json.Marshaler); ok {
			res[i] = e // e knows how to marshal itself
		} else {
			res[i] = e.Error() // Fallback to the error string
		}
	}
	return json.Marshal(res)
}

func (d Data) hashPassword() string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.MinCost)
	return string(hash)
}
func (d Data) PasswordMatches(hashedPwd string) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(d.Password))
	if err != nil {
		return false
	}
	return true
}

func (u ReqBody) Validate(v *validator.Validate) error {
	return v.Struct(u)
}

type RequestHandler func(w http.ResponseWriter, req *http.Request)

func SuccessResponse(s int, w http.ResponseWriter, r interface{}) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(ResponseData{Data: r})
}

func ErrorResponse(s int, w http.ResponseWriter, r error) error {
	w.WriteHeader(s)
	return json.NewEncoder(w).Encode(ResponseData{Data: JSONErrs([]error{r})})
}

func Credentials(id, token string) map[string]interface{} {
	return D{
		"account_id": id,
		"credentials": D{
			"jwt": token,
		},
	}
}

func AccountRoutes(r *mux.Router, db *sql.DB) {
	k := NewKong()
	v := validator.New()
	r.HandleFunc("/register", CreateAccount(db, k, v)).Methods("POST")
	r.HandleFunc("/login", Authenticate(db, k, v)).Methods("POST")
}

func CreateAccount(db *sql.DB, kong *Kong, validator *validator.Validate) RequestHandler {
	return func(w http.ResponseWriter, req *http.Request) {
		// Get Data form body
		reqBody := ReqBody{}
		err := json.NewDecoder(req.Body).Decode(&reqBody)
		//Could not decode json
		if err != nil {
			_ = ErrorResponse(http.StatusBadRequest, w, err)
			return
		}

		// Validate Data
		// Check validation rules in the struct tags of ReqBody
		err = reqBody.Validate(validator)
		if err != nil {
			_ = ErrorResponse(http.StatusBadRequest, w, err)
			return
		}

		// Create User
		err = db.QueryRow(
			"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
			reqBody.Data.EmailAddress,
			reqBody.Data.hashPassword(),
		).Scan(&reqBody.Data.ID)
		if err != nil {
			_ = ErrorResponse(http.StatusInternalServerError, w, err)
			return
		}

		// Register details in kong and get Credentials
		jwtCredentials := &KongJWTCredentials{}
		jwtCredentials, err = kong.CreateConsumerCredentials(reqBody.Data.ID)
		if err != nil {
			_ = ErrorResponse(http.StatusInternalServerError, w, err)
			return
		}

		// Save JWT Credentials inside database
		_, err = db.Exec(
			"UPDATE users SET jwt_credentials = $2 WHERE id = $1",
			reqBody.Data.ID, jwtCredentials)
		if err != nil {
			_ = ErrorResponse(http.StatusInternalServerError, w, err)
			return
		}

		// Generate JWT
		var token string
		token, err = jwtCredentials.GenerateJWT()
		if err != nil {
			_ = ErrorResponse(http.StatusInternalServerError, w, err)
			return
		}

		//send the response
		_ = SuccessResponse(http.StatusOK, w, Credentials(reqBody.Data.ID, token))
	}
}

func Authenticate(db *sql.DB, kong *Kong, validator *validator.Validate) RequestHandler {
	return func(w http.ResponseWriter, req *http.Request) {

		// Get Data form body
		reqBody := ReqBody{}
		err := json.NewDecoder(req.Body).Decode(&reqBody)
		//Could not decode json
		if err != nil {
			_ = ErrorResponse(http.StatusBadRequest, w, err)
			return
		}

		// Validate Data
		// Check validation rules in the struct tags of ReqBody
		err = reqBody.Validate(validator)
		if err != nil {
			_ = ErrorResponse(http.StatusBadRequest, w, err)
			return
		}

		// Find User or return 404

		jwtCredentials := new(KongJWTCredentials)
		var passwordHash string
		err = db.QueryRow(
			"SELECT id, email, password, jwt_credentials FROM users WHERE email = $1",
			reqBody.Data.EmailAddress,
		).Scan(&reqBody.Data.ID, &reqBody.Data.EmailAddress, &passwordHash, &jwtCredentials)
		if err != nil {
			_ = ErrorResponse(http.StatusNotFound, w, err)
			return
		}
		log.Println("User found")

		if !reqBody.Data.PasswordMatches(passwordHash) {
			// return the same error as the previous to prevent bad actors from knowing
			// which of the two submitted fields are wrong
			_ = ErrorResponse(http.StatusNotFound, w, err)
			return
		}
		log.Println("Passwords match")

		// Generate JWT
		var token string
		token, err = jwtCredentials.GenerateJWT()
		if err != nil {
			_ = ErrorResponse(http.StatusInternalServerError, w, err)
			return
		}
		log.Println("Passwords match")

		//send the response
		_ = SuccessResponse(http.StatusOK, w, Credentials(reqBody.Data.ID, token))
	}
}
