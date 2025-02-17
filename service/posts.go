package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"pictureloader/image_storage"
	"pictureloader/models"
	"strconv"
	"time"
	"unicode/utf8"
)

type PostRepositoryInterface interface {
	CreatePost(ctx context.Context, album *models.Post) error
	CreatePostAndImage(ctx context.Context, albumImage *models.PostImage) error
	GetPostData(ctx context.Context, albumID int) (string, map[string]string, error)
	GetUserPostIDs(ctx context.Context, userID int) ([]int, error)
	DeletePostByID(ctx context.Context, albumID int) error
	DeletePostImage(ctx context.Context, albumID int, imageID int) error
	IsOwnerOfPost(ctx context.Context, userID int, albumID int) error
	GetPostLikesCount(ctx context.Context, postID int) (int, error)
	LikePost(ctx context.Context, postID, userID int) error
	GetMostLikedPosts(ctx context.Context) ([]models.PostUnit, error)
}

type AlbumCacher interface {
	InvalidatePost(ctx context.Context, postID int) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type ImageWorker interface {
	GetImageIDBySK(ctx context.Context, imageSK string) (int, error)
}

type PostService struct {
	database      PostRepositoryInterface
	storage       image_storage.ImageStorage
	imageDatabase ImageWorker
	cache         AlbumCacher
}

func NewPostService(database PostRepositoryInterface, storage image_storage.ImageStorage,
	imageDatabase ImageWorker, cacher AlbumCacher) *PostService {
	return &PostService{
		database:      database,
		storage:       storage,
		imageDatabase: imageDatabase,
		cache:         cacher,
	}
}

func (als *PostService) CreatePost(ctx context.Context, post *models.Post) error {
	if post.Name == "" || utf8.RuneCountInString(post.Name) > 30 {
		return errors.New("invalid post name")
	}
	err := als.database.CreatePost(ctx, post)
	if err != nil {
		slog.Error("Create post", "error", err)
		return err
	}
	return nil
}

func (als *PostService) GetPost(ctx context.Context, postID int) (models.PostUnit, error) {

	cachedResult, err := als.cache.Get(ctx, strconv.Itoa(postID))
	if err != nil && !errors.Is(err, redis.Nil) {
		slog.Error("Get post", "error", err)
		return models.PostUnit{}, err
	}
	if cachedResult != "" {
		slog.Info("Get album from cache", "postID", postID)
		var album models.PostUnit
		json.Unmarshal([]byte(cachedResult), &album)
		return album, nil
	}

	likesCount, err := als.database.GetPostLikesCount(ctx, postID)
	if err != nil {
		slog.Error("Get post likes count", "error", err)
		return models.PostUnit{}, err
	}

	postName, images, err := als.database.GetPostData(ctx, postID)
	if err != nil {
		slog.Error("Get images from post", "error", err)
		return models.PostUnit{}, err
	}

	updatedImages := make(map[string]string)

	for key, value := range images {
		presignedURL, err := als.storage.GetFileURL(ctx, key)
		if err != nil {
			continue
		}
		updatedImages[value] = presignedURL
	}

	result := models.PostUnit{
		Name:   postName,
		Images: updatedImages,
		Likes:  likesCount,
	}

	resultJSON, _ := json.Marshal(result)

	err = als.cache.Set(ctx, strconv.Itoa(postID), string(resultJSON), time.Hour*5)
	if err != nil {
		slog.Error("Set cache for post", "error", err)
	}

	return result, nil
}

// GetUserPosts userID -> hashmap where key is album ID, value is models.PostUnit
func (als *PostService) GetUserPosts(ctx context.Context, userID int) (map[int]models.PostUnit, error) {
	postIDs, err := als.database.GetUserPostIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make(map[int]models.PostUnit)

	for _, albumID := range postIDs {
		album, err := als.GetPost(ctx, albumID)
		if err != nil {
			continue
		}
		result[albumID] = album
	}

	return result, nil
}

func (als *PostService) AppendImageToPost(ctx context.Context, postID int, imageSK string, userID int) error {
	err := als.database.IsOwnerOfPost(ctx, userID, postID)
	if err != nil {
		slog.Error("Database delete error", "error", err)
		return err
	}

	imageID, err := als.imageDatabase.GetImageIDBySK(ctx, imageSK)
	if err != nil {
		slog.Info("Add image to post", "error", err)
		return err
	}

	postImageModel := &models.PostImage{PostID: postID, ImageID: imageID}

	err = als.database.CreatePostAndImage(ctx, postImageModel)
	if err != nil {
		slog.Error("Add image to post", "error", err)
		return err
	}

	ok, err := als.cache.InvalidatePost(ctx, postID)
	if err != nil {
		slog.Error("Delete post", "error", err)
		return err
	}
	if ok == false {
		slog.Info("No such post", "postID", postID)
		return errors.New("no such post")
	}

	return nil
}

func (als *PostService) DeletePost(ctx context.Context, postID int, userID int) error {
	err := als.database.IsOwnerOfPost(ctx, userID, postID)
	if err != nil {
		slog.Error("Database delete error", "error", err)
		return err
	}

	err = als.database.DeletePostByID(ctx, postID)
	if err != nil {
		slog.Error("Delete post", "error", err)
		return err
	}

	ok, err := als.cache.InvalidatePost(ctx, postID)
	if err != nil {
		slog.Error("Delete post", "error", err)
		return err
	}
	if ok == false {
		slog.Info("No such post", "postID", postID)
		return errors.New("no such post")
	}

	return nil
}

func (als *PostService) DeleteImageFromPost(ctx context.Context, postID int, imageSK string, userID int) error {
	err := als.database.IsOwnerOfPost(ctx, userID, postID)
	if err != nil {
		slog.Error("Database delete error", "error", err)
		return err
	}

	imageID, err := als.imageDatabase.GetImageIDBySK(ctx, imageSK)
	if err != nil {
		slog.Info("Delete image from post", "error", err)
		return err
	}
	err = als.database.DeletePostImage(ctx, postID, imageID)
	if err != nil {
		slog.Error("Delete image from post", "error", err)
		return err
	}

	ok, err := als.cache.InvalidatePost(ctx, postID)
	if err != nil {
		slog.Error("Delete post", "error", err)
		return err
	}
	if ok == false {
		slog.Info("No such post", "postID", postID)
		return errors.New("no such post")
	}

	return nil
}

func (als *PostService) LikePost(ctx context.Context, postID, userID int) error {
	err := als.cache.Delete(ctx, strconv.Itoa(postID))
	if err != nil {
		slog.Error("Delete post", "error", err)
		return err
	}
	err = als.database.LikePost(ctx, postID, userID)
	if err != nil {
		slog.Error("Like post", "error", err)
		return err
	}
	return nil
}

func (als *PostService) GetMostLikedPosts(ctx context.Context) ([]models.PostUnit, error) {
	//пограничный случай если imageDesc у картинок одинаковые в одном посте, нужно будет что-то придумать todo
	posts, err := als.database.GetMostLikedPosts(ctx)

	for _, post := range posts {
		for k, v := range post.Images {
			imageURL, err := als.storage.GetFileURL(ctx, v)
			if err != nil {
				continue
			}
			delete(post.Images, k)
			post.Images[v] = imageURL
		}
	}
	if err != nil {
		slog.Error("Get most liked posts", "error", err)
		return nil, err
	}
	return posts, nil
}
