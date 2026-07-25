package postgres

import (
	"errors"
	"time"

	"github.com/gianpaoloaranha/go-social-network/internal/adapters/out/db/postgres/model"
	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) CreateComment(comment *domain.Comment) (*domain.Comment, error) {
	commentModel := commentFromDomain(*comment)
	if commentModel.ID == "" {
		commentModel.ID = uuid.NewString()
	}
	if commentModel.CreatedAt.IsZero() {
		commentModel.CreatedAt = time.Now().UTC()
	}

	if err := r.db.Create(&commentModel).Error; err != nil {
		return nil, err
	}

	createdComment := commentToDomain(commentModel)
	return &createdComment, nil
}

func (r *CommentRepository) GetCommentByID(id string) (*domain.Comment, error) {
	var comment model.Comment
	if err := r.db.First(&comment, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	domainComment := commentToDomain(comment)
	return &domainComment, nil
}

func (r *CommentRepository) GetCommentsByPostID(postID string) ([]domain.Comment, error) {
	var comments []model.Comment
	if err := r.db.
		Where("post_id = ?", postID).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		return nil, err
	}

	return commentsToDomain(comments), nil
}

func (r *CommentRepository) UpdateComment(comment *domain.Comment) (*domain.Comment, error) {
	commentModel := commentFromDomain(*comment)
	if err := r.db.Model(&model.Comment{}).
		Where("id = ?", commentModel.ID).
		Updates(map[string]any{
			"message": commentModel.Message,
		}).Error; err != nil {
		return nil, err
	}

	return r.GetCommentByID(commentModel.ID)
}

func (r *CommentRepository) DeleteComment(id string) error {
	return r.db.Delete(&model.Comment{}, "id = ?", id).Error
}
