package controllers

import (
	"encoding/json"
	"errors"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var GetSurvey = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	survey := models.GetSurveyByTemplateID(uint(id))
	if survey == nil {
		utils.Error(w, "Cannot load survey", errors.New("not found"), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(survey)
}

var SaveSurvey = func(w http.ResponseWriter, r *http.Request) {
	var survey models.Survey
	err := json.NewDecoder(r.Body).Decode(&survey)
	if nil != err {
		utils.Error(w, "Cannot decode survey", err, http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	oldSurvey := models.GetSurveyByTemplateID(uint(id))
	if oldSurvey != nil {
		survey.ID = oldSurvey.ID
	}
	survey.TemplateID = uint(id)
	err = models.SaveSurvey(&survey)
	if nil != err {
		utils.Error(w, "Cannot save survey", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(survey)
}

var DeleteSurvey = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	survey := models.GetSurveyByTemplateID(uint(id))
	if survey == nil {
		utils.Error(w, "Cannot load survey", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteSurvey(survey)
	if err != nil {
		utils.Error(w, "Cannot delete survey", err, http.StatusInternalServerError)
		return
	}
}
