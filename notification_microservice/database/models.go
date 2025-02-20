package database

import "time"

type LikesNotification struct {
	ID        int64 `gorm:"primary_key"`
	PostID    int
	Liker     int
	Liked     int
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
