package controllers

import (
	"encoding/json"
	"errors"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services/auth"
	"github.com/deploji/deploji-server/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var GetRepositories = func(w http.ResponseWriter, r *http.Request) {
	repositories := models.GetRepositories()
	if repositories == nil {
		utils.Error(w, "Cannot load repositories", errors.New("not found"), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(repositories)
}

var GetRepository = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	repository := models.GetRepository(uint(id))
	if repository == nil {
		utils.Error(w, "Cannot load repository", errors.New("not found"), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(repository)
}

var SaveRepositories = func(w http.ResponseWriter, r *http.Request) {
	var repository models.Repository
	err := json.NewDecoder(r.Body).Decode(&repository)
	if nil != err {
		utils.Error(w, "Cannot decode repository", err, http.StatusInternalServerError)
		return
	}
	if !auth.VerifyID(repository.ID, r, w, "id") {
		return
	}
	err = models.SaveRepository(&repository)
	if nil != err {
		utils.Error(w, "Cannot save repository", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(repository)
}

var DeleteRepository = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	repository := models.GetRepository(uint(id))
	if repository == nil {
		utils.Error(w, "Cannot load repository", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteRepository(repository)
	if err != nil {
		utils.Error(w, "Cannot delete repository", err, http.StatusInternalServerError)
		return
	}
}
