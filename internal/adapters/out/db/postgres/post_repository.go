package postgres

import (
	"errors"
	"time"

	"github.com/gianpaoloaranha/go-social-network/internal/adapters/out/db/postgres/model"
	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(post *domain.Post) (*domain.Post, error) {
	postModel := postFromDomain(*post)
	if postModel.ID == "" {
		postModel.ID = uuid.NewString()
	}
	if postModel.CreatedAt.IsZero() {
		postModel.CreatedAt = time.Now().UTC()
	}

	if err := r.db.Create(&postModel).Error; err != nil {
		return nil, err
	}

	createdPost := postToDomain(postModel)
	return &createdPost, nil
}

func (r *PostRepository) GetPosts() ([]domain.Post, error) {
	var posts []model.Post
	if err := r.db.Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}

	return postsToDomain(posts), nil
}

func (r *PostRepository) GetPostByID(id string) (*domain.Post, error) {
	var post model.Post
	if err := r.db.First(&post, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	domainPost := postToDomain(post)
	return &domainPost, nil
}

func (r *PostRepository) UpdatePost(post *domain.Post) (*domain.Post, error) {
	postModel := postFromDomain(*post)
	if err := r.db.Model(&model.Post{}).
		Where("id = ?", postModel.ID).
		Updates(map[string]any{
			"author_id":   postModel.AuthorID,
			"description": postModel.Description,
		}).Error; err != nil {
		return nil, err
	}

	return r.GetPostByID(postModel.ID)
}

func (r *PostRepository) DeletePost(id string) error {
	return r.db.Delete(&model.Post{}, "id = ?", id).Error
}
