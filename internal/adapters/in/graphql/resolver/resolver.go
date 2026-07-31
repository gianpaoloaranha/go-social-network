package resolver

import (
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/auth"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/comment"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/post"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	UserUsecase    user.UseCase
	PostUsecase    post.UseCase
	CommentUsecase comment.UseCase
	AuthUsecase    auth.Usecase
}
