package models

type Album struct {
	ID     int     `gorm:"primaryKey;autoIncrement" json:"album_id"`
	Name   string  `gorm:"not null" json:"name"`
	UserID int     `json:"user_id"`
	User   User    `gorm:"foreignKey:UserID;references:ID"`
	Images []Image `gorm:"many2many:album_images"`
}

type AlbumUnit struct {
	Name   string            `json:"name"`
	Images map[string]string `json:"images"`
}

// AlbumRegister uses only for swagger
type AlbumRegister struct {
	Name string `json:"name"`
}

type AlbumImage struct {
	AlbumID int   `gorm:"primaryKey"`
	ImageID int   `gorm:"primaryKey"`
	Album   Album `gorm:"foreignKey:AlbumID;references:ID;constraint:OnDelete:CASCADE"`
	Image   Image `gorm:"foreignKey:ImageID;references:ID;constraint:OnDelete:CASCADE"`
}
