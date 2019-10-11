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

var GetNotificationChannels = func(w http.ResponseWriter, r *http.Request) {
	notificationChannels := models.GetNotificationChannels()
	json.NewEncoder(w).Encode(notificationChannels)
}

var GetNotificationChannel = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	notificationChannel := models.GetNotificationChannel(uint(id))
	if notificationChannel == nil {
		utils.Error(w, "Cannot load notificationChannel", errors.New("not found"), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(notificationChannel)
}

var SaveNotificationChannels = func(w http.ResponseWriter, r *http.Request) {
	var notificationChannel models.NotificationChannel
	err := json.NewDecoder(r.Body).Decode(&notificationChannel)
	if nil != err {
		utils.Error(w, "Cannot decode notificationChannel", err, http.StatusInternalServerError)
		return
	}
	if !auth.VerifyID(notificationChannel.ID, r) {
		utils.Error(w, "updating model ID is forbidden", errors.New(""), http.StatusForbidden)
		return
	}
	err = models.SaveNotificationChannel(&notificationChannel)
	if nil != err {
		utils.Error(w, "Cannot save notificationChannel", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(notificationChannel)
}

var DeleteNotificationChannel = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	notificationChannel := models.GetNotificationChannel(uint(id))
	if notificationChannel == nil {
		utils.Error(w, "Cannot load notificationChannel", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteNotificationChannel(notificationChannel)
	if err != nil {
		utils.Error(w, "Cannot delete notificationChannel", err, http.StatusInternalServerError)
		return
	}
}
