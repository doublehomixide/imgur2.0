package models

type User struct {
	ID             int     `gorm:"primary_key;autoIncrement" json:"id"`
	Username       string  `gorm:"unique" json:"username"`
	Email          string  `gorm:"unique" json:"email"`
	Password       string  `json:"password"`
	ProfilePicture string  `gorm:"unique" json:"profilePictureStorageKey"`
	Images         []Image `json:"images"`
	Albums         []Post  `json:"albums"`
}

type UserProfile struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profile_picture"`
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
