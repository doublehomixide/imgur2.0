package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewPSQL(dbPath string) *gorm.DB {
	database, err := gorm.Open(postgres.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	err = database.AutoMigrate(&LikesNotification{})
	if err != nil {
		log.Fatalln(err)
	}
	return database
}
