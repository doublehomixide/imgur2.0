package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"pictureloader/database/postgres"
	"pictureloader/image_storage"
	"pictureloader/models"
	"strconv"
)

type ImageManager interface {
	UploadImage(ctx context.Context, userID int, URL string, description string) error
	GetUserImagesID(ctx context.Context, userID int) ([]string, error)
	GetImageDescription(ctx context.Context, imageURL string) (string, error)
	DeleteImage(ctx context.Context, imageID string) error
	IsOwnerOfPicture(ctx context.Context, userID int, imageSK string) error
}

type Cacher interface {
	DeleteImageFromPostCache(ctx context.Context, imageSK string) (bool, error)
}

type PictureLoader struct {
	storage  image_storage.ImageStorage
	database ImageManager
	cache    Cacher
}

func NewPictureLoader(storage image_storage.ImageStorage, database *postgres.ImageRepository, cache Cacher) *PictureLoader {
	return &PictureLoader{storage, database, cache}
}

func (p *PictureLoader) Upload(ctx context.Context, img models.ImageUnit, userID int, description string) (string, error) {
	url := strconv.Itoa(rand.Int())
	imgName, err := p.storage.UploadFile(ctx, img, url)
	if err != nil {
		slog.Error("S3 Upload Error", "error", err)
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}
	err = p.database.UploadImage(ctx, userID, url, description)
	if err != nil {
		slog.Error("Database upload error", "error", err)
		return "", fmt.Errorf("failed to upload image to database: %w", err)
	}
	return imgName, nil
}

func (p *PictureLoader) Download(ctx context.Context, imgURL string) (string, string, error) {
	img, err := p.storage.GetFileURL(ctx, imgURL)
	if err != nil {
		slog.Error("S3 error downloading file", "error", err)
		return "", "", fmt.Errorf("failed to get file StorageKey from S3: %w", err)
	}
	description, err := p.database.GetImageDescription(ctx, imgURL)
	if err != nil {
		slog.Error("Database download description error", "error", err)
		return "", "", fmt.Errorf("failed to get image description: %w", err)
	}
	return img, description, nil
}

func (p *PictureLoader) GetAllUserPictures(ctx context.Context, userID int) ([]string, error) {
	imageIDS, err := p.database.GetUserImagesID(ctx, userID)
	if err != nil {
		slog.Error("Database get user images id error", "error", err)
		return nil, err
	}
	imageURLS, err := p.storage.GetFileURLS(ctx, imageIDS)
	if err != nil {
		slog.Error("Storage get file urls error", "error", err)
		return nil, err
	}
	return imageURLS, err
}

func (p *PictureLoader) Delete(ctx context.Context, userID int, imgSK string) error {
	err := p.database.IsOwnerOfPicture(ctx, userID, imgSK)
	if err != nil {
		slog.Error("Database delete error", "error", err)
		return err
	}

	result, err := p.cache.DeleteImageFromPostCache(ctx, imgSK)
	if err != nil {
		slog.Error("Cache delete error", "error", err)
		return err
	}
	if !result {
		slog.Error("Cache delete error: no such key")
		return errors.New("no such key")
	}

	err = p.storage.DeleteFileByURL(ctx, imgSK)
	if err != nil {
		slog.Info("Storage delete error", "error", err)
		return err
	}
	err = p.database.DeleteImage(ctx, imgSK)
	if err != nil {
		slog.Info("Database delete error", "error", err)
		return err
	}
	return nil
}
