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
swagger:route POST /clients/{client_id}/accounts/new accounts createAccount
Return the id of the new account including the details of it
responses:
	201: accountResponse
	400: badRequestErrorResponse
	500: internalErrorResponse

CreateAccount handles POST request and returns the AccountID of the new Account
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
swagger:route GET /clients/{client_id}/accounts/{account_id} accounts getAccount
Return the id of the new client
responses:
	200: accountResponse
	400: badRequestErrorResponse
	404: notFoundErrorResponse
	500: internalErrorResponse

GetAccount handles GET requests and returns the details of the requested Account
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
swagger:route GET /clients/{client_id}/accounts/ accounts getAccounts
Return a list of the accounts listed by the {client_id}
responses:
	200: accountsResponse
	400: badRequestErrorResponse
	404: notFoundErrorResponse

GetAccounts handles GET requests and returns a list of the accounts listed by a client
**{client_id} should be replaced by an authentication token (JWT) and validated with middleware.
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
