package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
	/api/v1/clients/{client_id}/accounts/new

JSON Template:
	"COP & Savings"
	{
		"balance": 80000.00,
		"currency": "COP",
		"accountType": "Savings"
	}

	"USD & Checking"
	{
		"balance": 25000.00,
		"currency": "USD",
		"accountType": "Checking"
	}

	"MXN & CreditCard"
	{
		"balance": 65000.00,
		"currency": "MXN",
		"accountType": "Credit Card"
	}


"accountId"`
    Balance     float64     `json:"balance"`
    Currency    string      `json:"currency"`
    AccountType AccountType `json:"accountType"`
    Active      bool        `json:"active,omitempty"`
    ClientId    uint64      `json:"clientId"`
    CreatedAt   string      `json:"createdAt,omitempty"`
*/
func (a *App) CreateAccount(w http.ResponseWriter, r *http.Request) {
	clientId, err := strconv.Atoi(mux.Vars(r)["client_id"])
	LogError(err)

	var d Account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&d); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	d.ClientId = uint64(clientId)

	defer r.Body.Close()

	if err := d.CreateAccount(a.DB); err != nil {
		log.Println(err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, d)
}

/*
	/api/v1/clients/{client_id}/accounts/{account_id}

This function returns the information of an account based on the {account_id}.

**{client_id} should be replaced by an authentication token (JWT) and validated with middleware.
*/
func (a *App) GetAccount(w http.ResponseWriter, r *http.Request) {
	var d Account
	vars := mux.Vars(r)
	accountId, err := strconv.Atoi(vars["account_id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	d = Account{Id: uint64(accountId)}

	if err := d.GetAccount(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			RespondWithError(w, http.StatusNotFound, "This item was not found")
		default:
			RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	RespondWithJSON(w, http.StatusOK, d)
}

/*
	/api/v1/clients/{client_id}/accounts

This function returns a list of accounts owned by the {client_id}.
*/
func (a *App) GetAccounts(w http.ResponseWriter, r *http.Request) {
	var d Account
	clientId, err := strconv.Atoi(mux.Vars(r)["client_id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid client ID")
	}

	d.ClientId = uint64(clientId)

	accounts, err := d.GetAccounts(a.DB)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, accounts)
}
