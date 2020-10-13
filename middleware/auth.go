package middleware

import (
	"errors"
	"github.com/deploji/deploji-server/dto"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/services/auth"
	"github.com/deploji/deploji-server/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func AuthMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := services.GetJWTClaims(r)
	if user.Type == models.UserTypeAuditor &&
		(r.Method == http.MethodPost || r.Method == http.MethodPut ||r.Method == http.MethodPatch ||r.Method == http.MethodDelete) {
		utils.Error(rw, "Forbidden", errors.New(""), http.StatusForbidden)
		return
	}
	arr := strings.Split(r.URL.Path, "/")
	if len(arr) != 3 {
		next(rw, r)
		return
	}
	id, _ := strconv.ParseUint(arr[2], 10, 16)
	objectType := dto.ObjectType(arr[1])
	if !auth.Enforce(user, objectType, uint(id), getActionType(r.Method)) {
		logrus.Infof("AuthMiddleware Forbidden %d", len(arr))
		utils.Error(rw, "Forbidden", errors.New(""), http.StatusForbidden)
		return
	}
	next(rw, r)
}

func getActionType(method string) dto.ActionType {
	switch method {
	case http.MethodGet:
		return dto.ActionTypeRead
	case http.MethodPost:
		return dto.ActionTypeWrite
	case http.MethodPut:
		return dto.ActionTypeWrite
	case http.MethodPatch:
		return dto.ActionTypeWrite
	default:
		return dto.ActionTypeWrite
	}
}
