package user

import "github.com/gianpaoloaranha/go-social-network/internal/app/domain"

type CreateUserInput struct {
	Name     string
	Email    string
	Password string
}

type UpdateUserInput struct {
	ID       string
	Name     *string
	Email    *string
	Password *string
}

type UseCase interface {
	CreateUser(user CreateUserInput) (*domain.User, error)
	GetUsers() ([]domain.User, error)
	GetUserByID(id string) (*domain.User, error)
	UpdateUser(user UpdateUserInput) (*domain.User, error)
	DeleteUser(id string) error

	FollowUser(followerID, followeeID string) error
	UnfollowUser(followerID, followeeID string) error
	GetFollowers(userID string) ([]domain.User, error)
	GetFollowing(userID string) ([]domain.User, error)
}
