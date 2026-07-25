package resolver

import (
	"github.com/gianpaoloaranha/go-social-network/internal/adapters/in/graphql/generated/model"
	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
)

func userToGraphQL(user domain.User) *model.User {
	return &model.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Following: []*model.User{},
		Followers: []*model.User{},
		Posts:     []*model.Post{},
	}
}

func usersToGraphQL(users []domain.User) []*model.User {
	result := make([]*model.User, 0, len(users))
	for _, user := range users {
		result = append(result, userToGraphQL(user))
	}

	return result
}

func postToGraphQL(post domain.Post) *model.Post {
	return &model.Post{
		ID:          post.ID,
		Description: post.Description,
		CreatedAt:   post.CreatedAt,
		Comments:    []*model.Comment{},
	}
}

func postsToGraphQL(posts []domain.Post) []*model.Post {
	result := make([]*model.Post, 0, len(posts))
	for _, post := range posts {
		result = append(result, postToGraphQL(post))
	}

	return result
}

func commentToGraphQL(comment domain.Comment) *model.Comment {
	return &model.Comment{
		ID:        comment.ID,
		Message:   comment.Message,
		CreatedAt: comment.CreatedAt,
	}
}

func commentsToGraphQL(comments []domain.Comment) []*model.Comment {
	result := make([]*model.Comment, 0, len(comments))
	for _, comment := range comments {
		result = append(result, commentToGraphQL(comment))
	}

	return result
}
