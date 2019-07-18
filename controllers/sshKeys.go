package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/models"
	"log"
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
	log.Println(err)
	if nil != err {
		// Simplified
		log.Println(err)
		return
	}
	models.SaveSshKey(&key)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(key)
}

var DeleteSshKeys = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	key := models.GetSshKey(id)
	models.DeleteSshKey(key)
	w.Header().Add("Content-Type", "application/json")
}
