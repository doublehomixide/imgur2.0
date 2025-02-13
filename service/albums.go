package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log/slog"
	cache "pictureloader/caching"
	db "pictureloader/database"
	"pictureloader/image_storage"
	"pictureloader/models"
	"strconv"
	"time"
)

type AlbumService struct {
	database      db.AlbumRepositoryInterface
	storage       image_storage.ImageStorage
	imageDatabase db.ImageRepositoryInterface
	cache         cache.Cacher
}

func NewAlbumService(database db.AlbumRepositoryInterface, storage image_storage.ImageStorage,
	imageDatabase db.ImageRepositoryInterface, cacher cache.Cacher) *AlbumService {
	return &AlbumService{
		database:      database,
		storage:       storage,
		imageDatabase: imageDatabase,
		cache:         cacher,
	}
}

func (als *AlbumService) CreateAlbum(ctx context.Context, album *models.Album) error {
	err := als.database.CreateAlbum(ctx, album)
	if err != nil {
		slog.Error("Create album", "error", err)
		return err
	}
	return nil
}

func (als *AlbumService) GetAlbum(ctx context.Context, albumID int) (models.AlbumUnit, error) {

	cachedResult, err := als.cache.Get(ctx, strconv.Itoa(albumID))
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Error("Get album", "error", err)
		return models.AlbumUnit{}, err
	}
	if cachedResult != "" {
		slog.Info("Get album from cache", "albumID", albumID)
		var album models.AlbumUnit
		json.Unmarshal([]byte(cachedResult), &album)
		return album, nil
	}

	albumName, images, err := als.database.GetAlbumData(ctx, albumID)
	if err != nil {
		slog.Error("Get images from album", "error", err)
		return models.AlbumUnit{}, err
	}

	updatedImages := make(map[string]string)

	for key, value := range images {
		presignedURL, err := als.storage.GetFileURL(ctx, key)
		if err != nil {
			continue
		}
		updatedImages[value] = presignedURL
	}

	result := models.AlbumUnit{
		Name:   albumName,
		Images: updatedImages,
	}

	resultJSON, _ := json.Marshal(result)

	err = als.cache.Set(ctx, strconv.Itoa(albumID), string(resultJSON), time.Hour*5)
	if err != nil {
		slog.Error("Set album", "error", err)
		return models.AlbumUnit{}, err
	}

	return result, nil
}

// GetUserAlbums userID -> hashmap where key is album ID, value is models.AlbumUnit
func (als *AlbumService) GetUserAlbums(ctx context.Context, userID int) (map[int]models.AlbumUnit, error) {
	cachedResult, err := als.cache.Get(ctx, strconv.Itoa(userID)+"_albums")
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Error("Get albums", "error", err)
		return map[int]models.AlbumUnit{}, err
	}
	if cachedResult != "" {
		slog.Info("Get user albums from cache", "userID", userID)
		var userAlbums map[int]models.AlbumUnit
		json.Unmarshal([]byte(cachedResult), &userAlbums)
		return userAlbums, nil
	}

	albumIDs, err := als.database.GetUserAlbumIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make(map[int]models.AlbumUnit)

	for _, albumID := range albumIDs {
		album, err := als.GetAlbum(ctx, albumID)
		if err != nil {
			continue
		}
		result[albumID] = album
	}

	resultJSON, _ := json.Marshal(result)

	err = als.cache.Set(ctx, strconv.Itoa(userID)+"_albums", string(resultJSON), time.Hour*5)
	if err != nil {
		slog.Error("Set album", "error", err)
		return map[int]models.AlbumUnit{}, err
	}

	return result, nil
}

func (als *AlbumService) AppendImageToAlbum(ctx context.Context, albumID int, imageSK string, userID int) error {
	err := als.database.IsOwnerOfAlbum(ctx, userID, albumID)
	if err != nil {
		slog.Error("Database delete error", "error", err)
		return err
	}

	imageID, err := als.imageDatabase.GetImageIDBySK(ctx, imageSK)
	if err != nil {
		slog.Info("Add image to album", "error", err)
		return err
	}

	albumImageModel := &models.AlbumImage{AlbumID: albumID, ImageID: imageID}

	err = als.database.CreateAlbumAndImage(ctx, albumImageModel)
	if err != nil {
		slog.Error("Add image to album", "error", err)
		return err
	}

	err = als.cache.Delete(ctx, strconv.Itoa(albumID))
	if err != nil {
		slog.Error("Delete album", "error", err)
		return err
	}
	err = als.cache.Delete(ctx, strconv.Itoa(userID)+"_albums")
	if err != nil {
		slog.Error("Delete user albums", "error", err)
		return err
	}

	return nil
}

func (als *AlbumService) DeleteAlbum(ctx context.Context, albumID int, userID int) error {
	err := als.database.IsOwnerOfAlbum(ctx, userID, albumID)
	if err != nil {
		slog.Error("Database delete error", "error", err)
		return err
	}

	err = als.database.DeleteAlbumByID(ctx, albumID)
	if err != nil {
		slog.Error("Delete album", "error", err)
		return err
	}

	err = als.cache.Delete(ctx, strconv.Itoa(albumID))
	if err != nil {
		slog.Error("Delete album", "error", err)
		return err
	}
	err = als.cache.Delete(ctx, strconv.Itoa(userID)+"_albums")
	if err != nil {
		slog.Error("Delete user albums", "error", err)
		return err
	}

	return nil
}

func (als *AlbumService) DeleteImageFromAlbum(ctx context.Context, albumID int, imageSK string, userID int) error {
	err := als.database.IsOwnerOfAlbum(ctx, userID, albumID)
	if err != nil {
		slog.Error("Database delete error", "error", err)
		return err
	}

	imageID, err := als.imageDatabase.GetImageIDBySK(ctx, imageSK)
	if err != nil {
		slog.Info("Delete image from album", "error", err)
		return err
	}
	err = als.database.DeleteAlbumImage(ctx, albumID, imageID)
	if err != nil {
		slog.Error("Delete image from album", "error", err)
		return err
	}

	err = als.cache.Delete(ctx, strconv.Itoa(albumID))
	if err != nil {
		slog.Error("Delete album", "error", err)
		return err
	}
	err = als.cache.Delete(ctx, strconv.Itoa(userID)+"_albums")
	if err != nil {
		slog.Error("Delete user albums", "error", err)
		return err
	}

	return nil
}
