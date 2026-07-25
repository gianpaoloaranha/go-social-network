package postgres

import (
	"github.com/gianpaoloaranha/go-social-network/internal/adapters/out/db/postgres/model"
	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
)

func userToDomain(user model.User) domain.User {
	return domain.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func usersToDomain(users []model.User) []domain.User {
	domainUsers := make([]domain.User, 0, len(users))
	for _, user := range users {
		domainUsers = append(domainUsers, userToDomain(user))
	}

	return domainUsers
}

func userFromDomain(user domain.User) model.User {
	return model.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func postToDomain(post model.Post) domain.Post {
	return domain.Post{
		ID:          post.ID,
		AuthorID:    post.AuthorID,
		Description: post.Description,
		CreatedAt:   post.CreatedAt,
	}
}

func postsToDomain(posts []model.Post) []domain.Post {
	domainPosts := make([]domain.Post, 0, len(posts))
	for _, post := range posts {
		domainPosts = append(domainPosts, postToDomain(post))
	}

	return domainPosts
}

func postFromDomain(post domain.Post) model.Post {
	return model.Post{
		ID:          post.ID,
		AuthorID:    post.AuthorID,
		Description: post.Description,
		CreatedAt:   post.CreatedAt,
	}
}

func commentToDomain(comment model.Comment) domain.Comment {
	return domain.Comment{
		ID:        comment.ID,
		AuthorID:  comment.AuthorID,
		PostID:    comment.PostID,
		Message:   comment.Message,
		CreatedAt: comment.CreatedAt,
	}
}

func commentsToDomain(comments []model.Comment) []domain.Comment {
	domainComments := make([]domain.Comment, 0, len(comments))
	for _, comment := range comments {
		domainComments = append(domainComments, commentToDomain(comment))
	}

	return domainComments
}

func commentFromDomain(comment domain.Comment) model.Comment {
	return model.Comment{
		ID:        comment.ID,
		AuthorID:  comment.AuthorID,
		PostID:    comment.PostID,
		Message:   comment.Message,
		CreatedAt: comment.CreatedAt,
	}
}
