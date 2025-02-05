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
	UploadImage(userID int, URL string, description string) error
	GetUserImagesID(userID int) ([]string, error)
	GetImageDescription(imageURL string) (string, error)
	DeleteImage(imageID string) error
}
