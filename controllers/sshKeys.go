package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/models"
	"github.com/sotomskir/mastermind-server/utils"
	"net/http"
	"strconv"
)

var GetSshKeys = func(w http.ResponseWriter, r *http.Request) {
	keys := models.GetSshKeys()
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

var SaveSshKeys = func(w http.ResponseWriter, r *http.Request) {
	var key models.SshKey
	err := json.NewDecoder(r.Body).Decode(&key)
	if nil != err {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	key = *models.SaveSshKey(&key)
	if nil != &key {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

var DeleteSshKeys = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	key := models.GetSshKey(id)
	if key == nil {
		utils.Error(w, errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteSshKey(key)
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}
