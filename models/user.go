package models

type User struct {
	ID       int     `gorm:"primary_key;autoIncrement" json:"id"`
	Username string  `gorm:"unique" json:"username"`
	Email    string  `gorm:"unique" json:"email"`
	Password string  `json:"password"`
	Images   []Image `json:"images"`
	Albums   []Album `json:"albums"`
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
