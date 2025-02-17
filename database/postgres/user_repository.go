package postgres

import (
	"context"
	"gorm.io/gorm"
	"log"
	"pictureloader/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) CreateNewUser(ctx context.Context, user *models.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u *UserRepository) GetUserByID(ctx context.Context, id int) (*models.UserProfile, error) {
	var user models.UserProfile
	err := u.db.WithContext(ctx).Model(&models.User{}).First(&user, id).Error
	return &user, err
}

func (u *UserRepository) DeleteUserByID(ctx context.Context, id int) error {
	return u.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (u *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := u.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

func (u *UserRepository) ChangeUsernameByID(ctx context.Context, userID int, newUsername string) error {
	err := u.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("username", newUsername).Error
	return err
}

func (u *UserRepository) UpdatePasswordByID(ctx context.Context, userID int, newPass string) error {
	err := u.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("password", newPass).Error
	return err
}

func (u *UserRepository) UploadProfilePicture(ctx context.Context, userID int, imageSK string) error {
	err := u.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("profile_picture", imageSK).Error
	return err
}
