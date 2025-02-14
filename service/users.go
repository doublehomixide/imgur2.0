package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"log"
	"log/slog"
	"pictureloader/models"
)

type UserRepositoryInterface interface {
	CreateNewUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	DeleteUserByID(ctx context.Context, id int) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	ChangeUsernameByID(ctx context.Context, userID int, newUsername string) error
}

type UserService struct {
	database UserRepositoryInterface
}

func NewUserService(database UserRepositoryInterface) *UserService {
	return &UserService{database}
}

func (u *UserService) RegisterUser(ctx context.Context, user *models.User) error {
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
