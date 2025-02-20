package likes

import (
	"encoding/json"
	"log/slog"
)

type Message struct {
	PostID int `json:"post_id"`
	Liker  int `json:"liker"`
	Liked  int `json:"liked"`
}

type NotificationService struct {
	repo *LikeNotificationRepository
}

func NewNotificationService(repo *LikeNotificationRepository) *NotificationService {
	return &NotificationService{repo}
}

func (ns *NotificationService) ProcessLikeMessage(message []byte) error {
	var msg Message
	err := json.Unmarshal(message, &msg)
	if err != nil {
		slog.Info("Error unmarshalling message", "error", err)
	}
	err = ns.repo.CreateLikeNotification(msg.PostID, msg.Liker, msg.Liked)
	if err != nil {
		slog.Error("Error creating like notification", "error", err)
		return err
	}
	return nil
}

func (ns *NotificationService) GetAllLikeNotifications(postID int) ([]int, error) {
	result, err := ns.repo.GetAllLikeNotifications(postID)
	if err != nil {
		slog.Error("db error", "error", err)
		return nil, err
	}
	return result, err
}
