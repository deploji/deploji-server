package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/models"
	"github.com/sotomskir/mastermind-server/services"
	"github.com/sotomskir/mastermind-server/services/auth"
	"github.com/sotomskir/mastermind-server/utils"
	"net/http"
	"strconv"
)

var GetApplications = func(w http.ResponseWriter, r *http.Request) {
	applications, err := models.GetApplications()
	if err != nil {
		utils.Error(w, "Cannot load applications", err, http.StatusInternalServerError)
		return
	}
	applications = auth.FilterApplications(applications, services.GetJWTClaims(r))
	json.NewEncoder(w).Encode(applications)
}

var GetApplication = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	application := models.GetApplication(uint(id))
	if application == nil {
		utils.Error(w, "Cannot load application", errors.New("not found"), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(application)
}

var GetApplicationInventories = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	inventories := models.GetInventoriesByApplicationId(uint(id))
	if inventories == nil {
		utils.Error(w, "Cannot load inventories", errors.New("not found"), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "inventory/json")
	json.NewEncoder(w).Encode(inventories)
}

var SaveApplications = func(w http.ResponseWriter, r *http.Request) {
	var application models.Application
	err := json.NewDecoder(r.Body).Decode(&application)
	if nil != err {
		utils.Error(w, "Cannot decode application", err, http.StatusInternalServerError)
		return
	}
	err = models.SaveApplication(&application)
	if nil != err {
		utils.Error(w, "Cannot save application", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(application)
}

var DeleteApplication = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	application := models.GetApplication(uint(id))
	if application == nil {
		utils.Error(w, "Cannot load application", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteApplication(application)
	if err != nil {
		utils.Error(w, "Cannot delete application", err, http.StatusInternalServerError)
		return
	}
}
