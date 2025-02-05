package main

import (
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"log"
	"log/slog"
	"net/http"
	"os"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	cfg := config.Init()
	//minio and image storage init
	minioprov, err := minio.NewMinioProvider(cfg.MinioURL, cfg.MinioUSER, cfg.MinioPASSWORD, false)
	if err != nil {
		log.Fatalf("Failed to initialize Minio provider: %v", err)
	}
	slog.Info("Minio provider initialized")
	psqlDB := postgres.NewDataBase(cfg.PsqlDBPath)
	slog.Info("Postgres DB initialized")

	//database and service related to the database init
	userRepo := postgres.NewUserRepository(psqlDB)
	imageRepo := postgres.NewImageRepository(psqlDB)
	slog.Info("Image and User repositories initialized")

	imageService := service.NewPictureLoader(minioprov, imageRepo)
	userService := service.NewUserService(*userRepo)
	slog.Info("Image and User services initialized")

	picturesServer := rest.PictureNewServer(imageService)
	userServer := rest.NewUserServer(userService)

	slog.Info("User and Image server initialized")

	//router init
	mainRouter := mux.NewRouter()
	mainRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	rest.PictureRouter(mainRouter, picturesServer)
	rest.UserRouter(mainRouter, userServer)
	slog.Info("Routers are running")

	slog.Info("Starting server on port %s", cfg.ServerPort)
	if err = http.ListenAndServe(cfg.ServerPort, mainRouter); err != nil {
		slog.Error("PictureServer failed to start: %v", err)
		os.Exit(1)
	}
}
