package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"pictureloader/app_microservice/models"
)

func NewDataBase(dbPath string) *gorm.DB {
	database, err := gorm.Open(postgres.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	err = database.AutoMigrate(&models.User{}, &models.Image{}, &models.Post{}, &models.PostImage{}, &models.Like{})
	if err != nil {
		log.Fatalln(err)
	}
	return database
}
