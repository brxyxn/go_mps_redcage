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
swagger:route POST /clients clients createClient
Return the id of the new client including the details of it
responses:
	201: clientResponse
	400: badRequestErrorResponse
	500: internalErrorResponse

CreateClient handles POST request and returns the ClientID of the new Client
*/
func (c *Handlers) CreateClient(w http.ResponseWriter, r *http.Request) {
	u.LogInfo("Handling POST Clients /clients")

	var v *data.Client
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil {
		u.LogDebug("GetClient handler:", err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if v.Firstname == "" || v.Lastname == "" || v.Username == "" {
		u.LogDebug("GetClient handler: One or more required params are missing or empty value(s).")
		u.RespondWithError(w, http.StatusBadRequest, "One or more required params are missing or empty value(s).")
		return
	}

	if err := data.CreateClient(c.db, v); err != nil {
		u.LogDebug("GetClient handler:", err)
		u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RespondWithJSON(w, http.StatusCreated, v)
}

/*
swagger:route GET /clients/{client_id} clients getClient
Return the details of one client registry
responses:
	200: clientResponse
	400: badRequestErrorResponse
	404: notFoundErrorResponse
	500: internalErrorResponse

GetClient handles GET requests and returns the details of the requested Client
*/
func (c *Handlers) GetClient(w http.ResponseWriter, r *http.Request) {
	c.l.Println("Handling GET Clients /clients/:client_id")

	// db := c.db

	var values *data.Client

	clientId, err := strconv.Atoi(mux.Vars(r)["client_id"])
	if err != nil {
		u.LogDebug("GetClient handler:", err)
		u.RespondWithError(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	values = &data.Client{Id: uint64(clientId)}

	if err := data.GetClient(c.db, values); err != nil {
		switch err {
		case sql.ErrNoRows:
			u.LogDebug("GetClient handler:", err)
			u.RespondWithError(w, http.StatusNotFound, "This item was not found")
		default:
			u.LogDebug("GetClient handler:", err)
			u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	u.RespondWithJSON(w, http.StatusOK, values)
}
