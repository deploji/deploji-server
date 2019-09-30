package controllers

import (
	"encoding/json"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/utils"
	"log"
	"net/http"
)

var GetSettings = func(w http.ResponseWriter, r *http.Request) {
	settings := models.GetSettingGroups()
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

var SaveSettings = func(w http.ResponseWriter, r *http.Request) {
	var settings []models.Setting
	err := json.NewDecoder(r.Body).Decode(&settings)
	log.Println(err)
	if nil != err {
		utils.Error(w, "Cannot decode setting", err, http.StatusInternalServerError)
		return
	}
	err = models.SaveSettings(&settings)
	if nil != err {
		utils.Error(w, "Cannot save setting", err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}
