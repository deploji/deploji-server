package models

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sotomskir/mastermind-server/settings"
	"log"
	"os"
)

var db *gorm.DB

func init() {
	settings.Load()
	conn, err := gorm.Open(settings.Database.Type, settings.Database.URI)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.LogMode(false)
	db.AutoMigrate(
		&Project{},
		&SshKey{},
		&Application{},
		&ApplicationInventory{},
		&Inventory{},
		&Repository{},
		&DeploymentLog{},
		&Template{},
		&User{},
		&Setting{},
		&SettingGroup{},
		&Deployment{})

	driver, err := postgres.WithInstance(db.DB(), &postgres.Config{})
	fsrc, err := (&file.File{}).Open("file://migrations")
	if err != nil {
		log.Printf("Cannot open migrations file: %s", err)
		os.Exit(1)
	}
	m, err := migrate.NewWithInstance(
		"file",
		fsrc,
		"postgres",
		driver)
	if err != nil {
		log.Printf("Cannot create migrate instance: %s", err)
		os.Exit(1)
	}
	if err := m.Steps(2); err != nil {
		log.Printf("Migration error: %s", err)
	}
}

func GetDB() *gorm.DB {
	return db.Set("gorm:association_autoupdate", false)
}
