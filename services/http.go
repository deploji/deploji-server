package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetJson(url string, target interface{}) error {
	log.Printf("Request: %s", url)
	r, err := myClient.Get(url)
	if err != nil {
		fmt.Println("get")
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
