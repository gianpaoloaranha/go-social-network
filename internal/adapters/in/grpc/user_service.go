package grpc

import (
	"context"

	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/user"
	socialnetwork "github.com/gianpaoloaranha/go-social-network/proto/gen"
)

type grpcUserService struct {
	socialnetwork.UnimplementedUserServiceServer

	userUsecase user.UseCase
}

func NewGrpcUserService(userUsecase user.UseCase) *grpcUserService {
	return &grpcUserService{
		userUsecase: userUsecase,
	}
}

func (s *grpcUserService) CreateUser(ctx context.Context, req *socialnetwork.CreateUserRequest) (*socialnetwork.CreateUserResponse, error) {
	userInput := createUserRequestToUserInput(req)
	createdUser, err := s.userUsecase.CreateUser(*userInput)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return userToCreateUserResponse(createdUser), nil
}
