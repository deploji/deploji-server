package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/services/amqpService"
	"github.com/deploji/deploji-server/services/auth"
	"github.com/deploji/deploji-server/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var GetJobs = func(w http.ResponseWriter, r *http.Request) {
	page := utils.NewPage(r)
	filters := utils.NewFilters(r, []string{"application_id", "inventory_id"})
	jobs, paginator := models.GetJobs(page, filters)
	w.Header().Add("X-Total-Count", fmt.Sprintf("%d", paginator.TotalRecord))
	json.NewEncoder(w).Encode(jobs)
}

var GetLatestDeployments = func(w http.ResponseWriter, r *http.Request) {
	deployments := models.GetLatestDeployments()
	json.NewEncoder(w).Encode(deployments)
}

var GetJob = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	job := models.GetJob(uint(id))
	if job == nil {
		utils.Error(w, "Cannot load job", errors.New("not found"), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(job)
}

var SaveJobs = func(w http.ResponseWriter, r *http.Request) {
	var job models.Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if nil != err {
		utils.Error(w, "Cannot decode job", err, http.StatusInternalServerError)
		return
	}
	if !auth.VerifyID(job.ID, r) {
		utils.Error(w, "updating model ID is forbidden", errors.New(""), http.StatusForbidden)
		return
	}
	job.UserID = services.GetJWTClaims(r).UserID
	err = models.SaveJob(&job)
	if err != nil {
		utils.Error(w, "Cannot save job", err, http.StatusInternalServerError)
		return
	}
	err = amqpService.SendJob(job.ID, job.Type)
	if nil != err {
		utils.Error(w, "Cannot send job", err, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(job)
}

var GetJobLogs = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	jobLogs := models.GetJobLogs(uint(id))
	json.NewEncoder(w).Encode(jobLogs)
}
