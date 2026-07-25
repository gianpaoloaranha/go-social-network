package post

import "github.com/gianpaoloaranha/go-social-network/internal/app/domain"

type Repository interface {
	CreatePost(post *domain.Post) (*domain.Post, error)
	GetPosts() ([]domain.Post, error)
	GetPostsByAuthorID(authorID string) ([]domain.Post, error)
	GetPostByID(id string) (*domain.Post, error)
	UpdatePost(post *domain.Post) (*domain.Post, error)
	DeletePost(id string) error
}
