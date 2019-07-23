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
	router.HandleFunc("/deployments", controllers.GetDeployments).Methods("GET")
	router.HandleFunc("/deployments", controllers.SaveDeployments).Methods("POST")
	router.HandleFunc("/deployments/{id}", controllers.GetDeployment).Methods("GET")
	router.HandleFunc("/deployments/{id}/logs", controllers.GetDeploymentLogs).Methods("GET")
	router.HandleFunc("/inventories", controllers.GetInventories).Methods("GET")
	router.HandleFunc("/inventories", controllers.SaveInventories).Methods("POST")
	router.HandleFunc("/inventories", controllers.SaveInventories).Methods("PUT")
	router.HandleFunc("/inventories/{id}", controllers.GetInventory).Methods("GET")
	router.HandleFunc("/inventories/{id}", controllers.DeleteInventory).Methods("DELETE")
	router.HandleFunc("/repositories", controllers.GetRepositories).Methods("GET")
	router.HandleFunc("/repositories", controllers.SaveRepositories).Methods("POST")
	router.HandleFunc("/repositories", controllers.SaveRepositories).Methods("PUT")
	router.HandleFunc("/repositories/{id}", controllers.GetRepository).Methods("GET")
	router.HandleFunc("/repositories/{id}", controllers.DeleteRepository).Methods("DELETE")
	router.HandleFunc("/projects", controllers.GetProjects).Methods("GET")
	router.HandleFunc("/projects", controllers.SaveProjects).Methods("POST")
	router.HandleFunc("/projects", controllers.SaveProjects).Methods("PUT")
	router.HandleFunc("/projects/{id}", controllers.GetProject).Methods("GET")
	router.HandleFunc("/projects/{id}", controllers.DeleteProject).Methods("DELETE")
	router.HandleFunc("/projects/{id}/synchronize", controllers.SynchronizeProject).Methods("POST")
	router.HandleFunc("/projects/{id}/files", controllers.GetProjectFiles).Methods("GET")
	router.HandleFunc("/applications", controllers.GetApplications).Methods("GET")
	router.HandleFunc("/applications", controllers.SaveApplications).Methods("POST")
	router.HandleFunc("/applications", controllers.SaveApplications).Methods("PUT")
	router.HandleFunc("/applications/{id}", controllers.GetApplication).Methods("GET")
	router.HandleFunc("/applications/{id}", controllers.DeleteApplication).Methods("DELETE")
	router.HandleFunc("/versions", controllers.GetVersions).Queries("app", "{app}").Methods("GET")

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
