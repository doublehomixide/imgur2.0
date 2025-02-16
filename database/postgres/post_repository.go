package postgres

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"pictureloader/models"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (pr *PostRepository) CreatePost(ctx context.Context, post *models.Post) error {
	return pr.db.WithContext(ctx).Create(&post).Error
}

func (pr *PostRepository) CreatePostAndImage(ctx context.Context, postImage *models.PostImage) error {
	return pr.db.WithContext(ctx).Create(&postImage).Error
}

// GetPostData returns post name, hashmap where key is s3 storage StorageKey, value is description, error
func (pr *PostRepository) GetPostData(ctx context.Context, postID int) (string, map[string]string, error) {
	var post models.Post
	if err := pr.db.WithContext(ctx).Preload("Images").First(&post, postID).Error; err != nil {
		return "", nil, err
	}

	images := make(map[string]string)
	for _, image := range post.Images {
		images[image.StorageKey] = image.Description
	}

	return post.Name, images, nil
}

func (pr *PostRepository) GetUserPostIDs(ctx context.Context, userID int) ([]int, error) {
	var idSlice []int
	if err := pr.db.WithContext(ctx).Model(&models.Post{}).Where("user_id = ?", userID).Pluck("id", &idSlice).Error; err != nil {
		return nil, err
	}
	return idSlice, nil
}

func (pr *PostRepository) DeletePostByID(ctx context.Context, postID int) error {
	return pr.db.WithContext(ctx).Delete(&models.Post{}, postID).Error
}

func (pr *PostRepository) DeletePostImage(ctx context.Context, postID int, imageID int) error {
	return pr.db.WithContext(ctx).Delete(&models.PostImage{}, "post_id = ? AND image_id = ?", postID, imageID).Error
}

func (pr *PostRepository) IsOwnerOfPost(ctx context.Context, userID int, postID int) error {
	var trueUserID int
	pr.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", postID).Pluck("user_id", &trueUserID)
	if userID != trueUserID {
		return fmt.Errorf("user is not owner of this post %d", postID)
	}
	return nil
}

func (pr *PostRepository) GetPostLikesCount(ctx context.Context, postID int) (int, error) {
	var count int64
	err := pr.db.WithContext(ctx).Model(&models.Like{}).Where("post_id = ?", postID).Count(&count).Error
	return int(count), err
}

func (pr *PostRepository) LikePost(ctx context.Context, postID, userID int) error {
	err := pr.db.WithContext(ctx).Create(&models.Like{PostID: postID, UserID: userID}).Error
	if err != nil {
		return err
	}
	return nil
}
