package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/services"
	"github.com/sotomskir/mastermind-server/utils"
	"log"
	"net/http"
	"strconv"
)

var GetVersions = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("vars: %#v", vars)
	log.Printf("vars: %#v", r)

	id, _ := strconv.ParseUint(vars["app"], 10, 16)
	versions, err := services.GetVersions(uint(id))
	if err != nil {
		utils.Error(w, "Cannot load versions", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(versions)
}
