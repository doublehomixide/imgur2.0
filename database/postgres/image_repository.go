package postgres

import (
	"context"
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

func (i *ImageRepository) UploadImage(ctx context.Context, userID int, URL string, description string) error {
	imageModel := models.Image{StorageKey: URL, UserID: userID, Description: description}
	return i.db.WithContext(ctx).Create(&imageModel).Error
}

func (i *ImageRepository) GetUserImagesID(ctx context.Context, userID int) ([]string, error) {
	var imageIDS []string
	err := i.db.WithContext(ctx).Model(&models.Image{}).Where("user_id = ?", userID).Pluck("storage_key", &imageIDS).Error
	if err != nil {
		return nil, err
	}
	return imageIDS, nil
}

func (i *ImageRepository) GetImageDescription(ctx context.Context, imageURL string) (string, error) {
	var description string
	err := i.db.WithContext(ctx).Model(&models.Image{}).Where("storage_key = ?", imageURL).Pluck("description", &description).Error
	if err != nil {
		return "", err
	}
	return description, nil
}

func (i *ImageRepository) DeleteImage(ctx context.Context, imageID string) error {
	err := i.db.WithContext(ctx).Where("storage_key  = ?", imageID).Delete(&models.Image{}).Error
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (i *ImageRepository) GetImageIDBySK(ctx context.Context, imageSK string) (int, error) {
	var imageID int
	i.db.WithContext(ctx).Model(&models.Image{}).Where("storage_key  = ?", imageSK).Pluck("id", &imageID)
	if imageID == 0 {
		return 0, fmt.Errorf("no such image with SK %d", imageSK)
	}
	return imageID, nil
}

func (i *ImageRepository) IsOwnerOfPicture(ctx context.Context, userID int, imageSK string) error {
	var trueUserID int
	i.db.WithContext(ctx).Model(&models.Image{}).Where("storage_key = ?", imageSK).Pluck("user_id", &trueUserID)
	if userID != trueUserID {
		return fmt.Errorf("user is not owner of this picture %s", imageSK)
	}
	return nil
}

func (i *ImageRepository) GetImageLinkedPost(ctx context.Context, imageSK string) (int, error) {
	var postID int
	err := i.db.WithContext(ctx).
		Table("post_images").
		Select("post_id").
		Joins("JOIN images ON post_images.image_id = images.id").
		Where("images.storage_key = ?", imageSK).
		Limit(1).
		Pluck("post_id", &postID).Error
	if err != nil {
		return 0, err
	}
	return postID, nil
}
