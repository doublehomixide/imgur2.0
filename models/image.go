package models

import "io"

type Image struct {
	ID          int    `gorm:"primary_key" json:"id"`
	StorageKey  string `json:"storage_key" gorm:"not null"`
	UserID      int    `json:"user_id" gorm:"not null"`
	Description string `json:"description" gorm:"size:150"`
}

type ImageUnit struct {
	User
	Payload     io.Reader
	PayloadName string
	PayloadSize int64
}
