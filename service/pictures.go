package service

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	db "pictureloader/database"
	"pictureloader/database/postgres"
	"pictureloader/image_storage"
	"pictureloader/models"
	"strconv"
)

type PictureLoader struct {
	storage  image_storage.ImageStorage
	database db.ImageRepositoryInterface
}

func NewPictureLoader(storage image_storage.ImageStorage, database *postgres.ImageRepository) *PictureLoader {
	return &PictureLoader{storage, database}
}

func (p *PictureLoader) Upload(ctx context.Context, img models.ImageUnit, userID int, description string) (string, error) {
	url := strconv.Itoa(rand.Int())
	imgName, err := p.storage.UploadFile(ctx, img, url)
	if err != nil {
		slog.Error("S3 Upload Error", "error", err)
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}
	err = p.database.UploadImage(userID, url, description)
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
	description, err := p.database.GetImageDescription(imgURL)
	if err != nil {
		slog.Error("Database download description error", "error", err)
		return "", "", fmt.Errorf("failed to get image description: %w", err)
	}
	return img, description, nil
}

func (p *PictureLoader) GetAllUserPictures(ctx context.Context, userID int) ([]string, error) {
	imageIDS, err := p.database.GetUserImagesID(userID)
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

func (p *PictureLoader) Delete(ctx context.Context, imgName string) error {
	err := p.storage.DeleteFileByURL(ctx, imgName)
	if err != nil {
		slog.Info("Storage delete error", "error", err)
		return err
	}
	err = p.database.DeleteImage(imgName)
	if err != nil {
		slog.Info("Database delete error", "error", err)
		return err
	}
	return nil
}
