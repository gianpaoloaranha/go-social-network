package comment

import (
	"context"

	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
)

type CreateCommentInput struct {
	AuthorID string
	PostID   string
	Message  string
}

type UpdateCommentInput struct {
	ID      string
	Message *string
}

type UseCase interface {
	CreateComment(ctx context.Context, comment CreateCommentInput) (*domain.Comment, error)
	GetCommentByID(id string) (*domain.Comment, error)
	GetCommentsByPostID(postID string) ([]domain.Comment, error)
	UpdateComment(comment UpdateCommentInput) (*domain.Comment, error)
	DeleteComment(id string) error
	SubscribeAddedComment(ctx context.Context, postID string) (<-chan *domain.Comment, error)
}
