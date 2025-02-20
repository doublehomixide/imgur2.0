package users

import (
	"context"
	"github.com/stretchr/testify/assert"
	"pictureloader/app_microservice/models"
	"pictureloader/app_microservice/service"
	"testing"
)

func setupTest() (*service.UserService, *MockUsersRepository) {
	mockUserRepository := new(MockUsersRepository)
	userService := service.NewUserService(mockUserRepository)
	return userService, mockUserRepository
}

func TestUserService_RegisterUser(t *testing.T) {
	userService, mockUserRepository := setupTest()
	ctx := context.Background()

	user := models.User{
		ID:       1,
		Username: "vaflya",
		Email:    "vaflya@gmail.com",
		Password: "vaflya228",
		Images:   nil,
		Albums:   nil,
	}
	mockUserRepository.On("CreateNewUser", ctx, &user).Return(nil)
	err := userService.RegisterUser(ctx, &user)

	assert.NoError(t, err)
	mockUserRepository.AssertExpectations(t)
}

func TestUserService_RegisterUser_Error(t *testing.T) {
	userService, mockUserRepository := setupTest()
	ctx := context.Background()

	user := models.User{
		ID:       1,
		Username: "",
		Email:    "vaflya@gmail.com",
		Password: "vaflya228",
		Images:   nil,
		Albums:   nil,
	}
	err := userService.RegisterUser(ctx, &user)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid username")

	user = models.User{
		ID:       1,
		Username: "vaflya vaflya",
		Email:    "vaflya@gmail.com",
		Password: "vaflya228",
		Images:   nil,
		Albums:   nil,
	}
	err = userService.RegisterUser(ctx, &user)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "invalid username")

	mockUserRepository.AssertNotCalled(t, "CreateNewUser")
}
