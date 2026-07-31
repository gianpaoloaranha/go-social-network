package postgres

import (
	"errors"

	"github.com/gianpaoloaranha/go-social-network/internal/adapters/out/db/postgres/model"
	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	userModel := userFromDomain(*user)
	if userModel.ID == "" {
		userModel.ID = uuid.NewString()
	}

	if err := r.db.Create(&userModel).Error; err != nil {
		return nil, err
	}

	createdUser := userToDomain(userModel)
	return &createdUser, nil
}

func (r *UserRepository) GetUsers() ([]domain.User, error) {
	var users []model.User
	if err := r.db.Order("name ASC").Find(&users).Error; err != nil {
		return nil, err
	}

	return usersToDomain(users), nil
}

func (r *UserRepository) GetUserByID(id string) (*domain.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	domainUser := userToDomain(user)
	return &domainUser, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user model.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	domainUser := userToDomain(user)
	return &domainUser, nil
}

func (r *UserRepository) UpdateUser(user *domain.User) (*domain.User, error) {
	userModel := userFromDomain(*user)
	if err := r.db.Model(&model.User{}).
		Where("id = ?", userModel.ID).
		Updates(map[string]any{
			"name":     userModel.Name,
			"email":    userModel.Email,
			"password": userModel.Password,
		}).Error; err != nil {
		return nil, err
	}

	return r.GetUserByID(userModel.ID)
}

func (r *UserRepository) DeleteUser(id string) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}

func (r *UserRepository) FollowUser(followerID, followeeID string) error {
	userFollow := model.UserFollow{
		FollowerID: followerID,
		FolloweeID: followeeID,
	}

	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&userFollow).Error
}

func (r *UserRepository) UnfollowUser(followerID, followeeID string) error {
	return r.db.Delete(&model.UserFollow{}, "follower_id = ? AND followee_id = ?", followerID, followeeID).Error
}

func (r *UserRepository) GetFollowers(userID string) ([]domain.User, error) {
	var users []model.User
	if err := r.db.
		Joins("JOIN user_follows ON user_follows.follower_id = users.id").
		Where("user_follows.followee_id = ?", userID).
		Order("users.name ASC").
		Find(&users).Error; err != nil {
		return nil, err
	}

	return usersToDomain(users), nil
}

func (r *UserRepository) GetFollowing(userID string) ([]domain.User, error) {
	var users []model.User
	if err := r.db.
		Joins("JOIN user_follows ON user_follows.followee_id = users.id").
		Where("user_follows.follower_id = ?", userID).
		Order("users.name ASC").
		Find(&users).Error; err != nil {
		return nil, err
	}

	return usersToDomain(users), nil
}
