package models

type Post struct {
	ID     int     `gorm:"primaryKey;autoIncrement" json:"post_id"`
	Name   string  `gorm:"not null" json:"name"`
	UserID int     `json:"user_id" gorm:"not null"`
	Likes  []Like  `json:"likes" gorm:"foreignKey:PostID"`
	User   User    `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Images []Image `gorm:"many2many:post_images"`
}

type PostUnit struct {
	Name   string            `json:"name"`
	Images map[string]string `json:"images"`
	Likes  int               `json:"likes_count"`
}

// PostRegister uses only for swagger
type PostRegister struct {
	Name string `json:"name"`
}

type PostImage struct {
	PostID  int `gorm:"primaryKey"`
	ImageID int `gorm:"primaryKey"`
}

type Like struct {
	PostID int  `gorm:"primaryKey"`
	UserID int  `gorm:"primaryKey"`
	User   User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}
