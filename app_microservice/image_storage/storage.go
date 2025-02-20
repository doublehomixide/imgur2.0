package image_storage

import (
	"context"
	"pictureloader/app_microservice/models"
)

type ImageStorage interface {
	Connect() error                                                       // Инициализатор подключения
	UploadFile(context.Context, models.ImageUnit, string) (string, error) // Загрузка файлов
	GetFileURL(context.Context, string) (string, error)
	GetFileURLS(ctx context.Context, imageURLS []string) ([]string, error) // Скачивание файлов
	DeleteFileByURL(ctx context.Context, imageURL string) error
}
