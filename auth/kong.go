package main

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const JwtValidityHours = 3

type KongConsumer struct {
	Username string `json:"username"`
	CustomID string `json:"custom_id"`
	ID       string `json:"id"`
}

// KongStatus returns the status of the kong server
type Kong struct {
	Host   string
	Port   string
	Client *http.Client
}

func NewKong() *Kong {
	return &Kong{
		os.Getenv("KONG_HOST"),
		os.Getenv("KONG_ADMIN_PORT"),
		&http.Client{
			Timeout: time.Second * 10,
		},
	}
}

type KongJWTCredentials struct {
	ConsumerID string `json:"consumer_id"`
	CreatedAt  int64  `json:"created_at"`
	ID         string `json:"id"`
	Key        string `json:"key"`
	Secret     string `json:"secret"`
}
func (j KongJWTCredentials) Value() (driver.Value, error) {
	return json.Marshal(j)
}
func (j *KongJWTCredentials) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &j)
}

func (k KongJWTCredentials) GenerateJWT() (string, error) {
	expiresAt := time.Now().Add(time.Hour * JwtValidityHours).Unix()
	token := jwt.New(jwt.SigningMethodHS256)

	token.Claims = &jwt.StandardClaims{
		ExpiresAt: expiresAt,
		Id:        k.ConsumerID,
	}

	tokenString, err := token.SignedString([]byte(k.Secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (k Kong) CreateConsumerCredentials(userId string) (*KongJWTCredentials, error) {

	b := new(bytes.Buffer)

	c := KongConsumer{CustomID: userId}

	// Creating the consumer on Kong
	_ = json.NewEncoder(b).Encode(c)
	res, err := k.Client.Post(
		fmt.Sprintf("http://%s:%s/consumers", k.Host, k.Port),
		"application/json",
		b,
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&c)

	if err != nil {
		return nil, err
	}

	// Generating Credentials
	_ = json.NewEncoder(b).Encode(c)
	res, err = k.Client.Post(
		fmt.Sprintf("http://%s:%s/consumers/%s/jwtCredentials", k.Host, k.Port, c.CustomID),
		"application/json",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var jwtCredentials KongJWTCredentials
	err = json.NewDecoder(res.Body).Decode(&jwtCredentials)
	if err != nil {
		return nil, err
	}

	return &jwtCredentials, nil
}
