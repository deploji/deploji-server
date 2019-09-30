package middleware

import (
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/utils"
	"log"
	"net/http"
)

func JwtMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	jwtToken := services.TokenGetter(r)
	token, err := services.ParseToken(jwtToken)
	if err != nil {
		log.Printf("JWT error: %s", err)
		utils.Error(rw, "Unauthorized", err, http.StatusUnauthorized)
		return
	}
	if err := services.VerifyClaims(token, true, true, true); err != nil {
		log.Printf("JWT error: %s", err)
		utils.Error(rw, "Unauthorized", err, http.StatusUnauthorized)
		return
	}
	next(rw, r)
}
