package controllers

import (
	"encoding/json"
	"errors"
	"github.com/deploji/deploji-server/dto"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/services/auth"
	"github.com/deploji/deploji-server/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var GetTeams = func(w http.ResponseWriter, r *http.Request) {
	teams := models.GetTeams()
	if teams == nil {
		teams = make([]*models.Team, 0)
	}
	json.NewEncoder(w).Encode(teams)
}

var GetTeam = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	team := models.GetTeam(uint(id))
	if team == nil {
		utils.Error(w, "Cannot load team", errors.New("not found"), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(team)
}

var GetTeamUsers = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	users := auth.GetUsersForTeam(vars["id"])
	json.NewEncoder(w).Encode(users)
}

var SaveTeamUser = func(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if nil != err {
		utils.Error(w, "Cannot decode user", err, http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
	user  = *models.GetUser(user.ID)
	err = auth.AddUserToTeam(vars["id"], user.ID)
	if nil != err {
		utils.Error(w, "Error adding user: ", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

var DeleteTeamUser = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	err := auth.RemoveUserFromTeam(vars["id"], vars["userId"])
	if nil != err {
		utils.Error(w, "Error removing user: ", err, http.StatusInternalServerError)
		return
	}
}

var GetTeamPermissions = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	permissions := auth.GetPermissionsForTeam(uint(id))
	json.NewEncoder(w).Encode(permissions)
}

var GetUserPermissions = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	permissions := auth.GetPermissionsForUser(uint(id))
	json.NewEncoder(w).Encode(permissions)
}

var SaveTeam = func(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	if nil != err {
		utils.Error(w, "Cannot decode team", err, http.StatusInternalServerError)
		return
	}
	if !auth.VerifyID(team.ID, r) {
		utils.Error(w, "updating model ID is forbidden", errors.New(""), http.StatusForbidden)
		return
	}
	err = models.SaveTeam(&team)
	if nil != err {
		utils.Error(w, "Cannot save team", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(team)
}

var SavePermission = func(w http.ResponseWriter, r *http.Request) {
	var permission dto.Permission
	err := json.NewDecoder(r.Body).Decode(&permission)
	if nil != err {
		utils.Error(w, "Cannot decode permission", err, http.StatusBadRequest)
		return
	}
	if !auth.Enforce(services.GetJWTClaims(r), permission.ObjectType, permission.ObjectID, dto.ActionTypeAdmin) {
		utils.Error(w, "Cannot save team", err, http.StatusForbidden)
		return
	}
	err = auth.AddPermission(permission)
	if nil != err {
		utils.Error(w, "Cannot save team", err, http.StatusInternalServerError)
		return
	}
}

var GetPermissions = func(w http.ResponseWriter, r *http.Request) {
	filters := utils.NewFilters(r, []string{"SubjectType", "ObjectType", "ObjectID"})
	permissions := auth.GetPermissions(filters)
	json.NewEncoder(w).Encode(permissions)
}

var DeletePermission = func(w http.ResponseWriter, r *http.Request) {
	var permission dto.Permission
	err := json.NewDecoder(r.Body).Decode(&permission)
	if nil != err {
		utils.Error(w, "Cannot decode permission", err, http.StatusBadRequest)
		return
	}
	if !auth.Enforce(services.GetJWTClaims(r), permission.ObjectType, permission.ObjectID, dto.ActionTypeAdmin) {
		utils.Error(w, "Cannot save team", err, http.StatusForbidden)
		return
	}
	err = auth.RemovePermission(permission)
	if nil != err {
		utils.Error(w, "Cannot save team", err, http.StatusInternalServerError)
		return
	}
}

var DeleteTeam = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	team := models.GetTeam(uint(id))
	if team == nil {
		utils.Error(w, "Cannot load team", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteTeam(team)
	if err != nil {
		utils.Error(w, "Cannot delete team", err, http.StatusInternalServerError)
		return
	}
}
