package user

import "github.com/gianpaoloaranha/go-social-network/internal/app/domain"

type Repository interface {
	CreateUser(user *domain.User) (*domain.User, error)
	GetUsers() ([]domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUser(user *domain.User) (*domain.User, error)
	DeleteUser(id string) error

	FollowUser(followerID, followeeID string) error
	UnfollowUser(followerID, followeeID string) error
	GetFollowers(userID string) ([]domain.User, error)
	GetFollowing(userID string) ([]domain.User, error)
}
