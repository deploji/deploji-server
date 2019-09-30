package controllers

import (
	"encoding/json"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/utils"
	"github.com/gorilla/mux"
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
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versions)
}
