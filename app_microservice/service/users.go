package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
	"pictureloader/app_microservice/models"
	"strings"
)

type UserRepositoryInterface interface {
	CreateNewUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int) (*models.UserProfile, error)
	DeleteUserByID(ctx context.Context, id int) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	ChangeUsernameByID(ctx context.Context, userID int, newUsername string) error
	UpdatePasswordByID(ctx context.Context, userID int, newPassword string) error
	UploadProfilePicture(ctx context.Context, userID int, imageSK string) error
}

type UserStorageManager interface {
	GetFileURL(context.Context, string) (string, error)
}

type UserService struct {
	database UserRepositoryInterface
	storage  UserStorageManager
}

func NewUserService(database UserRepositoryInterface, storage UserStorageManager) *UserService {
	return &UserService{database, storage}
}

func (u *UserService) RegisterUser(ctx context.Context, user *models.User) error {
	if user.Username == "" || strings.Contains(user.Username, " ") {
		return errors.New("invalid username")
	}

	if user.Email == "" || strings.Contains(user.Email, " ") || !strings.Contains(user.Email, "@") {
		return errors.New("invalid email")
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	return u.database.CreateNewUser(ctx, user)
}

func (u *UserService) LoginUser(ctx context.Context, userLogin *models.UserLogin) (bool, int) {
	user, err := u.database.GetUserByUsername(ctx, userLogin.Username)
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

func (u *UserService) DeleteUserByID(ctx context.Context, userID int) error {
	err := u.database.DeleteUserByID(ctx, userID)
	if err != nil {
		slog.Error("DeleteUserByID error", "error", err)
		return err
	}
	return nil
}

func (u *UserService) UpdateUsername(ctx context.Context, userID int, username string) error {
	err := u.database.ChangeUsernameByID(ctx, userID, username)
	if err != nil {
		slog.Error("UpdateUsername error", "error", err)
		return err
	}
	return nil
}

func (u *UserService) UpdatePassword(ctx context.Context, userID int, newPassword string) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("UpdatePassword error", "error", err)
		return err
	}
	newPassword = string(hashedPass)
	err = u.database.UpdatePasswordByID(ctx, userID, newPassword)
	if err != nil {
		slog.Error("UpdateUsername error", "error", err)
		return err
	}
	return nil
}

func (u *UserService) GetUserByID(ctx context.Context, userID int) (*models.UserProfile, error) {
	user, err := u.database.GetUserByID(ctx, userID)
	if err != nil {
		slog.Error("GetUserByID error", "error", err)
		return nil, err
	}
	if user.ProfilePicture != "" {
		user.ProfilePicture, err = u.storage.GetFileURL(ctx, user.ProfilePicture)
		if err != nil {
			slog.Error("GetUserByID error", "error", err)
			return nil, err
		}
	}

	if user == nil {
		slog.Info("GetUserByID: User not found", "error", errors.New("user not found"))
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *UserService) UploadProfilePicture(ctx context.Context, userID int, imageSK string) error {
	err := u.database.UploadProfilePicture(ctx, userID, imageSK)
	if err != nil {
		slog.Error("UploadProfilePicture error", "error", err)
		return err
	}
	return nil
}
