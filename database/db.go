package db

import (
	"pictureloader/models"
)

type UserRepositoryInterface interface {
	CreateNewUser(user *models.User) error
	GetUserByID(ID int) (*models.User, error)
	DeleteUserByID(ID int) error
	GetUserByUsername(username string) (*models.User, error)
}

type ImageRepositoryInterface interface {
	UploadImage(userID int, URL string, description string) error
	GetUserImagesID(userID int) ([]string, error)
	GetImageDescription(imageURL string) (string, error)
	DeleteImage(imageID string) error
	GetImageIDBySK(imageSK string) (int, error)
}

type AlbumRepositoryInterface interface {
	CreateAlbum(album *models.Album) error
	CreateAlbumAndImage(albumImage *models.AlbumImage) error
	GetAlbumData(albumID int) (string, map[string]string, error)
	GetUserAlbumIDs(userID int) ([]int, error)
	DeleteAlbumByID(albumID int) error
	DeleteAlbumImage(albumID int, imageID int) error
}
