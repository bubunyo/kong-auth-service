package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
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

	// Creating Database Connection

	var (
		host     = os.Getenv("DB_HOST")
		dbPort     = 5432
		dbname   = os.Getenv("DB_DATABASE")
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		appPort = os.Getenv("PORT")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbPort, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("auth-app: failed to open DB: %v\n", err)
	}

	r := mux.NewRouter()

	// Redirect routes with trailing slash to the routes without tailing slashes
	// Ref - https://github.com/gorilla/mux/issues/30#issuecomment-21255847
	r = r.StrictSlash(true)

	// Setting the content type of all api responses to json
	r.Use(responseMiddleware)

	// mount Account Routes
	mount(db, r, "/accounts", AccountRoutes)

	// Supplemantary Api Routes
	r.HandleFunc("/healthcheck", healthCheck)


	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", appPort),
		// Set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Run server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Starting server on port %s", appPort)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = db.Close();
	_ = srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}
