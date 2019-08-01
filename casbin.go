package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v2"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	// Initialize a Gorm adapter and use it in a Casbin enforcer:
	// The adapter will use the MySQL database named "casbin".
	// If it doesn't exist, the adapter will create it automatically.
	// You can also use an already existing gorm instance with gormadapter.NewAdapterByDB(gormInstance)
	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", "localhost", "5432", "mastermind", "mastermind", "mastermind")
	a, err := gormadapter.NewAdapter("postgres", uri, true) // Your driver and data source.
	if err != nil {
		log.Printf("NewAdapter error: %s", err)
	}
	e, _ := casbin.NewEnforcer("rbac_model.conf", a)

	// Or you can use an existing DB "abc" like this:
	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.
	// a := gormadapter.NewAdapter("mysql", "mysql_username:mysql_password@tcp(127.0.0.1:3306)/abc", true)

	// Load the policy from DB.
	if err := e.LoadPolicy(); err != nil {
		log.Printf("LoadPolicy error: %s", err)
	}

	// Check the permission.
	hasPerm, err := e.Enforce("alice", "data1", "read")
	if err != nil {
		log.Printf("NewAdapter error: %s", err)
	}
	log.Printf("Perm: %t", hasPerm)
	// Modify the policy.
	if _, err := e.AddRoleForUser("alice", "data1"); err != nil {
		log.Printf("AddPolicy error: %s", err)
	}
	if _, err := e.AddGroupingPolicy("group1", "data1_admin"); err != nil {
		log.Printf("AddPolicy error: %s", err)
	}
	if _, err := e.AddGroupingPolicy("alice", "group1"); err != nil {
		log.Printf("AddPolicy error: %s", err)
	}
	if _, err := e.AddPolicy("group1", "data1", "read"); err != nil {
		log.Printf("AddPolicy error: %s", err)
	}
	p, _ := e.GetRolesForUser("alice")
	// e.RemovePolicy(...)
	allObjects := e.GetAllObjects()
	allSubjects := e.GetAllSubjects()
	allRoles := e.GetAllRoles()
	allActions := e.GetAllActions()
	g := e.GetGroupingPolicy()
	log.Printf("Objects: %#v", allObjects)
	log.Printf("Subjects: %#v", allSubjects)
	log.Printf("Roles: %#v", allRoles)
	log.Printf("Actions: %#v", allActions)
	log.Printf("P: %#v", p)
	log.Printf("g: %#v", g)

	// Save the policy back to DB.
	if err := e.SavePolicy(); err != nil {
		log.Printf("SavePolicy error: %s", err)
	}
}
