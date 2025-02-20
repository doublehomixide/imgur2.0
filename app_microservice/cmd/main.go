package main

import (
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"log"
	"log/slog"
	"net/http"
	"os"
	"pictureloader/app_microservice/broker"
	"pictureloader/app_microservice/caching/redis"
	config "pictureloader/app_microservice/cfg"
	postgres2 "pictureloader/app_microservice/database/postgres"
	_ "pictureloader/app_microservice/docs"
	"pictureloader/app_microservice/image_storage/minio"
	rest2 "pictureloader/app_microservice/rest"
	service2 "pictureloader/app_microservice/service"
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
	psqlDB := postgres2.NewDataBase(cfg.PsqlDBPath)
	slog.Info("Postgres DB initialized")

	//database and service related to the database init
	userRepo := postgres2.NewUserRepository(psqlDB)
	imageRepo := postgres2.NewImageRepository(psqlDB)
	postRepo := postgres2.NewPostRepository(psqlDB)
	slog.Info("Image and User repositories initialized")

	cache := redis.NewRedisClient(imageRepo)
	rabbitbroker := broker.NewRabbitBroker()
	defer rabbitbroker.Close()

	imageService := service2.NewPictureLoader(minioprov, imageRepo, cache)
	userService := service2.NewUserService(userRepo, minioprov)
	postService := service2.NewPostService(postRepo, minioprov, imageRepo, cache, *rabbitbroker)
	slog.Info("Image and User services initialized")

	picturesServer := rest2.PictureNewServer(imageService)
	userServer := rest2.NewUserServer(userService)
	albumServer := rest2.NewPostServer(*postService)

	slog.Info("User and Image server initialized")

	//router init
	mainRouter := mux.NewRouter()
	mainRouter.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	rest2.PictureRouter(mainRouter, picturesServer)
	rest2.UserRouter(mainRouter, userServer)
	rest2.PostRouter(mainRouter, albumServer)
	slog.Info("Routers are running")

	slog.Info("Starting server on port ", "port", cfg.ServerPort)
	if err = http.ListenAndServe(cfg.ServerPort, mainRouter); err != nil {
		slog.Error("PictureServer failed to start", "error", err)
		os.Exit(1)
	}
}
