package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/deploji/deploji-server/controllers"
	"github.com/deploji/deploji-server/middleware"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/services/amqpService"
	"github.com/deploji/deploji-server/services/auth"
	"github.com/deploji/deploji-server/settings"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"golang.org/x/net/context"
	"log"
	"os"
	"time"
)

func main() {
	ctx, done := context.WithCancel(context.Background())
	models.InitDatabase()
	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		settings.Database.Host,
		settings.Database.Port,
		settings.Database.Username,
		settings.Database.Name,
		settings.Database.Password)
	a, err := gormadapter.NewAdapter("postgres", uri, true)
	if err != nil {
		log.Printf("NewAdapter error: %s", err)
	}
	auth.E, err = casbin.NewSyncedEnforcer("rbac_model.conf", a)
	if err != nil {
		log.Printf("NewSyncedEnforcer error: %s", err)
		os.Exit(1)
	}
	auth.E.StartAutoLoadPolicy(time.Second * 60)

	go func() {
		amqpService.Publish(amqpService.Redial(ctx, os.Getenv("AMQP_URL")), amqpService.Jobs, "jobs")
		done()
	}()

	authRouter := mux.NewRouter()
	authRouter.HandleFunc("/ssh-keys", controllers.GetSshKeys).Methods("GET")
	authRouter.HandleFunc("/ssh-keys", controllers.SaveSshKeys).Methods("POST")
	authRouter.HandleFunc("/ssh-keys/{id}", controllers.SaveSshKeys).Methods("PUT")
	authRouter.HandleFunc("/ssh-keys/{id}", controllers.GetSshKey).Methods("GET")
	authRouter.HandleFunc("/ssh-keys/{id}", controllers.DeleteSshKeys).Methods("DELETE")
	authRouter.HandleFunc("/inventories", controllers.GetInventories).Methods("GET")
	authRouter.HandleFunc("/inventories", controllers.SaveInventory).Methods("POST")
	authRouter.HandleFunc("/inventories/{id}", controllers.SaveInventory).Methods("PUT")
	authRouter.HandleFunc("/inventories/{id}", controllers.GetInventory).Methods("GET")
	authRouter.HandleFunc("/inventories/{id}", controllers.DeleteInventory).Methods("DELETE")
	authRouter.HandleFunc("/templates", controllers.GetTemplates).Methods("GET")
	authRouter.HandleFunc("/templates", controllers.SaveTemplate).Methods("POST")
	authRouter.HandleFunc("/templates/{id}", controllers.SaveTemplate).Methods("PUT")
	authRouter.HandleFunc("/templates/{id}", controllers.GetTemplate).Methods("GET")
	authRouter.HandleFunc("/templates/{id}", controllers.DeleteTemplate).Methods("DELETE")
	authRouter.HandleFunc("/applications", controllers.GetApplications).Methods("GET")
	authRouter.HandleFunc("/applications", controllers.SaveApplications).Methods("POST")
	authRouter.HandleFunc("/applications/{id}", controllers.SaveApplications).Methods("PUT")
	authRouter.HandleFunc("/applications/{id}", controllers.GetApplication).Methods("GET")
	authRouter.HandleFunc("/applications/{id}/inventories", controllers.GetApplicationInventories).Methods("GET")
	authRouter.HandleFunc("/applications/{id}/notifications", controllers.GetApplicationNotifications).Methods("GET")
	authRouter.HandleFunc("/applications/{id}/notifications", controllers.SaveApplicationNotification).Methods("PUT")
	authRouter.HandleFunc("/applications/{id}", controllers.DeleteApplication).Methods("DELETE")
	authRouter.HandleFunc("/application-inventories/{id}", controllers.DeleteApplicationInventory).Methods("DELETE")
	authRouter.HandleFunc("/versions", controllers.GetVersions).Queries("app", "{app}").Methods("GET")
	authRouter.HandleFunc("/users", controllers.GetUsers).Methods("GET")
	authRouter.HandleFunc("/teams", controllers.GetTeams).Methods("GET")
	authRouter.HandleFunc("/permissions", controllers.GetPermissions).Methods("GET")
	authRouter.HandleFunc("/permissions", controllers.SavePermission).Methods("POST")
	authRouter.HandleFunc("/permissions/delete", controllers.DeletePermission).Methods("POST")
	authRouter.HandleFunc("/jobs", controllers.GetJobs).Methods("GET")
	authRouter.HandleFunc("/jobs", controllers.SaveJobs).Methods("POST")
	authRouter.HandleFunc("/jobs/latest-deployments", controllers.GetLatestDeployments).Methods("GET")
	authRouter.HandleFunc("/jobs/{id}", controllers.GetJob).Methods("GET")
	authRouter.HandleFunc("/jobs/{id}/logs", controllers.GetJobLogs).Methods("GET")
	authRouter.HandleFunc("/notification-channels", controllers.GetNotificationChannels).Methods("GET")

	adminRouter := mux.NewRouter()
	adminRouter.HandleFunc("/teams", controllers.SaveTeam).Methods("POST")
	adminRouter.HandleFunc("/teams/{id}", controllers.GetTeam).Methods("GET")
	adminRouter.HandleFunc("/teams/{id}/users", controllers.GetTeamUsers).Methods("GET")
	adminRouter.HandleFunc("/teams/{id}", controllers.SaveTeam).Methods("PUT")
	adminRouter.HandleFunc("/teams/{id}", controllers.DeleteTeam).Methods("DELETE")
	adminRouter.HandleFunc("/teams/{id}/users", controllers.SaveTeamUser).Methods("POST")
	adminRouter.HandleFunc("/teams/{id}/users/{userId}", controllers.DeleteTeamUser).Methods("DELETE")
	adminRouter.HandleFunc("/teams/{id}/permissions", controllers.GetTeamPermissions).Methods("GET")
	adminRouter.HandleFunc("/users", controllers.SaveUser).Methods("POST")
	adminRouter.HandleFunc("/users/{id}", controllers.SaveUser).Methods("PUT")
	adminRouter.HandleFunc("/users/{id}", controllers.GetUser).Methods("GET")
	adminRouter.HandleFunc("/users/{id}/permissions", controllers.GetUserPermissions).Methods("GET")
	adminRouter.HandleFunc("/settings", controllers.GetSettings).Methods("GET")
	adminRouter.HandleFunc("/settings", controllers.SaveSettings).Methods("PUT")
	adminRouter.HandleFunc("/repositories", controllers.GetRepositories).Methods("GET")
	adminRouter.HandleFunc("/repositories", controllers.SaveRepositories).Methods("POST")
	adminRouter.HandleFunc("/repositories/{id}", controllers.SaveRepositories).Methods("PUT")
	adminRouter.HandleFunc("/repositories/{id}", controllers.GetRepository).Methods("GET")
	adminRouter.HandleFunc("/repositories/{id}", controllers.DeleteRepository).Methods("DELETE")
	adminRouter.HandleFunc("/projects", controllers.GetProjects).Methods("GET")
	adminRouter.HandleFunc("/projects/synchronize-status", controllers.GetProjectsSyncStatus).Methods("GET")
	adminRouter.HandleFunc("/projects", controllers.SaveProjects).Methods("POST")
	adminRouter.HandleFunc("/projects/{id}", controllers.SaveProjects).Methods("PUT")
	adminRouter.HandleFunc("/projects/{id}", controllers.GetProject).Methods("GET")
	adminRouter.HandleFunc("/projects/{id}", controllers.DeleteProject).Methods("DELETE")
	adminRouter.HandleFunc("/projects/{id}/synchronize", controllers.SynchronizeProject).Methods("POST")
	adminRouter.HandleFunc("/projects/{id}/files", controllers.GetProjectFiles).Methods("GET")
	adminRouter.HandleFunc("/notification-channels", controllers.SaveNotificationChannels).Methods("POST")
	adminRouter.HandleFunc("/notification-channels/{id}", controllers.SaveNotificationChannels).Methods("PUT")
	adminRouter.HandleFunc("/notification-channels/{id}", controllers.GetNotificationChannel).Methods("GET")
	adminRouter.HandleFunc("/notification-channels/{id}", controllers.DeleteNotificationChannel).Methods("DELETE")

	openRouter := mux.NewRouter()
	openRouter.HandleFunc("/auth/authenticate", controllers.Authenticate).Methods("POST")
	openRouter.HandleFunc("/auth/refresh", controllers.Refresh).Methods("POST")

	authNegroni := negroni.New(
		negroni.HandlerFunc(middleware.JwtMiddleware),
		negroni.HandlerFunc(middleware.AuthMiddleware),
		negroni.HandlerFunc(middleware.HeadersMiddleware),
		negroni.Wrap(authRouter))

	adminNegroni := negroni.New(
		negroni.HandlerFunc(middleware.AdminOnlyMiddleware),
		negroni.Wrap(adminRouter))

	authRouter.PathPrefix("/").Handler(adminNegroni)
	openRouter.PathPrefix("/").Handler(authNegroni)
	n := negroni.Classic()
	n.UseHandler(openRouter)
	n.Run(":" + settings.Application.Port)
}
