package controllers

import (
	"encoding/json"
	"github.com/SherClockHolmes/webpush-go"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/utils"
	"log"
	"net/http"
)

var GetSettings = func(w http.ResponseWriter, r *http.Request) {
	settings := models.GetSettingGroups()
	json.NewEncoder(w).Encode(settings)
}

var GetFrontSettings = func(w http.ResponseWriter, r *http.Request) {
	settings := make(map[string]map[string]interface{}, 0)
	settings["PUSH"] = make(map[string]interface{}, 0)
	settings["PUSH"]["publicKey"] = models.GetSettingValue("PUSH", "publicKey", "")
	json.NewEncoder(w).Encode(settings)
}

var GenerateVapidKeys = func(w http.ResponseWriter, r *http.Request) {
	privateKey, publicKey, _ := webpush.GenerateVAPIDKeys()
	keys := make(map[string]string, 0)
	keys["PrivateKey"] = privateKey
	keys["PublicKey"] = publicKey
	json.NewEncoder(w).Encode(keys)
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
	json.NewEncoder(w).Encode(settings)
}
