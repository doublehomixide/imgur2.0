package users

import (
	"context"
	"github.com/stretchr/testify/mock"
	"pictureloader/models"
)

type MockUsersRepository struct {
	mock.Mock
}

func (m *MockUsersRepository) CreateNewUser(ctx context.Context, user *models.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *MockUsersRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUsersRepository) DeleteUserByID(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func (m *MockUsersRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUsersRepository) ChangeUsernameByID(ctx context.Context, userID int, newUsername string) error {
	return m.Called(ctx, userID, newUsername).Error(0)
}
