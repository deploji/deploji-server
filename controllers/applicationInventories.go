package controllers

import (
	"errors"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var DeleteApplicationInventory = func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseUint(vars["id"], 10, 16)
	applicationInventory := models.GetApplicationInventory(uint(id))
	if applicationInventory == nil {
		utils.Error(w, "Cannot load applicationInventory", errors.New("not found"), http.StatusNotFound)
		return
	}
	err := models.DeleteApplicationInventory(applicationInventory)
	if err != nil {
		utils.Error(w, "Cannot delete applicationInventory", err, http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "applicationInventory/json")
}
