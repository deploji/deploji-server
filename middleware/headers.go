package middleware

import (
	"net/http"
)

func HeadersMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw.Header().Add("Content-Type", "application/json")
	next(rw, r)
}
