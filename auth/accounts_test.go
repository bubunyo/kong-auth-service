package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/bubunyo/para-services/auth/db"
	"gopkg.in/go-playground/validator.v9"
	"html"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func sanitize(db *sql.DB) {
	_, _ = db.Exec("TRUNCATE users")
}
func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := make(map[string]bool)
	_ = json.Unmarshal([]byte(rr.Body.Bytes()), &expected)

	if !expected["ok"] {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAccountCreationRoute(t *testing.T) {
	const (
		createConsumerResponse = `{
    "custom_id": "123456",
    "created_at": 1567516830,
    "id": "75574239-4553-49c6-a0e2-fde643c57632",
    "tags": null,
    "username": "oneone@gmail.com2"
  }`
	)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && html.EscapeString(r.URL.Path) == "/consumers" {
			_, _ = w.Write([]byte(createConsumerResponse))
			return
		}
		if r.Method == "POST" &&
			html.EscapeString(r.URL.Path) == "/consumers/123456/jwtCredentials" {
			_, _ = w.Write([]byte(createConsumerResponse))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Not Found"))
	})

	// Creating a test client to control output responses
	httpClient, teardown := testingHTTPClient(h)
	defer teardown()

	var (
		conn = db.New()
		kong = &Kong{
			os.Getenv("KONG_HOST"),
			os.Getenv("KONG_ADMIN_PORT"),
			httpClient,
		}
		v = validator.New()
	)

	defer conn.DB.Close()
	sanitize(conn.DB)

	handler := http.HandlerFunc(accountCreationHandler(conn.DB, kong, v))

	reqBody := ReqBody{
		Data: Data{
			EmailAddress: "test@mail.com",
			Password:     "password",
		},
	}

	b, err := json.Marshal(&reqBody)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := CredentialsResponse("", "")
	_ = json.Unmarshal(rr.Body.Bytes(), &expected)

	resp, _ := expected["data"].(map[string]map[string]string)
	if resp["data"]["account_id"] != "" {
		t.Errorf("handler returned unexpected any body: got %v",
			resp["data"]["account_id"])
	}

	credResponse, _ := expected["data"].(map[string]map[string]map[string]string)
	if credResponse["data"]["credentials"]["jwt"] != "" {
		t.Errorf("handler returned unexpected any body: got %v",
			credResponse["data"]["credentials"]["jwt"])
	}
}

func TestAuthenticationRoute(t *testing.T) {
	const (
		createConsumerResponse = `{
    "custom_id": "123456",
    "created_at": 1567516830,
    "id": "75574239-4553-49c6-a0e2-fde643c57632",
    "tags": null,
    "username": "oneone@gmail.com2"
  }`
	)

	var (
		conn      = db.New()
		kong      = &Kong{}
		v         = validator.New()
		testEmail = "test@mail.com"
		password  = "password"
	)

	defer conn.DB.Close()
	sanitize(conn.DB)

	reqBody := ReqBody{
		Data: Data{
			EmailAddress: testEmail,
			Password:     password,
		},
	}

	// Setup User
	_ = conn.DB.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2)",
		testEmail, reqBody.Data.hashPassword(),
	)
	handler := http.HandlerFunc(authenticationHandler(conn.DB, kong, v))

	b, err := json.Marshal(&reqBody)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := CredentialsResponse("", "")
	_ = json.Unmarshal(rr.Body.Bytes(), &expected)

	resp, _ := expected["data"].(map[string]map[string]string)
	if resp["data"]["account_id"] != "" {
		t.Errorf("handler returned unexpected any body: got %v",
			resp["data"]["account_id"])
	}

	credResponse, _ := expected["data"].(map[string]map[string]map[string]string)
	if credResponse["data"]["credentials"]["jwt"] != "" {
		t.Errorf("handler returned unexpected any body: got %v",
			credResponse["data"]["credentials"]["jwt"])
	}
}
