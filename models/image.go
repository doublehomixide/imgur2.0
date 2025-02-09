package models

import "io"

type Image struct {
	ID          int    `gorm:"primary_key" json:"id"`
	StorageKey  string `json:"url"`
	UserID      int    `json:"user_id"`
	Description string `json:"description" gorm:"size:150"`
}

type ImageUnit struct {
	User
	Payload     io.Reader
	PayloadName string
	PayloadSize int64
}
