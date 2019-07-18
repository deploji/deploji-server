package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/models"
	"log"
	"net/http"
	"strconv"
)

var GetProjects = func(w http.ResponseWriter, r *http.Request) {
	projects := models.GetProjects()
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

var GetProject = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	project := models.GetProject(id)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

var SaveProjects = func(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	log.Println(err)
	if nil != err {
		// Simplified
		log.Println(err)
		return
	}
	err2 := models.SaveProject(&project)
	if nil != err2 {
		// Simplified
		log.Println(err2)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}
