package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"pictureloader/app_microservice/models"
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

func (pr *PostRepository) CreatePostAndImage(ctx context.Context, postID int, imageSK string) error {
	query := `
		INSERT INTO post_images (post_id, image_id)
		SELECT ?, id FROM images WHERE storage_key = ?`

	result := pr.db.WithContext(ctx).Exec(query, postID, imageSK)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no such image with SK %s", imageSK)
	}

	return nil
}

func (pr *PostRepository) GetPost(ctx context.Context, postID int) (models.PostUnit, error) {
	type postsDBStruct struct {
		Name   string          `json:"name"`
		Images json.RawMessage `json:"images"`
		Likes  int             `json:"likes_count"`
	}
	var post postsDBStruct
	err := pr.db.Model(&models.Post{}).WithContext(ctx).
		Raw(`SELECT posts.name AS Name,
       COALESCE(likes_count, 0) AS Likes,
       COALESCE(JSON_OBJECT_AGG(images.description, images.storage_key), '[]') AS Images
FROM posts
LEFT JOIN (
    SELECT post_id, COUNT(*) AS likes_count
    FROM likes
    GROUP BY post_id
) AS like_counts ON like_counts.post_id = posts.id
LEFT JOIN post_images ON post_images.post_id = posts.id
LEFT JOIN images ON post_images.image_id = images.id
WHERE posts.id = ?
GROUP BY posts.name, posts.id, like_counts.likes_count;
       `, postID).Scan(&post).Error
	if err != nil {
		return models.PostUnit{}, err
	}

	images := make(map[string]string)
	err = json.Unmarshal(post.Images, &images)
	if err != nil {
		return models.PostUnit{}, err
	}
	result := models.PostUnit{Name: post.Name, Images: images, Likes: post.Likes}

	return result, nil
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

func (pr *PostRepository) DeletePostImage(ctx context.Context, postID int, imageSK string) error {
	query := `DELETE FROM post_images 
		USING images 
		WHERE post_images.post_id = ? 
		AND post_images.image_id = images.id 
		AND images.storage_key = ?`
	result := pr.db.WithContext(ctx).Exec(query, postID, imageSK)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no such post-image relation for post_id %d and image_sk %s", postID, imageSK)
	}
	return nil
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

func (pr *PostRepository) UnlikePost(ctx context.Context, postID, userID int) error {
	return nil
}

func (pr *PostRepository) GetMostLikedPosts(ctx context.Context) ([]models.PostUnit, error) {
	//Получает с постгреса структуры формата Name, {imageDesc, imageSK}, likes и затем
	//анмаршалит  {imageDesc, imageSK} в map[string]string
	type postsDBStruct struct {
		Name   string          `json:"name"`
		Images json.RawMessage `json:"images"`
		Likes  int             `json:"likes_count"`
	}
	var postsDB []postsDBStruct
	err := pr.db.Model(&models.Post{}).WithContext(ctx).
		Raw(`
						SELECT 
				posts.name,
				COALESCE(likes_count, 0) AS likes,
				COALESCE(
					JSON_OBJECT_AGG(images.description, images.storage_key) 
					FILTER (WHERE images.id IS NOT NULL), 
					'{}'
				) AS images
			FROM posts
			LEFT JOIN (
				SELECT post_id, COUNT(*) AS likes_count
				FROM likes
				GROUP BY post_id
			) AS like_counts ON like_counts.post_id = posts.id
			LEFT JOIN post_images ON post_images.post_id = posts.id
			LEFT JOIN images ON images.id = post_images.image_id
			GROUP BY posts.name, likes_count
			ORDER BY likes_count DESC
			LIMIT 3;
	   `).Scan(&postsDB).Error
	if err != nil {
		return nil, err
	}

	var result []models.PostUnit

	for _, post := range postsDB {
		images := make(map[string]string)
		err := json.Unmarshal(post.Images, &images)
		if err != nil {
			continue
		}
		result = append(result, models.PostUnit{Name: post.Name, Images: images, Likes: post.Likes})
	}

	return result, nil
}

func (pr *PostRepository) GetPostOwner(ctx context.Context, postID int) (int, error) {
	var ownerID int
	err := pr.db.WithContext(ctx).Model(&models.Post{}).Where("id = ?", postID).Pluck("user_id", &ownerID).Error
	if err != nil {
		return 0, err
	}
	return ownerID, nil
}
