package main

import (
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	config "pictureloader/cfg"
	"pictureloader/database/postgres"
	_ "pictureloader/docs"
	"pictureloader/image_storage/minio"
	"pictureloader/rest"
	"pictureloader/service"
)

// @title Imgur 2.0 API
// @version 1.0
// @description API для загрузки и просмотра картинок с регистрацией
// @host localhost:8080
func main() {
	cfg := config.Init()
	//minio and image storage init
	minioprov, err := minio.NewMinioProvider(cfg.MinioURL, cfg.MinioUSER, cfg.MinioPASSWORD, false)
	if err != nil {
		log.Fatalf("Failed to initialize Minio provider: %v", err)
	}
	log.Println("Minio provider initialized")
	psqlDB := postgres.NewDataBase(cfg.PsqlDBPath)
	log.Print("Postgres DB initialized\n\n")

	//database and service related to the database init
	userRepo := postgres.NewUserRepository(psqlDB)
	imageRepo := postgres.NewImageRepository(psqlDB)
	log.Print("Image and User repositories initialized\n\n")

	imageService := service.NewPictureLoader(minioprov, imageRepo)
	userService := service.NewUserService(*userRepo)
	log.Print("Image and User services initialized\n\n")

	picturesServer := rest.PictureNewServer(imageService)
	userServer := rest.NewUserServer(userService)

	log.Print("User and Image server initialized\n\n")

	//router init
	mainRouter := mux.NewRouter()
	mainRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	rest.PictureRouter(mainRouter, picturesServer)
	rest.UserRouter(mainRouter, userServer)
	log.Println("Routers are running")

	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(cfg.ServerPort, mainRouter); err != nil {
		log.Fatalf("PictureServer failed to start: %v", err)
	}
}
