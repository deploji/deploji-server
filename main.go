package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sotomskir/mastermind-server/controllers"
	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/ssh-keys", controllers.GetSshKeys).Methods("GET")
	router.HandleFunc("/ssh-keys", controllers.SaveSshKeys).Methods("POST")
	router.HandleFunc("/ssh-keys", controllers.SaveSshKeys).Methods("PUT")
	router.HandleFunc("/ssh-keys/{id}", controllers.DeleteSshKeys).Methods("DELETE")
	router.HandleFunc("/projects", controllers.GetProjects).Methods("GET")
	router.HandleFunc("/projects/{id}", controllers.GetProject).Methods("GET")
	router.HandleFunc("/projects/{id}/synchronize", controllers.SynchronizeProject).Methods("POST")
	router.HandleFunc("/projects", controllers.SaveProjects).Methods("POST")
	router.HandleFunc("/projects", controllers.SaveProjects).Methods("PUT")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	err := http.ListenAndServe(":" + port, router)
	if err != nil {
		fmt.Print(err)
	}
}
