package minio

import (
	"context"
	"fmt"
	"log"
	"pictureloader/database/postgres"
	"pictureloader/image_storage/minio"
	"testing"
)

// TestMinioGetFileUrls isn't test, I know it
func TestMinioGetFileUrls(t *testing.T) {
	minioprov, err := minio.NewMinioProvider("localhost:9000", "minioadmin", "minioadmin", false)
	if err != nil {
		log.Fatalf("Failed to initialize Minio provider: %v", err)
	}
	log.Println("Minio provider initialized")

	psqlDB := postgres.NewDataBase("postgres://postgres:1000@localhost:5432/db?sslmode=disable")
	log.Print("Postgres DB initialized\n\n")

	imageREPO := postgres.NewImageRepository(psqlDB)
	userIDS, err := imageREPO.GetUserImagesID(4)
	fmt.Println(userIDS)

	ctx := context.Background()
	result, err := minioprov.GetFileURLS(ctx, userIDS)
	if err != nil {
		log.Fatalf("Failed to get file urls: %v", err)
	}
	log.Println(result)
}
