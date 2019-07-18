package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"os"
)

var db *gorm.DB

func init() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbType := os.Getenv("DB_TYPE")

	var dbUri string
	if dbType == "mysql" {
		dbUri = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbHost, dbName)
	} else if dbType == "postgres" {
		dbUri = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", dbHost, dbPort, username, dbName, password)
	} else {
		panic("unsupported database type")
	}

	fmt.Println(dbUri)

	conn, err := gorm.Open(dbType, dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.LogMode(true)
	db.AutoMigrate(&Project{}, &SshKey{})

}

func GetDB() *gorm.DB {
	return db
}
