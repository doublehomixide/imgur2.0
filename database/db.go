package db

import (
	"pictureloader/models"
)

type UserRepository interface {
	CreateNewUser(user *models.User) error
	GetUserByID(ID int) (*models.User, error)
	DeleteUserByID(ID int) error
	GetUserByUsername(username string) (*models.User, error)
}

type ImageRepository interface {
	UploadImage(userID int, URL string) error
	GetUserImagesID(userID int) ([]models.Image, error)
}
