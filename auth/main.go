package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bubunyo/para-services/auth/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type MountPoint func(r *mux.Router, db *sql.DB)

func mount(db *sql.DB, r *mux.Router, path string, mp MountPoint) {
	mp(r.PathPrefix(path).Subrouter(), db)
}

func responseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Setting responses to json
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	// Redirect routes with trailing slash to the routes without tailing slashes
	// Ref - https://github.com/gorilla/mux/issues/30#issuecomment-21255847
	r = r.StrictSlash(true)
	r.Use(responseMiddleware)

	// Creating Database Connection
	conn := db.New()

	mount(conn.DB, r, "/accounts", AccountRoutes)
	// Health Check
	r.HandleFunc("/healthcheck", healthCheck)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		// Set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	_ = conn.DB.Close();
	_ = srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}
