package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
)

/*
To initialize the routes and database connection you must
include the following information as strings and also
call Run setting the port to serve to the web.
*/
func (a *App) Initialize(host, port, user, password, dbname string) {
	connectionStr := fmt.Sprintf(
		"host=%s port=%v user=%s "+
			"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	a.DB, err = sql.Open("pgx", connectionStr)
	LogError(err)

	a.Router = mux.NewRouter()

	// Initializing routes
	a.initAccountRoutes()
	a.initClientRoutes()
	a.initTransactionRoutes()

	// Creating tables if not exists
	a.createTable(clientQuery)
	a.createTable(accountQuery)
	a.createTable(transactionQuery)

	log.Println("API initialized!")
}

/*
In order to run the application and serve it to the web
add the port, eg: Run("3000").
Don't include colons.
*/
func (a *App) Run(port string) {
	log.Printf("Running at localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), a.Router))
}

/* Initializing routes to send/retreive Client's data. */
func (a *App) initClientRoutes() {
	a.Router.HandleFunc("/api/v1/clients/new", a.CreateClient).Methods("POST")            // done
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}", a.GetClient).Methods("GET") // done
}

/* Initializing routes to send/retreive Account's data. */
func (a *App) initAccountRoutes() {
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts", a.GetAccounts).Methods("GET")
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts/new", a.CreateAccount).Methods("POST")
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts/{account_id:[0-9]+}", a.GetAccount).Methods("GET")
}

/* Initializing routes to send/retreive Transaction's data. */
func (a *App) initTransactionRoutes() {
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts/{account_id:[0-9]+}/transactions", a.GetTransactions).Methods("GET")
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts/{account_id:[0-9]+}/transactions/new", a.CreateTransaction).Methods("POST")
}

/* This function handles the creation of tables in DB. */
func (a *App) createTable(query string) {
	_, err := a.DB.Exec(query)
	LogErrorMsg(err, "Error creating table.")
}
