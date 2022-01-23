package models

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	model "github.com/brxyxn/go_mps_redcage/structs"
	"github.com/brxyxn/go_mps_redcage/utils"
	"github.com/gorilla/mux"
)

func (a *App) CreateClient(w http.ResponseWriter, r *http.Request) {
	var p ClientRef
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p.Client); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := p.CreateClient(a.DB); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, p)
}

// /api/v1/client/{id}
func (a *App) GetClient(w http.ResponseWriter, r *http.Request) {
	var p ClientRef
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid client ID")
		return
	}
	p.Client = &model.Client{Id: uint64(id)}

	if err := p.GetClient(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			utils.RespondWithError(w, http.StatusNotFound, "Client not found")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, p)
}
