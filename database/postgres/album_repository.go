package postgres

import (
	"fmt"
	"gorm.io/gorm"
	"pictureloader/models"
)

type AlbumRepository struct {
	db *gorm.DB
}

func NewAlbumRepository(db *gorm.DB) *AlbumRepository {
	return &AlbumRepository{db: db}
}

func (ar *AlbumRepository) CreateAlbum(album *models.Album) error {
	return ar.db.Create(&album).Error
}

func (ar *AlbumRepository) CreateAlbumAndImage(albumImage *models.AlbumImage) error {
	return ar.db.Create(&albumImage).Error
}

// GetAlbumData returns album name, hashmap where key is s3 storage StorageKey, value is description, error
func (ar *AlbumRepository) GetAlbumData(albumID int) (string, map[string]string, error) {
	var album models.Album
	if err := ar.db.Preload("Images").First(&album, albumID).Error; err != nil {
		return "", nil, err
	}

	images := make(map[string]string)
	for _, image := range album.Images {
		images[image.StorageKey] = image.Description
	}

	return album.Name, images, nil
}

func (ar *AlbumRepository) GetUserAlbumIDs(userID int) ([]int, error) {
	var idSlice []int
	if err := ar.db.Model(&models.Album{}).Where("user_id = ?", userID).Pluck("id", &idSlice).Error; err != nil {
		return nil, err
	}
	return idSlice, nil
}

func (ar *AlbumRepository) DeleteAlbumByID(albumID int) error {
	return ar.db.Delete(&models.Album{}, albumID).Error
}

func (ar *AlbumRepository) DeleteAlbumImage(albumID int, imageID int) error {
	return ar.db.Delete(&models.AlbumImage{}, "album_id = ? AND image_id = ?", albumID, imageID).Error
}

func (ar *AlbumRepository) IsOwnerOfAlbum(userID int, albumID int) error {
	var trueUserID int
	ar.db.Model(&models.Album{}).Where("id = ?", albumID).Pluck("user_id", &trueUserID)
	if userID != trueUserID {
		return fmt.Errorf("user is not owner of this album %d", albumID)
	}
	return nil
}
