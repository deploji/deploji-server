package main

import (
	"github.com/deploji/deploji-server/controllers"
	"github.com/deploji/deploji-server/services"
	"github.com/deploji/deploji-server/services/amqpService"
	"github.com/deploji/deploji-server/settings"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"golang.org/x/net/context"
	"os"
)

func main() {
	ctx, done := context.WithCancel(context.Background())

	go func() {
		amqpService.Publish(amqpService.Redial(ctx, os.Getenv("AMQP_URL")), amqpService.Jobs, "jobs")
		done()
	}()

	openRouter := mux.NewRouter()
	authRouter := mux.NewRouter()

	authRouter.HandleFunc("/ssh-keys", controllers.GetSshKeys).Methods("GET")
	authRouter.HandleFunc("/ssh-keys", controllers.SaveSshKeys).Methods("POST")
	authRouter.HandleFunc("/ssh-keys", controllers.SaveSshKeys).Methods("PUT")
	authRouter.HandleFunc("/ssh-keys/{id}", controllers.GetSshKey).Methods("GET")
	authRouter.HandleFunc("/ssh-keys/{id}", controllers.DeleteSshKeys).Methods("DELETE")
	authRouter.HandleFunc("/jobs", controllers.GetJobs).Methods("GET")
	authRouter.HandleFunc("/jobs", controllers.SaveJobs).Methods("POST")
	authRouter.HandleFunc("/jobs/latest-deployments", controllers.GetLatestDeployments).Methods("GET")
	authRouter.HandleFunc("/jobs/{id}", controllers.GetJob).Methods("GET")
	authRouter.HandleFunc("/jobs/{id}/logs", controllers.GetJobLogs).Methods("GET")
	authRouter.HandleFunc("/inventories", controllers.GetInventories).Methods("GET")
	authRouter.HandleFunc("/inventories", controllers.SaveInventories).Methods("POST")
	authRouter.HandleFunc("/inventories", controllers.SaveInventories).Methods("PUT")
	authRouter.HandleFunc("/inventories/{id}", controllers.GetInventory).Methods("GET")
	authRouter.HandleFunc("/inventories/{id}", controllers.DeleteInventory).Methods("DELETE")
	authRouter.HandleFunc("/repositories", controllers.GetRepositories).Methods("GET")
	authRouter.HandleFunc("/repositories", controllers.SaveRepositories).Methods("POST")
	authRouter.HandleFunc("/repositories", controllers.SaveRepositories).Methods("PUT")
	authRouter.HandleFunc("/repositories/{id}", controllers.GetRepository).Methods("GET")
	authRouter.HandleFunc("/repositories/{id}", controllers.DeleteRepository).Methods("DELETE")
	authRouter.HandleFunc("/projects", controllers.GetProjects).Methods("GET")
	authRouter.HandleFunc("/projects/synchronize-status", controllers.GetProjectsSyncStatus).Methods("GET")
	authRouter.HandleFunc("/projects", controllers.SaveProjects).Methods("POST")
	authRouter.HandleFunc("/projects", controllers.SaveProjects).Methods("PUT")
	authRouter.HandleFunc("/projects/{id}", controllers.GetProject).Methods("GET")
	authRouter.HandleFunc("/projects/{id}", controllers.DeleteProject).Methods("DELETE")
	authRouter.HandleFunc("/projects/{id}/synchronize", controllers.SynchronizeProject).Methods("POST")
	authRouter.HandleFunc("/projects/{id}/files", controllers.GetProjectFiles).Methods("GET")
	authRouter.HandleFunc("/templates", controllers.GetTemplates).Methods("GET")
	authRouter.HandleFunc("/templates", controllers.SaveTemplates).Methods("POST")
	authRouter.HandleFunc("/templates", controllers.SaveTemplates).Methods("PUT")
	authRouter.HandleFunc("/templates/{id}", controllers.GetTemplate).Methods("GET")
	authRouter.HandleFunc("/templates/{id}", controllers.DeleteTemplate).Methods("DELETE")
	authRouter.HandleFunc("/applications", controllers.GetApplications).Methods("GET")
	authRouter.HandleFunc("/applications", controllers.SaveApplications).Methods("POST")
	authRouter.HandleFunc("/applications", controllers.SaveApplications).Methods("PUT")
	authRouter.HandleFunc("/applications/{id}", controllers.GetApplication).Methods("GET")
	authRouter.HandleFunc("/applications/{id}/inventories", controllers.GetApplicationInventories).Methods("GET")
	authRouter.HandleFunc("/applications/{id}", controllers.DeleteApplication).Methods("DELETE")
	authRouter.HandleFunc("/application-inventories/{id}", controllers.DeleteApplicationInventory).Methods("DELETE")
	authRouter.HandleFunc("/versions", controllers.GetVersions).Queries("app", "{app}").Methods("GET")
	authRouter.HandleFunc("/auth/users", controllers.SaveUser).Methods("POST")
	authRouter.HandleFunc("/auth/users", controllers.SaveUser).Methods("PUT")
	authRouter.HandleFunc("/auth/users", controllers.GetUsers).Methods("GET")
	authRouter.HandleFunc("/auth/users/{id}", controllers.GetUser).Methods("GET")
	authRouter.HandleFunc("/settings", controllers.GetSettings).Methods("GET")
	authRouter.HandleFunc("/settings", controllers.SaveSettings).Methods("PUT")
	openRouter.HandleFunc("/auth/authenticate", controllers.Authenticate).Methods("POST")
	openRouter.HandleFunc("/auth/refresh", controllers.Refresh).Methods("POST")

	an := negroni.New(negroni.HandlerFunc(services.JwtMiddleware), negroni.Wrap(authRouter))
	openRouter.PathPrefix("/").Handler(an)
	n := negroni.Classic()
	n.UseHandler(openRouter)
	n.Run(":"+settings.Application.Port)
}
