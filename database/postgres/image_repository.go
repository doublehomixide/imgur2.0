package postgres

import (
	"gorm.io/gorm"
	"pictureloader/models"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (i *ImageRepository) UploadImage(userID int, URL string) error {
	imageModel := models.Image{URL: URL, UserID: userID}
	return i.db.Create(imageModel).Error
}

func (i *ImageRepository) GetUserImagesID(userID int) ([]string, error) {
	var imageIDS []string
	err := i.db.Model(&models.Image{}).Where("user_id = ?", userID).Pluck("url", &imageIDS).Error
	if err != nil {
		return nil, err
	}
	return imageIDS, nil
}
