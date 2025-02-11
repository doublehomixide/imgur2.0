package postgres

import (
	"fmt"
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
	imageModel := models.Image{StorageKey: URL, UserID: userID, Description: description}
	return i.db.Create(&imageModel).Error
}

func (i *ImageRepository) GetUserImagesID(userID int) ([]string, error) {
	var imageIDS []string
	err := i.db.Model(&models.Image{}).Where("user_id = ?", userID).Pluck("storage_key", &imageIDS).Error
	if err != nil {
		return nil, err
	}
	return imageIDS, nil
}

func (i *ImageRepository) GetImageDescription(imageURL string) (string, error) {
	var description string
	err := i.db.Model(&models.Image{}).Where("storage_key = ?", imageURL).Pluck("description", &description).Error
	if err != nil {
		return "", err
	}
	return description, nil
}

func (i *ImageRepository) DeleteImage(imageID string) error {
	err := i.db.Where("storage_key  = ?", imageID).Delete(&models.Image{}).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (i *ImageRepository) GetImageIDBySK(imageSK string) (int, error) {
	var imageID int
	i.db.Model(&models.Image{}).Where("storage_key  = ?", imageSK).Pluck("id", &imageID)
	if imageID == 0 {
		return 0, fmt.Errorf("no such image with SK %d", imageSK)
	}
	return imageID, nil
}

// IsOwnerOfPicture userID, imageSK -> true if userID == imageSK.userID else false
func (i *ImageRepository) IsOwnerOfPicture(userID int, imageSK string) error {
	var trueUserID int
	i.db.Model(&models.Image{}).Where("storage_key = ?", imageSK).Pluck("user_id", &trueUserID)
	if userID != trueUserID {
		return fmt.Errorf("user is not owner of this picture %s", imageSK)
	}
	return nil
}
