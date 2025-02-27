package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"pictureloader/app_microservice/broker"
	"pictureloader/app_microservice/image_storage"
	"pictureloader/app_microservice/models"
	"unicode/utf8"
)

type PostRepositoryInterface interface {
	CreatePost(ctx context.Context, album *models.Post) error
	CreatePostAndImage(ctx context.Context, postID int, imageSK string) error
	GetUserPostIDs(ctx context.Context, userID int) ([]int, error)
	DeletePostByID(ctx context.Context, albumID int) error
	DeletePostImage(ctx context.Context, postID int, imageSK string) error
	IsOwnerOfPost(ctx context.Context, userID int, albumID int) error
	LikePost(ctx context.Context, postID, userID int) error
	GetMostLikedPosts(ctx context.Context) ([]models.PostUnit, error)
	GetPostOwner(ctx context.Context, postID int) (int, error)
	GetPost(ctx context.Context, postID int) (models.PostUnit, error)
}

type AlbumCacher interface {
	InvalidatePost(ctx context.Context, postID int) (bool, error)
	SetMostLikedPosts(ctx context.Context, posts []models.PostUnit) error
	GetMostLikedPosts(ctx context.Context) ([]models.PostUnit, error)
	GetPost(ctx context.Context, postID int) (models.PostUnit, bool, error)
	SetPost(ctx context.Context, postID int, resultJSON string) (bool, error)
}

type PostService struct {
	database PostRepositoryInterface
	storage  image_storage.ImageStorage
	cache    AlbumCacher
	broker   broker.RabbitBroker
}

func NewPostService(database PostRepositoryInterface, storage image_storage.ImageStorage,
	cacher AlbumCacher, broker broker.RabbitBroker) *PostService {
	return &PostService{
		database: database,
		storage:  storage,
		cache:    cacher,
		broker:   broker,
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
	cachedResult, isExist, err := als.cache.GetPost(ctx, postID)
	if err != nil {
		slog.Error("Get post", "error", err)
	}

	if !isExist {
		slog.Info("No such key in redis", "postID", postID)
	}

	if err == nil && isExist {
		return cachedResult, nil
	}

	post, err := als.database.GetPost(ctx, postID)
	if err != nil {
		slog.Error("Get post", "error", err)
		return models.PostUnit{}, err
	}

	for k, v := range post.Images {
		imageURL, err := als.storage.GetFileURL(ctx, v)
		if err != nil {
			slog.Error("Get image URL", "error", err)
			continue
		}
		post.Images[k] = imageURL
	}

	resultJSON, _ := json.Marshal(post)

	ok, err := als.cache.SetPost(ctx, postID, string(resultJSON))
	if err != nil {
		slog.Error("Set cache for post", "error", err)
	}
	if err == nil && !ok {
		slog.Info("Set post no such post", "postID", postID)
	}

	return post, nil
}

// GetUserPosts userID -> hashmap where key is album ID, value is models.PostUnit
func (als *PostService) GetUserPosts(ctx context.Context, userID int) (map[int]models.PostUnit, error) {
	postIDs, err := als.database.GetUserPostIDs(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make(map[int]models.PostUnit)

	//n+1 problem todo
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

	err = als.database.CreatePostAndImage(ctx, postID, imageSK)
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
	//todo можно в 1 запрос думаю
	err := als.database.IsOwnerOfPost(ctx, userID, postID)
	if err != nil {
		slog.Error("Database delete error", "error", err)
		return err
	}

	err = als.database.DeletePostImage(ctx, postID, imageSK)
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
	result, err := als.cache.InvalidatePost(ctx, postID)
	if err != nil {
		slog.Error("Delete post", "error", err)
		return err
	}
	if result == false {
		slog.Error("No such post (cache)", "postID", postID)
	}

	postOwnerID, err := als.database.GetPostOwner(ctx, postID)
	if err != nil {
		slog.Error("Get post owner", "error", err)
		return err
	}

	err = als.database.LikePost(ctx, postID, userID)
	if err != nil {
		slog.Error("Like post", "error", err)
		return err
	}

	err = als.broker.PublishNewLike(postID, userID, postOwnerID)
	if err != nil {
		slog.Error("Like post", "broker error", err)
	}

	return nil
}

func (als *PostService) GetMostLikedPosts(ctx context.Context) ([]models.PostUnit, error) {
	//пограничный случай если imageDesc у картинок одинаковые в одном посте, нужно будет что-то придумать todo
	posts, err := als.cache.GetMostLikedPosts(ctx)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			slog.Info("Most liked posts arent cached yet", "error", err)
		} else {
			slog.Error("Get most liked posts", "error", err)
		}
	} else {
		return posts, nil
	}

	//если кеш не найден
	posts, err = als.database.GetMostLikedPosts(ctx)
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

	err = als.cache.SetMostLikedPosts(ctx, posts)
	if err != nil {
		slog.Error("Cache most liked posts", "error", err)
		return nil, err
	}
	return posts, nil
}
