package controllers

import (
	"encoding/json"
	"github.com/deploji/deploji-server/dto"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/utils"
	"net/http"
)

var SavePushSubscription = func(w http.ResponseWriter, r *http.Request) {
	var sub dto.PushSubscriptionDTO
	err := json.NewDecoder(r.Body).Decode(&sub)
	if nil != err {
		utils.Error(w, "Cannot decode sub", err, http.StatusInternalServerError)
		return
	}
	jwt := services.GetJWTClaims(r)
	err = models.SavePushSubscription(sub.Endpoint, sub.ExpirationTime, sub.Keys.Auth, sub.Keys.P256dh, jwt.UserID)
	if nil != err {
		utils.Error(w, "Cannot save sub", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(sub)
}
