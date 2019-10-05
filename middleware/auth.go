package middleware

import (
	"errors"
	"github.com/deploji/deploji-server/dto"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/services/auth"
	"github.com/deploji/deploji-server/utils"
	"net/http"
	"strconv"
	"strings"
)

func AuthMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	arr := strings.Split(r.URL.Path, "/")
	if len(arr) < 3 {
		next(rw, r)
		return
	}
	id, _ := strconv.ParseUint(arr[2], 10, 16)
	user := services.GetJWTClaims(r)
	if !auth.Enforce(user, dto.ObjectType(arr[1]), uint(id), dto.ActionType(r.Method)) {
		utils.Error(rw, "Forbidden", errors.New(""), http.StatusForbidden)
		return
	}
	next(rw, r)
}
