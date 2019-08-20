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

var GetJobs = func(w http.ResponseWriter, r *http.Request) {
	page := utils.NewPage(r)
	filters := utils.NewFilters(r, []string{"application_id", "inventory_id"})
	jobs, paginator := models.GetJobs(page, filters)
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("X-Total-Count", fmt.Sprintf("%d", paginator.TotalRecord))
	json.NewEncoder(w).Encode(jobs)
}

var GetLatestDeployments = func(w http.ResponseWriter, r *http.Request) {
	deployments := models.GetLatestDeployments()
	w.Header().Add("Content-Type", "application/json")
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
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

var SaveJobs = func(w http.ResponseWriter, r *http.Request) {
	var job models.Job
	err := json.NewDecoder(r.Body).Decode(&job)
	if nil != err {
		utils.Error(w, "Cannot decode job", err, http.StatusInternalServerError)
		return
	}
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

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

var GetJobLogs = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	jobLogs := models.GetJobLogs(uint(id))
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobLogs)
}
