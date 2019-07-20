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

var GetInventories = func(w http.ResponseWriter, r *http.Request) {
	inventories := models.GetInventories()
	if inventories == nil {
		utils.Error(w, errors.New("not found"), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventories)
}

var GetInventory = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	inventory := models.GetInventory(uint(id))
	if inventory == nil {
		utils.Error(w, errors.New("not found"), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

var SaveInventories = func(w http.ResponseWriter, r *http.Request) {
	var inventory models.Inventory
	err := json.NewDecoder(r.Body).Decode(&inventory)
	log.Println(err)
	if nil != err {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	err = models.SaveInventory(&inventory)
	if nil != err {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inventory)
}

var DeleteInventory = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	inventory := models.GetInventory(uint(id))
	if inventory == nil {
		utils.Error(w, errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteInventory(inventory)
	if err != nil {
		utils.Error(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
}
