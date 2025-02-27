package rest

import (
	"context"
	"encoding/json"
	"github.com/go-chi/httprate"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"net/http"
	"pictureloader/app_microservice/models"
	"pictureloader/app_microservice/safety/jwtutils"
	"pictureloader/app_microservice/service"
	"strconv"
	"time"
)

type PostServer struct {
	service service.PostService
}

func NewPostServer(service service.PostService) *PostServer {
	return &PostServer{service: service}
}

func PostRouter(api *mux.Router, server *PostServer) {
	jwtUtils := jwtutils.UtilsJWT{}
	router := api.PathPrefix("/posts").Subrouter()
	router.HandleFunc("", server.CreatePostHandler).Methods("POST")
	router.HandleFunc("/my", server.GetMyPosts).Methods("GET")
	router.HandleFunc("/most-liked", server.GetMostLikedPosts).Methods("GET")
	router.HandleFunc("/{postID}", server.GetPost).Methods("GET")
	router.HandleFunc("/{postID}/like", server.LikePostHandler).Methods("POST")
	router.HandleFunc("/{postID}/{imageSK}", server.AddImageToPost).Methods("POST")
	router.HandleFunc("/{postID}", server.DeletePost).Methods("DELETE")
	router.HandleFunc("/{postID}/{imageSK}", server.DeletePostImage).Methods("DELETE")
	router.Use(jwtUtils.AuthMiddleware)
	router.Use(httprate.LimitByRealIP(3, 3*time.Second))
}

// CreatePostHandler creates a new post for the user.
// @Summary Create a new post
// @Description Creates a new post for the user.
// @Tags Posts
// @Accept json
// @Produce json
// @Param post body models.PostRegister true "Post"
// @Router /posts [post]
func (ps *PostServer) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)

	var post models.Post
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	post.UserID = userID

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err = ps.service.CreatePost(ctx, &post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"status":"Post created"}`))
}

// GetPost retrieves images from a specific post.
// @Summary Get post details
// @Description Retrieves images and details from a specific post.
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path int true "Post ID"
// @Router /posts/{postID} [get]
func (ps *PostServer) GetPost(w http.ResponseWriter, r *http.Request) {
	postIDstr := mux.Vars(r)["postID"]
	postID, err := strconv.Atoi(postIDstr)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := ps.service.GetPost(ctx, postID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	encoder.Encode(result)
}

// AddImageToPost adds an image to a post.
// @Summary Add an image to a post
// @Description Adds an image to the specified post.
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path int true "Post ID"
// @Param imageSK path string true "Image Storage Key"
// @Router /posts/{postID}/{imageSK} [post]
func (ps *PostServer) AddImageToPost(w http.ResponseWriter, r *http.Request) {
	postIDstr := mux.Vars(r)["postID"]
	imageSK := mux.Vars(r)["imageSK"]
	postID, _ := strconv.Atoi(postIDstr)

	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := ps.service.AppendImageToPost(ctx, postID, imageSK, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"status":"Image added"}`))
}

// GetMyPosts retrieves all posts of the user.
// @Summary Get all posts of the user
// @Description Retrieves all posts of the currently authenticated user.
// @Tags Posts
// @Accept json
// @Produce json
// @Router /posts/my [get]
func (ps *PostServer) GetMyPosts(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	result, err := ps.service.GetUserPosts(ctx, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	encoder.Encode(result)
}

// DeletePost removes a specific post and its images from the system.
// @Summary Delete a post
// @Description Removes a specific post and all its images from the system.
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path int true "Post ID"
// @Router /posts/{postID} [delete]
func (ps *PostServer) DeletePost(w http.ResponseWriter, r *http.Request) {
	postIDstr := mux.Vars(r)["postID"]
	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err = ps.service.DeletePost(ctx, postID, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"status":"Post deleted"}`))
}

// DeletePostImage removes an image from a post.
// @Summary Delete an image from a post
// @Description Removes a specific image from the specified post.
// @Tags Posts
// @Accept json
// @Produce json
// @Param postID path int true "Post ID"
// @Param imageSK path string true "Image Storage Key"
// @Router /posts/{postID}/{imageSK} [delete]
func (ps *PostServer) DeletePostImage(w http.ResponseWriter, r *http.Request) {
	postIDstr := mux.Vars(r)["postID"]
	imageSK := mux.Vars(r)["imageSK"]
	postID, _ := strconv.Atoi(postIDstr)

	claims := r.Context().Value("claims").(jwt2.MapClaims)
	sub := claims["sub"].(float64)
	userID := int(sub)

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := ps.service.DeleteImageFromPost(ctx, postID, imageSK, userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"status":"Image deleted"}`))
}

// LikePostHandler handles liking a post.
// @Summary     Like a post
// @Description Increments the like count of a post and invalidates cache
// @Tags        Posts
// @Accept      json
// @Produce     json
// @Param       postID path int true "Post ID"
// @Router      /posts/{postID}/like [post]
func (ps *PostServer) LikePostHandler(w http.ResponseWriter, r *http.Request) {
	postIDstr := mux.Vars(r)["postID"]
	postID, err := strconv.Atoi(postIDstr)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	sub := r.Context().Value("claims").(jwt2.MapClaims)["sub"]
	userID := int(sub.(float64))

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err = ps.service.LikePost(ctx, postID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"Post liked successfully"}`))
}

// GetMostLikedPosts returns the most liked posts.
// @Summary     Get most liked posts
// @Description Returns a list of the most liked posts, ordered by like count in descending order.
// @Tags        Posts
// @Accept      json
// @Produce     json
// @Router      /posts/most-liked [get]
func (ps *PostServer) GetMostLikedPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := ps.service.GetMostLikedPosts(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	encoder.Encode(posts)
}
