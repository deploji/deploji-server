package settings

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

var Application ApplicationSettings
var Auth AuthSettings
var Database DatabaseSettings
var Amqp AmqpSettings

type ApplicationSettings struct {
	Port string
}

type AuthSettings struct {
	JWTSecret  string
	RefreshTTL time.Duration
	TTL        time.Duration
}

type AmqpSettings struct {
	URL string
}

type DatabaseSettings struct {
	Username string
	Password string
	Name     string
	Host     string
	Port     string
	Type     string
	URI      string
}

func Load() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	Application.Port = os.Getenv("PORT")
	Auth.JWTSecret = os.Getenv("JWT_SECRET")
	ttl, _ := strconv.ParseInt(os.Getenv("JWT_TTL"), 10, 0)
	refreshTtl, _ := strconv.ParseInt(os.Getenv("JWT_REFRESH_TTL"), 10, 0)
	Auth.RefreshTTL = time.Duration(refreshTtl) * time.Hour
	Auth.TTL = time.Duration(ttl) * time.Minute
	Amqp.URL = os.Getenv("AMQP_URL")
	Database.Username = os.Getenv("DB_USER")
	Database.Password = os.Getenv("DB_PASS")
	Database.Name = os.Getenv("DB_NAME")
	Database.Host = os.Getenv("DB_HOST")
	Database.Port = os.Getenv("DB_PORT")
	Database.Type = os.Getenv("DB_TYPE")

	if Database.Type == "mysql" {
		Database.URI = fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			Database.Username,
			Database.Password,
			Database.Host,
			Database.Name)
	} else if Database.Type == "postgres" {
		Database.URI = fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			Database.Host,
			Database.Port,
			Database.Username,
			Database.Name,
			Database.Password)
	} else {
		log.Printf("unsupported gorm database type: %s", Database.Type)
		os.Exit(1)
	}
}
