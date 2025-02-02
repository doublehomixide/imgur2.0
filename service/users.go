package service

import (
	"pictureloader/database/postgres"
	"pictureloader/models"
)

type UserService struct {
	database postgres.UserRepository
}

func NewUserService(database postgres.UserRepository) *UserService {
	return &UserService{database}
}

func (u *UserService) RegisterUser(user *models.User) error {
	user.HashedPassword = user.HashedPassword + "_HashedPassword"
	return u.database.CreateNewUser(user)
}

func (u *UserService) LoginUser(username string, password string) (bool, int) {
	user, err := u.database.GetUserByUsername(username)
	if user == nil || err != nil {
		return false, -1
	}
	if user.HashedPassword == password+"_HashedPassword" {
		return true, user.ID
	}
	return false, -1
}
