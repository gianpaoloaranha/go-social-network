package grpc

import (
	"context"

	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/post"
	socialnetwork "github.com/gianpaoloaranha/go-social-network/proto/gen"
)

type grpcPostService struct {
	socialnetwork.UnimplementedPostServiceServer

	postUsecase post.UseCase
}

func NewGrpcPostService(postUsecase post.UseCase) *grpcPostService {
	return &grpcPostService{
		postUsecase: postUsecase,
	}
}

func (s *grpcPostService) CreatePost(ctx context.Context, req *socialnetwork.CreatePostRequest) (*socialnetwork.CreatePostResponse, error) {
	panic("implement me")
}