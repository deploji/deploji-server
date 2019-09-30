package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/deploji/deploji-server/dto"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/utils"
	"log"
	"net/http"
)

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var credentials dto.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil || credentials.Username == "" || credentials.Password == "" {
		utils.Error(w, "Bad request", fmt.Errorf("bad request"), http.StatusBadRequest)
		return
	}
	user := models.GetUserByUsername(credentials.Username)
	if user == nil || user.IsActive == false {
		utils.Error(w, "Unauthorized", fmt.Errorf("user not found or inactive"), http.StatusBadRequest)
		return
	}

	if models.GetSettingBoolValue("LDAP", "enabled", false) {
		log.Println("LDAP enabled")
		authenticated, _ := services.AuthenticateLDAP(user, credentials.Password)
		if authenticated == true {
			token, err := services.GenerateToken(user)
			if err != nil {
				utils.Error(w, "JWT error", err, http.StatusBadRequest)
				return
			}
			json.NewEncoder(w).Encode(token)
			return
		}
	}

	authenticated, _ := services.AuthenticateDatabase(user, credentials.Password)
	if authenticated == true {
		token, err := services.GenerateToken(user)
		if err != nil {
			utils.Error(w, "JWT error", err, http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(token)
		return
	}
	utils.Error(w, "Unauthorized", err, http.StatusBadRequest)
}

var Refresh = func(w http.ResponseWriter, r *http.Request) {
	token, err := services.RefreshToken(r)
	if err != nil {
		log.Printf("JWT error: %s", err)
		utils.Error(w, "Unauthorized", err, http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(token)
}
