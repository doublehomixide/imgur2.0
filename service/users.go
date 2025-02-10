package service

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
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
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	return u.database.CreateNewUser(user)
}

func (u *UserService) LoginUser(userLogin *models.UserLogin) (bool, int) {
	user, err := u.database.GetUserByUsername(userLogin.Username)
	if err != nil {
		slog.Error("Login user (get user by username) error", "error", err)
		return false, -1
	}
	if user == nil {
		slog.Info("LoginUser: User not found", "error", err)
		return false, -1
	}

	// Сравниваем хэш пароля с введенным паролем
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
	if err != nil {
		log.Printf("Incorrect password for user: %s", userLogin.Username)
		return false, -1
	}

	return true, user.ID
}

func (u *UserService) DeleteUserByID(userID int) error {
	err := u.database.DeleteUserByID(userID)
	if err != nil {
		slog.Error("DeleteUserByID error", "error", err)
		return err
	}
	return nil
}

func (u *UserService) UpdateUsername(userID int, username string) error {
	err := u.database.ChangeUsernameByID(userID, username)
	if err != nil {
		slog.Error("UpdateUsername error", "error", err)
		return err
	}
	return nil
}
