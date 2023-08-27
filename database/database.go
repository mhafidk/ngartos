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
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		fmt.Println("Error parsing str to int")
	}

	host := config.Config("DB_HOST")
	user := config.Config("DB_USER")
	name := config.Config("DB_NAME")
	password := config.Config("DB_PASSWORD")
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