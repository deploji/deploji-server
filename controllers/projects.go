package controllers

import (
	"encoding/json"
	"errors"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/services/amqpService"
	"github.com/deploji/deploji-server/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var GetProjects = func(w http.ResponseWriter, r *http.Request) {
	projects := models.GetProjects()
	json.NewEncoder(w).Encode(projects)
}

var GetProjectsSyncStatus = func(w http.ResponseWriter, r *http.Request) {
	jobs := models.GetLatestSCMPulls()
	json.NewEncoder(w).Encode(jobs)
}

var GetProject = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	project := models.GetProject(uint(id))
	if project == nil {
		utils.Error(w, "Cannot load project", errors.New("not found"), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(project)
}

var GetProjectFiles = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	files, err := services.GetProjectFiles(uint(id))
	if err != nil {
		utils.Error(w, "Cannot load project files", err, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(files)
}

var SaveProjects = func(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	log.Println(err)
	if nil != err {
		utils.Error(w, "Cannot decode project", err, http.StatusInternalServerError)
		return
	}
	err = models.SaveProject(&project)
	if nil != err {
		utils.Error(w, "Cannot save project", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(project)
}

var DeleteProject = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	project := models.GetProject(uint(id))
	if project == nil {
		utils.Error(w, "Cannot load project", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteProject(project)
	if err != nil {
		utils.Error(w, "Cannot delete project", err, http.StatusInternalServerError)
		return
	}
}

var SynchronizeProject = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	job := models.Job{ProjectID: uint(id), Type: models.TypeSCMPull}
	if err := models.SaveJob(&job); err != nil {
		utils.Error(w, "Cannot save job", err, http.StatusInternalServerError)
		return
	}
	if err := amqpService.SendJob(job.ID, models.TypeSCMPull); err != nil {
		utils.Error(w, "Cannot send job", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(job)
}
