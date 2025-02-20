package albums

import (
	"context"
	"github.com/stretchr/testify/assert"
	"pictureloader/app_microservice/models"
	"pictureloader/app_microservice/service"
	"testing"
)

type mockAlbumService struct {
}

func setupTest() (*service.PostService, *MockAlbumRepository, *MockImageWorker, *MockImageStorage, *MockCacher) {
	mockAlbumRepo := new(MockAlbumRepository)
	mockImageWorker := new(MockImageWorker)
	mockImageStorage := new(MockImageStorage)
	mockCacher := new(MockCacher)

	albumService := service.NewPostService(mockAlbumRepo, mockImageStorage, mockImageWorker, mockCacher)

	return albumService, mockAlbumRepo, mockImageWorker, mockImageStorage, mockCacher
}

func TestAlbumService_CreateAlbum(t *testing.T) {
	ctx := context.Background()
	albumService, mockAlbumRepo, _, _, _ := setupTest()

	album := &models.Post{ID: 1, Name: "Test Post", UserID: 1}

	mockAlbumRepo.On("CreatePost", ctx, album).Return(nil)

	err := albumService.CreatePost(ctx, album)

	assert.NoError(t, err)
	mockAlbumRepo.AssertExpectations(t)
}

func TestAlbumService_CreateAlbum_Error(t *testing.T) {
	ctx := context.Background()
	albumService, mockAlbumRepo, _, _, _ := setupTest()

	album := &models.Post{ID: 1, Name: "1234567890123456789012345678901", UserID: 1} //31 symbol

	err := albumService.CreatePost(ctx, album)

	assert.Error(t, err)
	assert.Equal(t, "invalid album name", err.Error())

	mockAlbumRepo.AssertExpectations(t)
	mockAlbumRepo.AssertNotCalled(t, "CreatePost")
}
