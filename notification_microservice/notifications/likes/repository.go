package likes

import (
	"gorm.io/gorm"
	"pictureloader/notification_microservice/database"
	"time"
)

func NewPSQLNotificationsRepository(db *gorm.DB) *LikeNotificationRepository {
	return &LikeNotificationRepository{db}
}

type LikeNotificationRepository struct {
	DB *gorm.DB
}

func (np *LikeNotificationRepository) CreateLikeNotification(postID, likerID, likedID int) error {
	err := np.DB.Create(&database.LikesNotification{PostID: postID, Liker: likerID, Liked: likedID}).Error
	if err != nil {
		return err
	}
	return nil
}

type LikeNotification struct {
	PostID    int       `json:"post_id"`
	Liker     int       `json:"liker"`
	CreatedAt time.Time `json:"created_at"`
}

func (np *LikeNotificationRepository) GetAllLikeNotifications(likedID int) ([]LikeNotification, error) {
	var likeNotif []LikeNotification
	err := np.DB.Model(&database.LikesNotification{}).Where("liked = ?", likedID).Find(&likeNotif).Error
	if err != nil {
		return nil, err
	}
	return likeNotif, nil
}
