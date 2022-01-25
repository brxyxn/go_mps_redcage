package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
	/api/v1/clients/{client_id}/accounts/{account_id}/transactions/new

This function is used to create a transaction registry based on the
transactionType field parsed in the JSON request.

Note: if the transaction is a deposit or a withdraw
then "senderAccountId" must be 0 (zero).

JSON templates:
	"deposit"
	{
		"amount": 100000.00,
		"transactionType": 1,
		"description": "This is a deposit",
		"receiverAccountId": 1,
		"senderAccountId": 0
	}

	"withdraw"
	{
		"amount": 120.00,
		"transactionType": 2,
		"description": "This is a withdraw",
		"receiverAccountId": 2,
		"senderAccountId": 0
	}

	"transfer"
	{
		"amount": 100.00,
		"transactionType": 3,
		"description": "This is a transfer or EFT",
		"receiverAccountId": 2,
		"senderAccountId": 1
	}

*/
func (a *App) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var d Transaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&d); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := d.CreateTransaction(a.DB); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, d)
}

/*
	/api/v1/clients/{client_id}/accounts/{account_id}/transactions

This function returns a list of transactions from an account.
*/
func (a *App) GetTransactions(w http.ResponseWriter, r *http.Request) {
	var d Transaction

	accountId, err := strconv.Atoi(mux.Vars(r)["account_id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid account ID")
	}

	d.ReceiverAccountId = uint64(accountId)

	transactions, err := d.GetTransactions(a.DB)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, transactions)
}
