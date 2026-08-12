package grpc

import (
	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/post"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/user"
	socialnetwork "github.com/gianpaoloaranha/go-social-network/proto/gen"
)

func createUserRequestToUserInput(req *socialnetwork.CreateUserRequest) *user.CreateUserInput {
	return &user.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}

func userToCreateUserResponse(user *domain.User) *socialnetwork.CreateUserResponse {
	return &socialnetwork.CreateUserResponse{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
	}
}

func createPostRequestToPostInput(req *socialnetwork.CreatePostRequest) *post.CreatePostInput {
	return &post.CreatePostInput{
		AuthorID:    req.AuthorId,
		Description: req.Description,
	}
}

func postToCreatePostResponse(p *domain.Post) *socialnetwork.CreatePostResponse {
	return &socialnetwork.CreatePostResponse{
		PostId:      p.ID,
		AuthorId:    p.AuthorID,
		Description: p.Description,
	}
}
