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
swagger:route POST /clients/{client_id}/accounts accounts createAccount
Return the id of the new account including the details of it
responses:
	201: accountResponse
	400: badRequestErrorResponse
	500: internalErrorResponse

CreateAccount handles POST request and returns the AccountID of the new Account
*/
func (a *Handlers) CreateAccount(w http.ResponseWriter, r *http.Request) {
	u.LogInfo("Handling POST Accounts /clients/:client_id/accounts")

	clientId, err := strconv.Atoi(mux.Vars(r)["client_id"])
	if err != nil {
		u.LogDebug("CreateAccount Handler", err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid client_id")
		return
	}

	var v *data.Account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil {
		u.LogDebug("Error decoding data", err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	v.ClientId = uint64(clientId)

	defer r.Body.Close()

	if err := data.CreateAccount(a.db, v); err != nil {
		u.LogDebug(err)
		u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	u.RespondWithJSON(w, http.StatusCreated, v)
}

/*
swagger:route GET /clients/{client_id}/accounts/{account_id} accounts getAccount
Return the id of the new client
responses:
	200: accountGetResponse
	400: badRequestErrorResponse
	404: notFoundErrorResponse
	500: internalErrorResponse

GetAccount handles GET requests and returns the details of the requested Account
**{client_id} should be replaced by an authentication token (JWT) and validated with middleware.
*/
func (a *Handlers) GetAccount(w http.ResponseWriter, r *http.Request) {
	u.LogInfo("Handling GET Accounts /clients/:client_id/accounts/:account_id")

	var v *data.Account
	vars := mux.Vars(r)
	accountId, err := strconv.Atoi(vars["account_id"])
	if err != nil {
		u.LogDebug("GetAccount Handler: ", err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	v = &data.Account{Id: uint64(accountId)}

	if err := data.GetAccount(a.db, v); err != nil {
		switch err {
		case sql.ErrNoRows:
			u.LogDebug("GetAccount Handler: ", err)
			u.RespondWithError(w, http.StatusNotFound, "This item was not found")
		default:
			u.LogDebug("GetAccount Handler: ", err)
			u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	u.RespondWithJSON(w, http.StatusOK, v)
}

/*
swagger:route GET /clients/{client_id}/accounts accounts getAccounts
Return a list of the accounts listed by the {client_id}
responses:
	200: accountsResponse
	400: badRequestErrorResponse
	404: notFoundErrorResponse

GetAccounts handles GET requests and returns a list of the accounts listed by a client
**{client_id} should be replaced by an authentication token (JWT) and validated with middleware.
*/
func (a *Handlers) GetAccounts(w http.ResponseWriter, r *http.Request) {
	u.LogInfo("Handling GET Accounts /clients/:client_id/accounts")

	clientId, err := strconv.Atoi(mux.Vars(r)["client_id"])
	if err != nil {
		u.LogDebug("GetAccounts Handler:", err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	accounts := data.GetAccounts(a.db, clientId)
	if accounts == nil {
		u.LogDebug("GetAccounts Handler:", err)
		u.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	u.RespondWithJSON(w, http.StatusOK, accounts)
}
