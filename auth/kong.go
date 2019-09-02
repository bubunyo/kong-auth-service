package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

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

type KongJWTCredentials struct {
	ConsumerID string `json:"consumer_id"`
	CreatedAt  int64  `json:"created_at"`
	ID         string `json:"id"`
	Key        string `json:"key"`
	Secret     string `json:"secret"`
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

func (k Kong) Status() (D, error) {
	res, err := k.Client.Get(fmt.Sprintf("http://%s:%s", k.Host, k.Port))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var body D

	err = json.NewDecoder(res.Body).Decode(&body)

	if err == nil {
		return nil, err
	}

	return body, nil
}

func (k Kong) CreateUserCredentials(userId string) (*KongJWTCredentials, error) {

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
		fmt.Sprintf("http://%s:%s/consumers/%s/jwt", k.Host, k.Port, c.CustomID),
		"application/json",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var jwt KongJWTCredentials
	err = json.NewDecoder(res.Body).Decode(&jwt)
	if err != nil {
		return nil, err
	}

	return &jwt, nil
}
