package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Message(status bool, message string) (map[string]interface{}) {
	return map[string]interface{} {"status" : status, "message" : message}
}

func Respond(w http.ResponseWriter, data map[string] interface{})  {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, err error, status int) {
	logrus.Errorln(err.Error())
	http.Error(w, err.Error(), status)
}
