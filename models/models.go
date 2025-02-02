package models

import "io"

type User struct {
	ID             int    `gorm:"primary_key;autoIncrement" json:"id"`
	Username       string `gorm:"unique" json:"username"`
	Email          string `gorm:"unique" json:"email"`
	HashedPassword string `json:"hashed_password"`
	Images         []Image
}

type Image struct {
	URL    string `gorm:"primary_key" json:"url"`
	UserID int    `json:"user_id"`
}

type ImageUnit struct {
	User
	Payload     io.Reader
	PayloadName string
	PayloadSize int64
}
