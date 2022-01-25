package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

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
func (a *App) CreateClient(w http.ResponseWriter, r *http.Request) {
	var d Client
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&d); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if d.Firstname == "" || d.Lastname == "" || d.Username == "" {
		RespondWithError(w, http.StatusBadRequest, "One or more required params are missing or empty value(s).")
		return
	}

	if err := d.CreateClient(a.DB); err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusCreated, d)
}

/*
	/api/v1/clients/{client_id}

This function returns the information of a client based on the {client_id}.
*/
func (a *App) GetClient(w http.ResponseWriter, r *http.Request) {
	var d Client

	clientId, err := strconv.Atoi(mux.Vars(r)["client_id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid client ID")
		return
	}

	d = Client{Id: uint64(clientId)}

	if err := d.GetClient(a.DB); err != nil {
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
