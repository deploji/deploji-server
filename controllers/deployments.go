package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/models"
	"github.com/sotomskir/mastermind-server/services/amqpService"
	"github.com/sotomskir/mastermind-server/utils"
	"log"
	"net/http"
	"strconv"
)

var GetDeployments = func(w http.ResponseWriter, r *http.Request) {
	deployments := models.GetDeployments()
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployments)
}

var GetDeployment = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	deployment := models.GetDeployment(uint(id))
	if deployment == nil {
		utils.Error(w, errors.New("not found"), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployment)
}

var SaveDeployments = func(w http.ResponseWriter, r *http.Request) {
	var deployment models.Deployment
	err := json.NewDecoder(r.Body).Decode(&deployment)
	log.Println(err)
	if nil != err {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	err = models.SaveDeployment(&deployment)
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	err = amqpService.Send(deployment)
	if nil != err {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployment)
}

var GetDeploymentLogs = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	deploymentLogs := models.GetDeploymentLogs(uint(id))
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deploymentLogs)
}
