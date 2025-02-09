package service

import (
	"context"
	"log/slog"
	db "pictureloader/database"
	"pictureloader/image_storage"
	"pictureloader/models"
)

type AlbumService struct {
	database      db.AlbumRepositoryInterface
	storage       image_storage.ImageStorage
	imageDatabase db.ImageRepositoryInterface
}

func NewAlbumService(database db.AlbumRepositoryInterface, storage image_storage.ImageStorage, imageDatabase db.ImageRepositoryInterface) *AlbumService {
	return &AlbumService{
		database:      database,
		storage:       storage,
		imageDatabase: imageDatabase,
	}
}

func (als *AlbumService) CreateAlbum(album *models.Album) error {
	err := als.database.CreateAlbum(album)
	if err != nil {
		slog.Error("Create album", "error", err)
		return err
	}
	return nil
}

func (als *AlbumService) GetAlbum(albumID int) (models.AlbumUnit, error) {
	albumName, images, err := als.database.GetAlbumData(albumID)
	if err != nil {
		slog.Error("Get images from album", "error", err)
		return models.AlbumUnit{}, err
	}

	updatedImages := make(map[string]string)

	for key, value := range images {
		presignedURL, err := als.storage.GetFileURL(context.Background(), key)
		if err != nil {
			continue
		}
		updatedImages[value] = presignedURL
	}

	result := models.AlbumUnit{
		Name:   albumName,
		Images: updatedImages,
	}

	return result, nil
}

// GetUserAlbums userID -> hashmap where key is album ID, value is models.AlbumUnit
func (als *AlbumService) GetUserAlbums(userID int) (map[int]models.AlbumUnit, error) {
	albumIDs, err := als.database.GetUserAlbumIDs(userID)
	if err != nil {
		return nil, err
	}

	result := make(map[int]models.AlbumUnit)

	for _, albumID := range albumIDs {
		album, err := als.GetAlbum(albumID)
		if err != nil {
			continue
		}
		result[albumID] = album
	}
	return result, nil
}

func (als *AlbumService) AppendImageToAlbum(albumID int, imageSK string) error {

	imageID, err := als.imageDatabase.GetImageIDBySK(imageSK)
	if err != nil {
		slog.Info("Add image to album", "error", err)
		return err
	}

	albumImageModel := &models.AlbumImage{AlbumID: albumID, ImageID: imageID}

	err = als.database.CreateAlbumAndImage(albumImageModel)
	if err != nil {
		slog.Error("Add image to album", "error", err)
		return err
	}
	return nil
}

func (als *AlbumService) DeleteAlbum(albumID int) error {
	err := als.database.DeleteAlbumByID(albumID)
	if err != nil {
		slog.Error("Delete album", "error", err)
		return err
	}
	return nil
}
