package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	broker_package "pictureloader/notification_microservice/broker"
	"pictureloader/notification_microservice/database"
	"pictureloader/notification_microservice/notifications/likes"
)

func main() {

	dbConn := database.NewPSQL("postgres://postgres:1000@localhost:5432/db2?sslmode=disable")
	LikesNotifications, likesService := likes.NewLikesNotifications(dbConn)

	broker := broker_package.NewRabbitMQ(likesService)
	go broker.ListenLikes()
	defer broker.Connection.Close()
	defer broker.Channel.Close()

	mainRouter := mux.NewRouter()
	likes.LikesNotificationsRouter(mainRouter, LikesNotifications)
	err := http.ListenAndServe(":8081", mainRouter)
	if err != nil {
		log.Fatal(err)
	}
}
