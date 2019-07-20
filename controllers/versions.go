package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/services"
	"github.com/sotomskir/mastermind-server/utils"
	"net/http"
	"strconv"
)

var GetVersions = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["app"], 10, 16)
	versions, err := services.GetVersions(uint(id))
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versions)
}
