package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/brxyxn/go_mps_redcage/data"
	u "github.com/brxyxn/go_mps_redcage/utils"
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
func (a *Handlers) CreateAccount(w http.ResponseWriter, r *http.Request) {
	a.l.Println("Handling POST Accounts /clients/:id/accounts/new")

	clientId, err := strconv.Atoi(mux.Vars(r)["client_id"])
	if err != nil {
		a.l.Println(err)
	}

	var v *data.Account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil {
		a.l.Println("Error decoding data", err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	v.ClientId = uint64(clientId)

	defer r.Body.Close()

	if err := data.CreateAccount(a.db, v); err != nil {
		a.l.Println(err)
		u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	u.RespondWithJSON(w, http.StatusCreated, v)
}

/*
	/api/v1/clients/{client_id}/accounts/{account_id}

This function returns the information of an account based on the {account_id}.

**{client_id} should be replaced by an authentication token (JWT) and validated with middleware.
*/
func (a *Handlers) GetAccount(w http.ResponseWriter, r *http.Request) {
	a.l.Println("Handling GET Account")

	var v *data.Account
	vars := mux.Vars(r)
	accountId, err := strconv.Atoi(vars["account_id"])
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	v = &data.Account{Id: uint64(accountId)}

	if err := data.GetAccount(a.db, v); err != nil {
		switch err {
		case sql.ErrNoRows:
			u.RespondWithError(w, http.StatusNotFound, "This item was not found")
		default:
			u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	u.RespondWithJSON(w, http.StatusOK, v)
}

/*
	/api/v1/clients/{client_id}/accounts

This function returns a list of accounts owned by the {client_id}.
*/
func (a *Handlers) GetAccounts(w http.ResponseWriter, r *http.Request) {
	a.l.Println("Handling GET Accounts")

	// var v *data.Account
	clientId, err := strconv.Atoi(mux.Vars(r)["client_id"])
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, "Invalid client ID")
	}

	// v = &data.Account{ClientId: uint64(clientId)}

	accounts := data.GetAccounts(a.db, clientId)
	if accounts == nil {
		u.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	u.RespondWithJSON(w, http.StatusOK, accounts)
}

/*
	v := Decimal{}

	v.ParseDecimal(fmt.Sprintf("%s %s", p.Currency, p.Balance))

	log.Println(v.Amount, v.Money)
*/
