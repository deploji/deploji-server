package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/dto"
	"github.com/sotomskir/mastermind-server/services/auth"
	"github.com/sotomskir/mastermind-server/utils"
	"net/http"
	"strconv"
)

var GetGroups = func(w http.ResponseWriter, r *http.Request) {
	groups := auth.GetGroups()
	w.Header().Add("Content-Type", "group/json")
	json.NewEncoder(w).Encode(groups)
}

var GetGroup = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := vars["id"]
	group := auth.GetGroup(id)
	if group == nil {
		utils.Error(w, "Cannot load group", errors.New("not found"), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "group/json")
	json.NewEncoder(w).Encode(group)
}

var SaveGroups = func(w http.ResponseWriter, r *http.Request) {
	var group dto.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if nil != err {
		utils.Error(w, "Cannot decode group", err, http.StatusInternalServerError)
		return
	}
	err = auth.SaveGroup(&group)
	if nil != err {
		utils.Error(w, "Cannot save group", err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "group/json")
	json.NewEncoder(w).Encode(group)
}

var DeleteGroup = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	group := auth.GetGroup(uint(id))
	if group == nil {
		utils.Error(w, "Cannot load group", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := auth.DeleteGroup(group)
	if err != nil {
		utils.Error(w, "Cannot delete group", err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "group/json")
}
