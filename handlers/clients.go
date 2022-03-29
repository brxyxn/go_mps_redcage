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
	/api/v1/clients/new

JSON Template:
	{
		"firstName": "John",
		"lastName": "Doe",
		"userName": "john.doe"
	}

*/
func (c *Handlers) CreateClient(w http.ResponseWriter, r *http.Request) {
	c.l.Println("Handle POST Clients")

	var v *data.Client
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&v); err != nil {
		u.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if v.Firstname == "" || v.Lastname == "" || v.Username == "" {
		u.RespondWithError(w, http.StatusBadRequest, "One or more required params are missing or empty value(s).")
		return
	}

	if err := data.CreateClient(c.db, v); err != nil {
		u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	u.RespondWithJSON(w, http.StatusCreated, v)
}

/*
	/api/v1/clients/{client_id}

This function returns the information of a client based on the {client_id}.
*/
func (c *Handlers) GetClient(w http.ResponseWriter, r *http.Request) {
	c.l.Println("Handling GET Clients")

	// db := c.db

	var values *data.Client

	clientId, err := strconv.Atoi(mux.Vars(r)["client_id"])
	if err != nil {
		u.RespondWithError(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	values = &data.Client{Id: uint64(clientId)}

	if err := data.GetClient(c.db, values); err != nil {
		switch err {
		case sql.ErrNoRows:
			u.RespondWithError(w, http.StatusNotFound, "This item was not found")
		default:
			u.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	u.RespondWithJSON(w, http.StatusOK, values)
}
