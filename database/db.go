package db

import (
	"context"
	"pictureloader/models"
)

type UserRepositoryInterface interface {
	CreateNewUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	DeleteUserByID(ctx context.Context, id int) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	ChangeUsernameByID(ctx context.Context, userID int, newUsername string) error
}

type ImageRepositoryInterface interface {
	UploadImage(ctx context.Context, userID int, URL string, description string) error
	GetUserImagesID(ctx context.Context, userID int) ([]string, error)
	GetImageDescription(ctx context.Context, imageURL string) (string, error)
	DeleteImage(ctx context.Context, imageID string) error
	GetImageIDBySK(ctx context.Context, imageSK string) (int, error)
	IsOwnerOfPicture(ctx context.Context, userID int, imageSK string) error
}

type AlbumRepositoryInterface interface {
	CreateAlbum(ctx context.Context, album *models.Album) error
	CreateAlbumAndImage(ctx context.Context, albumImage *models.AlbumImage) error
	GetAlbumData(ctx context.Context, albumID int) (string, map[string]string, error)
	GetUserAlbumIDs(ctx context.Context, userID int) ([]int, error)
	DeleteAlbumByID(ctx context.Context, albumID int) error
	DeleteAlbumImage(ctx context.Context, albumID int, imageID int) error
	IsOwnerOfAlbum(ctx context.Context, userID int, albumID int) error
}
