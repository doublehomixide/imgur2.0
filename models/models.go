package models

import "io"

type User struct {
	ID       int    `gorm:"primary_key;autoIncrement" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Images   []Image
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
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
