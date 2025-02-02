package postgres

import (
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

func (u *UserRepository) CreateNewUser(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := u.db.First(&user, id).Error
	return &user, err
}

func (u *UserRepository) DeleteUserByID(id int) error {
	return u.db.Delete(&models.User{}, id).Error
}

func (u *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := u.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Println(err)
	}
	return &user, nil
}
