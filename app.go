package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/brxyxn/go_mps_redcage/handlers"
	u "github.com/brxyxn/go_mps_redcage/utils"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type App struct {
	Router   *mux.Router
	db       *sql.DB
	l        *log.Logger
	bindAddr string
}

func (a *App) initRoutes() {
	h := handlers.NewHandlers(a.db, a.l)
	// Client routes
	a.Router.HandleFunc("/api/v1/clients", h.CreateClient).Methods("POST")
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}", h.GetClient).Methods("GET")

	// Account routes
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts", h.GetAccounts).Methods("GET")
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts", h.CreateAccount).Methods("POST")
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts/{account_id:[0-9]+}", h.GetAccount).Methods("GET")

	// Transaction routes
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts/{account_id:[0-9]+}/transactions", h.GetTransactions).Methods("GET")
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts/{account_id:[0-9]+}/transactions", h.CreateTransaction).Methods("POST")

	// Serving Documentation Web Server
	// host:port/docs
	opts := middleware.RedocOpts{SpecURL: "/docs/swagger.yaml"}

	sh := middleware.Redoc(opts, nil)

	a.Router.Handle("/docs/swagger.yaml", http.FileServer(http.Dir("./")))
	a.Router.Handle("/docs", sh)
}

/*
To initialize the routes and database connection you must
include the following information as strings and also
call Run setting the port to serve to the web.
*/
func (a *App) Initialize(host, port, user, password, dbname string) {
	// connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbname)
	connectionStr := fmt.Sprintf(
		"host=%s port=%v user=%s "+
			"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	a.db, err = sql.Open("pgx", connectionStr)
	if err != nil {

		u.LogInfo("Error opening a new connection to the DB.", err)
	}

	a.Router = mux.NewRouter() // Make sure this is set before the server is started

	// We try to validate the connection to the DB is correct, otherwise, the app
	// will restart itself, this is a temporary solution because the postgres image usually
	// is initialized after golang's image.
	err = a.db.Ping()
	if err != nil {
		a.db.Close()
		a.l.Fatal(err)
	} else {
		u.LogInfo("(Optional)", "Creating and seeding tables to initializate DB.")

		// Executing SQL statements to create tables and seed DB.
		sqlDir := "db/docker_postgres_init.sql"
		query, err := ioutil.ReadFile(sqlDir)
		if err != nil {
			u.LogError(fmt.Sprintf("Error while reading %s file.", sqlDir), err)
		}

		if _, err := a.db.Exec(string(query)); err != nil {
			a.l.Panic("Unable to run SQL statements.", err)
		}
	}
}

/*
Runs the new server.
*/
func (a *App) Run() {
	// Initializing routes
	a.initRoutes()

	// Creating a new server
	srv := http.Server{
		Addr:         a.bindAddr,        // configure the bind address
		Handler:      a.Router,          // set the default handler
		ErrorLog:     a.l,               // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Starting the server
	go func() {
		u.LogInfo("Running server on port", a.bindAddr)

		err := srv.ListenAndServe()
		if err != nil {
			a.l.Printf("Server Status: %s\n", err)
			os.Exit(1)
		}
	}()

	// Creating channel
	cs := make(chan os.Signal, 1)

	signal.Notify(cs, os.Interrupt, os.Kill)
	// signal.Notify(sigchan, os.Kill) // If running on Windows

	sigchan := <-cs
	u.LogDebug("Signal received:", sigchan)

	ctx, fn := context.WithTimeout(context.Background(), 30*time.Second)
	defer fn()
	srv.Shutdown(ctx)
}
