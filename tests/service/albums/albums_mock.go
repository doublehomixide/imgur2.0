package albums

import (
	"context"
	"github.com/stretchr/testify/mock"
	"pictureloader/models"
	"time"
)

type MockAlbumRepository struct {
	mock.Mock
}

func (m *MockAlbumRepository) CreateAlbum(ctx context.Context, album *models.Album) error {
	args := m.Called(ctx, album)
	return args.Error(0)
}

func (m *MockAlbumRepository) CreateAlbumAndImage(ctx context.Context, albumImage *models.AlbumImage) error {
	args := m.Called(ctx, albumImage)
	return args.Error(0)
}

func (m *MockAlbumRepository) GetAlbumData(ctx context.Context, albumID int) (string, map[string]string, error) {
	args := m.Called(ctx, albumID)
	return args.String(0), args.Get(1).(map[string]string), args.Error(2)
}
func (m *MockAlbumRepository) GetUserAlbumIDs(ctx context.Context, userID int) ([]int, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]int), args.Error(1)
}

func (m *MockAlbumRepository) DeleteAlbumByID(ctx context.Context, albumID int) error {
	args := m.Called(ctx, albumID)
	return args.Error(0)
}
func (m *MockAlbumRepository) DeleteAlbumImage(ctx context.Context, albumID int, imageID int) error {
	args := m.Called(ctx, albumID, imageID)
	return args.Error(0)
}

func (m *MockAlbumRepository) IsOwnerOfAlbum(ctx context.Context, userID int, albumID int) error {
	args := m.Called(ctx, userID, albumID)
	return args.Error(0)
}

////////////////////

type MockImageWorker struct {
	mock.Mock
}

func (m *MockImageWorker) GetImageIDBySK(ctx context.Context, imageSK string) (int, error) {
	args := m.Called(ctx, imageSK)
	return args.Int(0), args.Error(1)
}

////////////////////

type MockImageStorage struct {
	mock.Mock
}

func (m *MockImageStorage) Connect() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockImageStorage) UploadFile(ctx context.Context, img models.ImageUnit, name string) (string, error) {
	args := m.Called(ctx, img, name)
	return args.String(0), args.Error(1)
}

func (m *MockImageStorage) GetFileURL(ctx context.Context, sk string) (string, error) {
	args := m.Called(ctx, sk)
	return args.String(0), args.Error(1)
}

func (m *MockImageStorage) GetFileURLS(ctx context.Context, imageURLS []string) ([]string, error) {
	args := m.Called(ctx, imageURLS)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockImageStorage) DeleteFileByURL(ctx context.Context, imageURL string) error {
	args := m.Called(ctx, imageURL)
	return args.Error(0)
}

////////////////////

type MockCacher struct {
	mock.Mock
}

func (m *MockCacher) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}
func (m *MockCacher) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockCacher) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}
