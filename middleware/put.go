package middleware

import (
	"encoding/json"
	"errors"
	"github.com/deploji/deploji-server/utils"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
)

func PutMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		next(rw, r)
		return
	}

	arr := strings.Split(r.URL.Path, "/")
	if len(arr) < 3 {
		next(rw, r)
		return
	}
	id, _ := strconv.ParseUint(arr[2], 10, 16)

	var model gorm.Model
	err := json.NewDecoder(r.Body).Decode(&model)
	if nil != err {
		utils.Error(rw, "Cannot decode model", err, http.StatusBadRequest)
		return
	}

	if model.ID != uint(id) {
		utils.Error(rw, "Updating model ID is forbidden", errors.New(""), http.StatusBadRequest)
		return
	}
	next(rw, r)
}
