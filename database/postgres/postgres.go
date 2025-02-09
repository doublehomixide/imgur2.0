package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"pictureloader/models"
)

func NewDataBase(dbPath string) *gorm.DB {
	database, err := gorm.Open(postgres.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	err = database.AutoMigrate(&models.User{}, &models.Image{}, &models.Album{}, &models.AlbumImage{})
	if err != nil {
		log.Fatalln(err)
	}
	return database
}
