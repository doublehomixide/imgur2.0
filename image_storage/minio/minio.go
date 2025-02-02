package minio

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"pictureloader/image_storage"
)

type minioAuthData struct {
	url      string
	user     string
	password string
	token    string
	ssl      bool
}

type MinioProvider struct {
	minioAuthData
	client *minio.Client
}

const bucketName = "usersphotos"

// NewMinioProvider инициализирует нового клиента Minio
func NewMinioProvider(minioURL string, minioUser string, minioPassword string, ssl bool) (image_storage.ImageStorage, error) {
	client, err := minio.New(minioURL, &minio.Options{
		Creds:  credentials.NewStaticV4(minioUser, minioPassword, ""),
		Secure: ssl,
	})
	if err != nil {
		return nil, err
	}

	return &MinioProvider{
		minioAuthData: minioAuthData{
			password: minioPassword,
			url:      minioURL,
			user:     minioUser,
			ssl:      ssl,
		},
		client: client,
	}, nil
}

func (m *MinioProvider) Connect() error {
	var err error
	m.client, err = minio.New(m.url, &minio.Options{
		Creds:  credentials.NewStaticV4(m.user, m.password, ""),
		Secure: m.ssl,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return err
}
