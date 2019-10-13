package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetJson(url string, target interface{}) error {
	if reflect.TypeOf(http.DefaultTransport).ConvertibleTo(reflect.TypeOf(&http.Transport{})) {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	log.Printf("Request: %s", url)
	r, err := myClient.Get(url)
	if err != nil {
		fmt.Println("get")
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
