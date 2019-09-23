package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/dto"
	"github.com/sotomskir/mastermind-server/models"
	"github.com/sotomskir/mastermind-server/services/auth"
	"github.com/sotomskir/mastermind-server/utils"
	"log"
	"net/http"
	"strconv"
)

var GetTeams = func(w http.ResponseWriter, r *http.Request) {
	teams := models.GetTeams()
	if teams == nil {
		utils.Error(w, "Cannot load team", errors.New("not found"), http.StatusNotFound)
		return
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
	users := auth.GetUsersForRole(vars["id"])
	json.NewEncoder(w).Encode(users)
}

var SaveTeamUser = func(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	log.Println(err)
	if nil != err {
		utils.Error(w, "Cannot decode user", err, http.StatusInternalServerError)
		return
	}
	vars := mux.Vars(r)
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
	permissions := auth.GetPermissionsForTeam(vars["id"])
	json.NewEncoder(w).Encode(permissions)
}

var SaveTeam = func(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	log.Println(err)
	if nil != err {
		utils.Error(w, "Cannot decode team", err, http.StatusInternalServerError)
		return
	}
	err = models.SaveTeam(&team)
	if nil != err {
		utils.Error(w, "Cannot save team", err, http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(team)
}

var SaveTeamPermission = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var permission dto.Permission
	err := json.NewDecoder(r.Body).Decode(&permission)
	log.Println(err)
	if nil != err {
		utils.Error(w, "Cannot decode permission", err, http.StatusInternalServerError)
		return
	}
	err = auth.AddPermissionToTeam(vars["id"], permission.ObjectType, permission.ObjectID, permission.Role)
	if nil != err {
		utils.Error(w, "Cannot save team", err, http.StatusInternalServerError)
		return
	}
}

var DeleteTeamPermission = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var permission dto.Permission
	err := json.NewDecoder(r.Body).Decode(&permission)
	log.Println(err)
	if nil != err {
		utils.Error(w, "Cannot decode permission", err, http.StatusInternalServerError)
		return
	}
	err = auth.RemovePermissionFromTeam(vars["id"], permission.ObjectType, permission.ObjectID, permission.Role)
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
