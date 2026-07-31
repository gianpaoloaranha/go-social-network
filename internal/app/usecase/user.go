package usecase

import (
	"net/mail"

	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/user"
)

type userUsecase struct {
	userRepository user.Repository
}

func NewUserUsecase(userRepository user.Repository) user.UseCase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

// CreateUser creates a new user in the system.
func (uc *userUsecase) CreateUser(input user.CreateUserInput) (*domain.User, error) {
	if input.Name == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Name is required")
	}

	if input.Email == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Email is required")
	}

	if _, err := mail.ParseAddress(input.Email); err != nil {
		return nil, domain.NewError(domain.ErrInvalidInput, "Invalid email format")
	}

	if input.Password == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Password is required")
	}

	hashedPassword, err := domain.Hash(input.Password)
	if err != nil {
		return nil, domain.NewError(domain.ErrInternal, "Failed to hash password")
	}

	createdUser, err := uc.userRepository.CreateUser(&domain.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not create user", err)
	}

	return createdUser, nil
}

// GetUsers retrieves all users from the system.
func (uc *userUsecase) GetUsers() ([]domain.User, error) {
	users, err := uc.userRepository.GetUsers()
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve users", err)
	}

	return users, nil
}

// GetUserByID retrieves a user by their ID.
func (uc *userUsecase) GetUserByID(id string) (*domain.User, error) {
	if id == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "User ID is required")
	}

	user, err := uc.userRepository.GetUserByID(id)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve user", err)
	}

	if user == nil {
		return nil, domain.NewError(domain.ErrNotFound, "User not found")
	}

	return user, nil
}

// UpdateUser updates an existing user's information.
func (uc *userUsecase) UpdateUser(input user.UpdateUserInput) (*domain.User, error) {
	if input.ID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "User ID is required")
	}

	if input.Name == nil && input.Email == nil && input.Password == nil {
		return nil, domain.NewError(domain.ErrInvalidInput, "At least one field is required")
	}

	if input.Name != nil && *input.Name == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Name is required")
	}

	if input.Email != nil && *input.Email == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Email is required")
	}

	if input.Password != nil && *input.Password == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Password is required")
	}

	userToUpdate, err := uc.userRepository.GetUserByID(input.ID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve user", err)
	}

	if userToUpdate == nil {
		return nil, domain.NewError(domain.ErrNotFound, "User not found")
	}

	if input.Name != nil {
		userToUpdate.Name = *input.Name
	}

	if input.Email != nil {
		if _, err := mail.ParseAddress(*input.Email); err != nil {
			return nil, domain.NewError(domain.ErrInvalidInput, "Invalid email format")
		}

		userToUpdate.Email = *input.Email
	}

	if input.Password != nil {
		hashedPassword, err := domain.Hash(*input.Password)
		if err != nil {
			return nil, domain.NewError(domain.ErrInternal, "Failed to hash password")
		}

		userToUpdate.Password = string(hashedPassword)
	}

	updatedUser, err := uc.userRepository.UpdateUser(userToUpdate)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not update user", err)
	}

	return updatedUser, nil
}

// DeleteUser deletes a user from the system by their ID.
func (uc *userUsecase) DeleteUser(id string) error {
	if id == "" {
		return domain.NewError(domain.ErrInvalidInput, "User ID is required")
	}

	user, err := uc.userRepository.GetUserByID(id)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not retrieve user", err)
	}

	if user == nil {
		return domain.NewError(domain.ErrNotFound, "User not found")
	}

	err = uc.userRepository.DeleteUser(id)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not delete user", err)
	}

	return nil
}

// FollowUser allows a user to follow another user.
func (uc *userUsecase) FollowUser(followerID, followeeID string) error {
	if followerID == "" {
		return domain.NewError(domain.ErrInvalidInput, "Follower ID is required")
	}

	if followeeID == "" {
		return domain.NewError(domain.ErrInvalidInput, "Followee ID is required")
	}

	if followerID == followeeID {
		return domain.NewError(domain.ErrInvalidInput, "User cannot follow themselves")
	}

	follower, err := uc.userRepository.GetUserByID(followerID)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not retrieve follower", err)
	}

	if follower == nil {
		return domain.NewError(domain.ErrNotFound, "Follower not found")
	}

	followee, err := uc.userRepository.GetUserByID(followeeID)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not retrieve followee", err)
	}

	if followee == nil {
		return domain.NewError(domain.ErrNotFound, "Followee not found")
	}

	err = uc.userRepository.FollowUser(followerID, followeeID)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not follow user", err)
	}

	return nil
}

// UnfollowUser allows a user to unfollow another user.
func (uc *userUsecase) UnfollowUser(followerID, followeeID string) error {
	if followerID == "" {
		return domain.NewError(domain.ErrInvalidInput, "Follower ID is required")
	}

	if followeeID == "" {
		return domain.NewError(domain.ErrInvalidInput, "Followee ID is required")
	}

	if followerID == followeeID {
		return domain.NewError(domain.ErrInvalidInput, "User cannot unfollow themselves")
	}

	follower, err := uc.userRepository.GetUserByID(followerID)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not retrieve follower", err)
	}

	if follower == nil {
		return domain.NewError(domain.ErrNotFound, "Follower not found")
	}

	followee, err := uc.userRepository.GetUserByID(followeeID)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not retrieve followee", err)
	}

	if followee == nil {
		return domain.NewError(domain.ErrNotFound, "Followee not found")
	}

	err = uc.userRepository.UnfollowUser(followerID, followeeID)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not unfollow user", err)
	}

	return nil
}

// GetFollowers retrieves the list of followers for a given user.
func (uc *userUsecase) GetFollowers(userID string) ([]domain.User, error) {
	if userID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "User ID is required")
	}

	user, err := uc.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve user", err)
	}

	if user == nil {
		return nil, domain.NewError(domain.ErrNotFound, "User not found")
	}

	followers, err := uc.userRepository.GetFollowers(userID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve followers", err)
	}

	return followers, nil
}

// GetFollowing retrieves the list of users that a given user is following.
func (uc *userUsecase) GetFollowing(userID string) ([]domain.User, error) {
	if userID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "User ID is required")
	}

	user, err := uc.userRepository.GetUserByID(userID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve user", err)
	}

	if user == nil {
		return nil, domain.NewError(domain.ErrNotFound, "User not found")
	}

	following, err := uc.userRepository.GetFollowing(userID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve following", err)
	}

	return following, nil
}
