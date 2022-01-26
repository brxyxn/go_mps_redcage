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

	// connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, dbname)
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
	a.initRoutes()

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

func (a *App) initRoutes() {
	a.Router.HandleFunc("/api/v1/clients/new", a.CreateClient).Methods("POST")            // done
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}", a.GetClient).Methods("GET") // done

	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts", a.GetAccounts).Methods("GET")
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts/new", a.CreateAccount).Methods("POST")
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts/{account_id:[0-9]+}", a.GetAccount).Methods("GET")

	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts/{account_id:[0-9]+}/transactions", a.GetTransactions).Methods("GET")
	a.Router.HandleFunc("/api/v1/clients/{client_id:[0-9]+}/accounts/{account_id:[0-9]+}/transactions/new", a.CreateTransaction).Methods("POST")
}

/* This function handles the creation of tables in DB. */
// func (a *App) createTable(query string) {
// 	_, err := a.DB.Exec(query)
// 	LogErrorMsg(err, "Error creating table.")
// }

// func (a *App) seedTables(query string) {
// 	_, err := a.DB.Exec(query)
// 	LogErrorMsg(err, "Error seeding tables.")
// }
