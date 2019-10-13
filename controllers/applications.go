package controllers

import (
	"encoding/json"
	"errors"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/services/auth"
	"github.com/deploji/deploji-server/utils"
	"github.com/gorilla/mux"
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
	auth.InsertApplicationPermissions(application, services.GetJWTClaims(r))
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

var GetApplicationNotifications = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	notifications := models.GetApplicationNotifications(uint(id))
	if notifications == nil {
		utils.Error(w, "Cannot load notifications", errors.New("not found"), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "inventory/json")
	json.NewEncoder(w).Encode(notifications)
}

var SaveApplicationNotification = func(w http.ResponseWriter, r *http.Request) {
	var applicationNotification models.ApplicationNotification
	err := json.NewDecoder(r.Body).Decode(&applicationNotification)
	if nil != err {
		utils.Error(w, "Cannot decode applicationNotification", err, http.StatusInternalServerError)
		return
	}
	if err := models.SaveApplicationNotification(&applicationNotification); nil != err {
		utils.Error(w, "Cannot save applicationNotification", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(applicationNotification)
}

var SaveApplications = func(w http.ResponseWriter, r *http.Request) {
	var application models.Application
	err := json.NewDecoder(r.Body).Decode(&application)
	if nil != err {
		utils.Error(w, "Cannot decode application", err, http.StatusInternalServerError)
		return
	}
	if !auth.VerifyID(application.ID, r) {
		utils.Error(w, "updating model ID is forbidden", errors.New(""), http.StatusForbidden)
		return
	}
	err = models.SaveApplication(&application)
	if nil != err {
		utils.Error(w, "Cannot save application", err, http.StatusInternalServerError)
		return
	}
	auth.AddOwnerPermissions(r, application)
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
