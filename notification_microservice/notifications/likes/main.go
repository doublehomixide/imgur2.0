package likes

import "gorm.io/gorm"

func NewLikesNotifications(db *gorm.DB) (*LikesServer, *NotificationService) {
	notificationRepository := NewPSQLNotificationsRepository(db)
	notificationService := NewNotificationService(notificationRepository)
	notificationServer := NewLikesServer(notificationService)
	return notificationServer, notificationService
}
