package likes

import (
	"gorm.io/gorm"
	"pictureloader/notification_microservice/database"
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

func (np *LikeNotificationRepository) GetAllLikeNotifications(postID int) ([]int, error) {
	var userIDs []int
	err := np.DB.Model(&database.LikesNotification{}).Where("post_id = ?", postID).Pluck("user_id", &userIDs).Error
	if err != nil {
		return nil, err
	}
	return userIDs, nil
}
