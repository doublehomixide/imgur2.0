package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"pictureloader/models"
	"sync"
	"time"
)

// UploadFile - Отправляет файл в minio
func (m *MinioProvider) UploadFile(ctx context.Context, object models.ImageUnit, imageName string) (string, error) {
	_, err := m.client.PutObject(
		ctx,
		bucketName, // Константа с именем бакета
		imageName,
		object.Payload,
		object.PayloadSize,
		minio.PutObjectOptions{ContentType: "image/png"},
	)
	return imageName, err
}

// GetFileURL - Получает StorageKey файла с minio
func (m *MinioProvider) GetFileURL(ctx context.Context, imageURL string) (string, error) {
	imgLink, err := m.client.PresignedGetObject(ctx, bucketName, imageURL, time.Minute, nil)
	return imgLink.String(), err
}

func (m *MinioProvider) GetFileURLS(ctx context.Context, imageURLS []string) ([]string, error) {
	var wg sync.WaitGroup
	var imgChan = make(chan string, len(imageURLS))
	for _, imageURL := range imageURLS {
		wg.Add(1)
		go func(img string) {
			defer wg.Done()

			imgLink, err := m.GetFileURL(ctx, img)
			if err != nil {
				return
			}
			imgChan <- imgLink
		}(imageURL)
	}
	go func() {
		wg.Wait()
		close(imgChan)
	}()
	var result []string
	for imgLink := range imgChan {
		if imgLink != "" {
			result = append(result, imgLink)
		}
	}
	return result, nil
}

func (m *MinioProvider) DeleteFileByURL(ctx context.Context, imageURL string) error {
	err := m.client.RemoveObject(ctx, bucketName, imageURL, minio.RemoveObjectOptions{})
	return err
}
