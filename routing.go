package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/brxyxn/go_mps_redcage/models"
	"github.com/gorilla/mux"
)

type AppRef struct {
	*models.App
}

const (
	host = "queenie.db.elephantsql.com"
	port = "5432"
)

func (a *AppRef) Initialize(user, password, dbname string) {
	connectionStr := fmt.Sprintf("host=%s port=%v user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	a.DB, err = sql.Open("pgx", connectionStr)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	// Initializing routes
	// a.InitAccountRoutes()
	a.InitClientRoutes()
	// a.InitTransactionRoutes()
}

func (a *AppRef) Run(addr string) {
	log.Fatal(http.ListenAndServe(":3005", a.Router))
}

/* Initializing routes to retreive Client's data */
func (a *AppRef) InitClientRoutes() {
	// a.Router.HandleFunc("/api/v1/client", a.GetClients).Methods("GET")
	a.Router.HandleFunc("/api/v1/client/new", a.CreateClient).Methods("POST")
	a.Router.HandleFunc("/api/v1/client/{id:[0-9]+}", a.GetClient).Methods("GET")
	// a.Router.HandleFunc("/api/v1/client/{id:[0-9]+}", a.updateClient).Methods("PUT")    // Create flag
}

/* Initializing routes to retreive Account's data */
func (a *AppRef) InitAccountRoutes() {
	// a.Router.HandleFunc("/api/v1/account/{id:[0-9]+}", a.getAccount).Methods("GET")
	// a.Router.HandleFunc("/api/v1/client/{id:[0-9]+}/account", a.getClientAccounts).Methods("GET")
	// a.Router.HandleFunc("/api/v1/client/{id:[0-9]+}/account/create", a.createClientAccount).Methods("POST")
	// a.Router.HandleFunc("/api/v1/client/{id:[0-9]+}/account/{id:[0-9]+}", a.getClientAccount).Methods("GET")
}

/* Initializing routes to retreive Transaction's data */
func (a *AppRef) InitTransactionRoutes() {
	// a.Router.HandleFunc("/api/v1/account/{id:[0-9]+}/transaction", a.getTransactions).Methods("GET")
	// a.Router.HandleFunc("/api/v1/account/{id:[0-9]+}/transaction/create", a.createTransaction).Methods("POST")
}
