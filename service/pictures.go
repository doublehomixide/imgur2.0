package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	db "pictureloader/database"
	"pictureloader/database/postgres"
	"pictureloader/image_storage"
	"pictureloader/models"
	"strconv"
)

type PictureLoader struct {
	storage  image_storage.ImageStorage
	database db.ImageRepository
}

func NewPictureLoader(storage image_storage.ImageStorage, database *postgres.ImageRepository) *PictureLoader {
	return &PictureLoader{storage, database}
}

func (p *PictureLoader) Upload(ctx context.Context, img models.ImageUnit, userID int, description string) (string, error) {
	url := strconv.Itoa(rand.Int())
	imgName, err := p.storage.UploadFile(ctx, img, url)
	if err != nil {
		log.Fatal(err)
	}
	err = p.database.UploadImage(userID, url, description)
	if err != nil {
		log.Fatal(err)
	}
	return imgName, nil
}

func (p *PictureLoader) Download(ctx context.Context, imgURL string) (string, string, error) {
	img, err := p.storage.GetFileURL(ctx, imgURL)
	if err != nil {
		log.Fatalf("Error downloading file: %v", err)
	}
	description, err := p.database.GetImageDescription(imgURL)
	if err != nil {
		dscErr := fmt.Sprintf("Error getting description: %v", err)
		return img, dscErr, err
	}
	return img, description, nil
}

func (p *PictureLoader) GetAllUserPictures(ctx context.Context, userID int) ([]string, error) {
	imageIDS, err := p.database.GetUserImagesID(userID)
	if err != nil {
		return nil, err
	}
	imageURLS, err := p.storage.GetFileURLS(ctx, imageIDS)
	if err != nil {
		return nil, err
	}
	return imageURLS, err
}

func (p *PictureLoader) Delete(ctx context.Context, imgName string) error {
	err := p.storage.DeleteFileByURL(ctx, imgName)
	if err != nil {
		return err
	}
	err = p.database.DeleteImage(imgName)
	if err != nil {
		return err
	}
	return nil
}
