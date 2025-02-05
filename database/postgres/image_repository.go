package postgres

import (
	"gorm.io/gorm"
	"log"
	"pictureloader/models"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (i *ImageRepository) UploadImage(userID int, URL string, description string) error {
	imageModel := models.Image{URL: URL, UserID: userID, Description: description}
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

func (i *ImageRepository) GetImageDescription(imageURL string) (string, error) {
	var description string
	err := i.db.Model(&models.Image{}).Where("url = ?", imageURL).Pluck("description", &description).Error
	if err != nil {
		return "", err
	}
	return description, nil
}

func (i *ImageRepository) DeleteImage(imageID string) error {
	err := i.db.Where("url = ?", imageID).Delete(&models.Image{}).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
