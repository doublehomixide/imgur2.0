package likes

import (
	"encoding/json"
	"fmt"
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

func (ns *NotificationService) GetAllLikeNotifications(userID int) ([]string, error) {
	likeNotif, err := ns.repo.GetAllLikeNotifications(userID)
	if err != nil {
		slog.Error("db error", "error", err)
		return nil, err
	}

	var result []string

	for _, el := range likeNotif {
		result = append(result, fmt.Sprintf("User %d liked your post number %d at %02d-%02d-%04d %02d:%02d",
			el.Liker, el.PostID,
			el.CreatedAt.Day(), el.CreatedAt.Month(), el.CreatedAt.Year(),
			el.CreatedAt.Hour(), el.CreatedAt.Minute()))
	}

	return result, err
}
