package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/models"
	"github.com/sotomskir/mastermind-server/utils"
	"log"
	"net/http"
	"strconv"
)

var GetApplications = func(w http.ResponseWriter, r *http.Request) {
	applications, err := models.GetApplications()
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applications)
}

var GetApplication = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	application := models.GetApplication(uint(id))
	if application == nil {
		utils.Error(w, errors.New("not found"), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(application)
}

var SaveApplications = func(w http.ResponseWriter, r *http.Request) {
	var application models.Application
	err := json.NewDecoder(r.Body).Decode(&application)
	log.Println(err)
	if nil != err {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	err = models.SaveApplication(&application)
	if nil != err {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(application)
}

var DeleteApplication = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	application := models.GetApplication(uint(id))
	if application == nil {
		utils.Error(w, errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteApplication(application)
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}
