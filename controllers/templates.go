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

var GetTemplates = func(w http.ResponseWriter, r *http.Request) {
	templates := models.GetTemplates()
	w.Header().Add("Content-Type", "application/json")
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
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(template)
}

var SaveTemplates = func(w http.ResponseWriter, r *http.Request) {
	var template models.Template
	err := json.NewDecoder(r.Body).Decode(&template)
	log.Println(err)
	if nil != err {
		utils.Error(w, "Cannot decode template", err, http.StatusInternalServerError)
		return
	}
	err = models.SaveTemplate(&template)
	if nil != err {
		utils.Error(w, "Cannot save template", err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
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
	w.Header().Add("Content-Type", "application/json")
}
