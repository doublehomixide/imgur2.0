package albums

import (
	"context"
	"github.com/stretchr/testify/assert"
	"pictureloader/models"
	"pictureloader/service"
	"testing"
)

type mockAlbumService struct {
}

func setupTest() (*service.AlbumService, *MockAlbumRepository, *MockImageWorker, *MockImageStorage, *MockCacher) {
	mockAlbumRepo := new(MockAlbumRepository)
	mockImageWorker := new(MockImageWorker)
	mockImageStorage := new(MockImageStorage)
	mockCacher := new(MockCacher)

	albumService := service.NewAlbumService(mockAlbumRepo, mockImageStorage, mockImageWorker, mockCacher)

	return albumService, mockAlbumRepo, mockImageWorker, mockImageStorage, mockCacher
}

func TestAlbumService_CreateAlbum(t *testing.T) {
	ctx := context.Background()
	albumService, mockAlbumRepo, _, _, _ := setupTest()

	album := &models.Album{ID: 1, Name: "Test Album", UserID: 1}

	mockAlbumRepo.On("CreateAlbum", ctx, album).Return(nil)

	err := albumService.CreateAlbum(ctx, album)

	assert.NoError(t, err)
	mockAlbumRepo.AssertExpectations(t)
}

func TestAlbumService_CreateAlbum_Error(t *testing.T) {
	ctx := context.Background()
	albumService, mockAlbumRepo, _, _, _ := setupTest()

	album := &models.Album{ID: 1, Name: "1234567890123456789012345678901", UserID: 1} //31 symbol

	err := albumService.CreateAlbum(ctx, album)

	assert.Error(t, err)
	assert.Equal(t, "invalid album name", err.Error())

	mockAlbumRepo.AssertExpectations(t)
	mockAlbumRepo.AssertNotCalled(t, "CreateAlbum")
}
