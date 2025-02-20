package likes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type LikesServer struct {
	service *NotificationService
}

func NewLikesServer(service *NotificationService) *LikesServer {
	return &LikesServer{service}
}

func LikesNotificationsRouter(api *mux.Router, server *LikesServer) {
	router := api.PathPrefix("/likes").Subrouter()
	router.HandleFunc("/{postID}", server.GetLikesNotifications).Methods("GET")
}

func (server *LikesServer) GetLikesNotifications(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postID"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	userIDs, err := server.service.GetAllLikeNotifications(postID)
	if err != nil {
		http.Error(w, "Failed to get notifications", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userIDs)
}
