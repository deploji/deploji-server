package controllers

import (
	"encoding/json"
	"errors"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/services/auth"
	"github.com/deploji/deploji-server/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var GetTemplates = func(w http.ResponseWriter, r *http.Request) {
	jwt := services.GetJWTClaims(r)
	templates := auth.FilterTemplates(models.GetTemplates(), jwt)
	json.NewEncoder(w).Encode(templates)
}

var GetTemplate = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	template := models.GetTemplate(uint(id))
	if template == nil {
		utils.Error(w, "Cannot load template", errors.New("not found"), http.StatusNotFound)
		return
	}
	auth.InsertTemplatePermissions(template, services.GetJWTClaims(r))
	json.NewEncoder(w).Encode(template)
}

var SaveTemplate = func(w http.ResponseWriter, r *http.Request) {
	var template models.Template
	err := json.NewDecoder(r.Body).Decode(&template)
	log.Println(err)
	if nil != err {
		utils.Error(w, "Cannot decode template", err, http.StatusInternalServerError)
		return
	}
	if !auth.VerifyID(template.ID, r, w, "id") {
		return
	}
	err = models.SaveTemplate(&template)
	if nil != err {
		utils.Error(w, "Cannot save template", err, http.StatusInternalServerError)
		return
	}
	auth.AddOwnerPermissions(r, template)
	json.NewEncoder(w).Encode(template)
}

var DeleteTemplate = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	template := models.GetTemplate(uint(id))
	if template == nil {
		utils.Error(w, "Cannot load template", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteTemplate(template)
	if err != nil {
		utils.Error(w, "Cannot delete template", err, http.StatusInternalServerError)
		return
	}
}

var GetTemplateNotifications = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	notifications := auth.FilterTemplateNotifications(models.GetTemplateNotifications(uint(id)), services.GetJWTClaims(r))
	if notifications == nil {
		utils.Error(w, "Cannot load notifications", errors.New("not found"), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "inventory/json")
	json.NewEncoder(w).Encode(notifications)
}

var SaveTemplateNotification = func(w http.ResponseWriter, r *http.Request) {
	var templateNotification models.TemplateNotification
	err := json.NewDecoder(r.Body).Decode(&templateNotification)
	if nil != err {
		utils.Error(w, "Cannot decode templateNotification", err, http.StatusInternalServerError)
		return
	}
	if err := models.SaveTemplateNotification(&templateNotification); nil != err {
		utils.Error(w, "Cannot save templateNotification", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(templateNotification)
}
