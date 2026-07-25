package post

import "github.com/gianpaoloaranha/go-social-network/internal/app/domain"

type CreatePostInput struct {
	Description string
	AuthorID    string
}

type UpdatePostInput struct {
	ID          string
	Description *string
	AuthorID    *string
}

type UseCase interface {
	CreatePost(post CreatePostInput) (*domain.Post, error)
	GetPosts() ([]domain.Post, error)
	GetPostsByAuthorID(authorID string) ([]domain.Post, error)
	GetPostByID(id string) (*domain.Post, error)
	UpdatePost(post UpdatePostInput) (*domain.Post, error)
	DeletePost(id string) error
}
