package models

type Post struct {
	ID     int     `gorm:"primaryKey;autoIncrement" json:"post_id"`
	Name   string  `gorm:"not null" json:"name"`
	UserID int     `json:"user_id" gorm:"not null"`
	User   User    `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Images []Image `gorm:"many2many:post_images"`
}

type PostUnit struct {
	Name   string            `json:"name"`
	Images map[string]string `json:"images"`
}

// PostRegister uses only for swagger
type PostRegister struct {
	Name string `json:"name"`
}

type PostImage struct {
	PostID  int   `gorm:"primaryKey"`
	ImageID int   `gorm:"primaryKey"`
	Post    Post  `gorm:"foreignKey:PostID;references:ID;constraint:OnDelete:CASCADE"`
	Image   Image `gorm:"foreignKey:ImageID;references:ID;constraint:OnDelete:CASCADE"`
}
