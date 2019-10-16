package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/services/auth"
	"github.com/deploji/deploji-server/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var GetUsers = func(w http.ResponseWriter, r *http.Request) {
	users := models.GetUsers()
	if users == nil {
		utils.Error(w, "Cannot load user", errors.New("not found"), http.StatusNotFound)
		return
	}
	for _, user := range users {
		user.Password = ""
	}
	json.NewEncoder(w).Encode(users)
}

var GetUser = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	user := models.GetUser(uint(id))
	if user == nil {
		utils.Error(w, "Cannot load user", errors.New("not found"), http.StatusNotFound)
		return
	}
	user.Password = ""
	json.NewEncoder(w).Encode(user)
}

var SaveUser = func(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	log.Println(err)
	if nil != err {
		utils.Error(w, "Cannot decode user", err, http.StatusInternalServerError)
		return
	}
	user.Password, err = services.HashPassword(user.Password)
	if err != nil {
		utils.Error(w, "Error hashing password", fmt.Errorf(""), http.StatusInternalServerError)
		return
	}
	if !auth.VerifyID(user.ID, r, w, "id") {
		return
	}
	err = models.SaveUser(&user)
	if nil != err {
		utils.Error(w, "Cannot save user", err, http.StatusInternalServerError)
		return
	}
	user.Password = ""
	json.NewEncoder(w).Encode(user)
}

var DeleteUser = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	user := models.GetUser(uint(id))
	if user == nil {
		utils.Error(w, "Cannot load user", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteUser(user)
	if err != nil {
		utils.Error(w, "Cannot delete user", err, http.StatusInternalServerError)
		return
	}
}
