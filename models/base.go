package models

import (
	"fmt"
	"github.com/deploji/deploji-server/settings"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
)

var db *gorm.DB

func InitDatabase() {
	settings.Load()
	conn, err := gorm.Open(settings.Database.Type, settings.Database.URI)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.LogMode(os.Getenv("GORM_LOG_MODE") == "true")
	db.AutoMigrate(
		&Project{},
		&SshKey{},
		&Application{},
		&ApplicationInventory{},
		&Inventory{},
		&Job{},
		&JobLog{},
		&Repository{},
		&Template{},
		&Team{},
		&User{},
		&Setting{},
		&NotificationChannel{},
		&ApplicationNotification{},
		&TemplateNotification{},
		&ProjectNotification{},
		&PushSubscription{},
		&Survey{},
		&SurveyInput{},
		&SettingGroup{})

	driver, err := postgres.WithInstance(db.DB(), &postgres.Config{})
	runMigrations(driver)
}

func runMigrations(driver database.Driver) error {
	fsrc, err := (&file.File{}).Open("file://migrations")
	if err != nil {
		log.Printf("Cannot open migrations file: %s", err)
		return err
	}
	m, err := migrate.NewWithInstance(
		"file",
		fsrc,
		"postgres",
		driver)
	if err != nil {
		log.Printf("Cannot create migrate instance: %s", err)
		return err
	}
	if err := m.Steps(2); err != nil {
		log.Printf("Migration error: %s", err)
		return err
	}
	return nil
}

func GetDB() *gorm.DB {
	return db.Set("gorm:association_autoupdate", false)
}
