package service

import (
	"context"
	"log"
	"math/rand"
	"pictureloader/database/postgres"
	"pictureloader/image_storage"
	"pictureloader/models"
	"strconv"
)

type PictureLoader struct {
	storage  image_storage.ImageStorage
	database *postgres.ImageRepository
}

func NewPictureLoader(storage image_storage.ImageStorage, database *postgres.ImageRepository) *PictureLoader {
	return &PictureLoader{storage, database}
}

func (p *PictureLoader) Upload(ctx context.Context, img models.ImageUnit, userID int) (string, error) {
	url := strconv.Itoa(rand.Int())
	imgName, err := p.storage.UploadFile(ctx, img, url)
	if err != nil {
		log.Fatal(err)
	}
	err = p.database.UploadImage(userID, url)
	if err != nil {
		log.Fatal(err)
	}
	return imgName, nil
}

func (p *PictureLoader) Download(ctx context.Context, imgName string) (string, error) {
	img, err := p.storage.GetFileURL(ctx, imgName)
	if err != nil {
		log.Fatalf("Error downloading file: %v", err)
	}
	return img, nil
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
