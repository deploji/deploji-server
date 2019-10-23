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

var GetSurveyInputs = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	surveyInputs, err := models.GetSurveyInputsByTemplateID(uint(id))
	if err != nil {
		utils.Error(w, "", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(surveyInputs)
}

var GetSurveyInput = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["inputId"], 10, 16)
	surveyInput := models.GetSurveyInput(uint(id))
	if surveyInput == nil {
		utils.Error(w, "Cannot load surveyInput", errors.New("not found"), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(surveyInput)
}

var SaveSurveyInput = func(w http.ResponseWriter, r *http.Request) {
	var surveyInput models.SurveyInput
	err := json.NewDecoder(r.Body).Decode(&surveyInput)
	if nil != err {
		utils.Error(w, "Cannot decode surveyInput", err, http.StatusInternalServerError)
		return
	}
	if !auth.VerifyID(surveyInput.ID, r, w, "inputId") {
		return
	}
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	survey := models.GetSurveyByTemplateID(uint(id))
	if survey == nil {
		utils.Error(w, "survey not found", errors.New(""), http.StatusNotFound)
		return
	}
	surveyInput.SurveyID = survey.ID
	err = models.SaveSurveyInput(&surveyInput)
	if nil != err {
		utils.Error(w, "Cannot save surveyInput", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(surveyInput)
}

var DeleteSurveyInput = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["inputId"], 10, 16)
	surveyInput := models.GetSurveyInput(uint(id))
	if surveyInput == nil {
		utils.Error(w, "Cannot load surveyInput", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteSurveyInput(surveyInput)
	if err != nil {
		utils.Error(w, "Cannot delete surveyInput", err, http.StatusInternalServerError)
		return
	}
}
