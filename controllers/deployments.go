package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/models"
	"github.com/sotomskir/mastermind-server/services/amqpService"
	"github.com/sotomskir/mastermind-server/utils"
	"net/http"
	"strconv"
)

var GetDeployments = func(w http.ResponseWriter, r *http.Request) {
	page := utils.NewPage(r)
	filters := utils.NewFilters(r, []string{"application_id", "inventory_id"})
	deployments, paginator := models.GetDeployments(page, filters)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("X-Total-Count", fmt.Sprintf("%d", paginator.TotalRecord))
	json.NewEncoder(w).Encode(deployments)
}

var GetLatestDeployments = func(w http.ResponseWriter, r *http.Request) {
	deployments := models.GetLatestDeployments()
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployments)
}

var GetDeployment = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	deployment := models.GetDeployment(uint(id))
	if deployment == nil {
		utils.Error(w, "Cannot load deployment", errors.New("not found"), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployment)
}

var SaveDeployments = func(w http.ResponseWriter, r *http.Request) {
	var deployment models.Deployment
	err := json.NewDecoder(r.Body).Decode(&deployment)
	if nil != err {
		utils.Error(w, "Cannot decode deployment", err, http.StatusInternalServerError)
		return
	}
	err = models.SaveDeployment(&deployment)
	if err != nil {
		utils.Error(w, "Cannot save deployment", err, http.StatusInternalServerError)
		return
	}
	err = amqpService.SendDeployment(deployment)
	if nil != err {
		utils.Error(w, "Cannot send deployment", err, http.StatusInternalServerError)
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
