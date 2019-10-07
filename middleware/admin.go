package middleware

import (
"errors"
"github.com/deploji/deploji-server/models"
"github.com/deploji/deploji-server/services"
"github.com/deploji/deploji-server/utils"
"net/http"
)

func AdminOnlyMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user := services.GetJWTClaims(r)
	if user.Type != models.UserTypeAdmin {
		utils.Error(rw, "Forbidden", errors.New(""), http.StatusForbidden)
		return
	}
	next(rw, r)
}
