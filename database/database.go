package database

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mhafidk/ngartos/config"
	"github.com/mhafidk/ngartos/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func Connect() {
	var pEnv, hostEnv, userEnv, nameEnv, passwordEnv string

	environment := config.Config("ENVIRONMENT")
	if environment == "DEV" {
		pEnv = "DB_PORT_DEV"
		hostEnv = "DB_HOST_DEV"
		userEnv = "DB_USER_DEV"
		nameEnv = "DB_NAME_DEV"
		passwordEnv = "DB_PASSWORD_DEV"
	} else {
		pEnv = "DB_PORT"
		hostEnv = "DB_HOST"
		userEnv = "DB_USER"
		nameEnv = "DB_NAME"
		passwordEnv = "DB_PASSWORD"
	}

	p := config.Config(pEnv)
	host := config.Config(hostEnv)
	user := config.Config(userEnv)
	name := config.Config(nameEnv)
	password := config.Config(passwordEnv)

	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		fmt.Println("Error parsing str to int")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", host, user, password, name, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running migrations")
	db.AutoMigrate(&model.User{}, &model.Topic{}, &model.Exercise{}, &model.Bookmark{})

	DB = Dbinstance{
		Db: db,
	}
}
