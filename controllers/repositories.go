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

var GetRepositories = func(w http.ResponseWriter, r *http.Request) {
	repositories := models.GetRepositories()
	if repositories == nil {
		utils.Error(w, errors.New("not found"), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repositories)
}

var GetRepository = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	repository := models.GetRepository(uint(id))
	if repository == nil {
		utils.Error(w, errors.New("not found"), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repository)
}

var SaveRepositories = func(w http.ResponseWriter, r *http.Request) {
	var repository models.Repository
	err := json.NewDecoder(r.Body).Decode(&repository)
	log.Println(err)
	if nil != err {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	err = models.SaveRepository(&repository)
	if nil != err {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repository)
}

var DeleteRepository = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	repository := models.GetRepository(uint(id))
	if repository == nil {
		utils.Error(w, errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteRepository(repository)
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}
