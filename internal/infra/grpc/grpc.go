package grpc

import (
	googlegrpc "google.golang.org/grpc"
	"gorm.io/gorm"

	grpcadapter "github.com/gianpaoloaranha/go-social-network/internal/adapters/in/grpc"
	"github.com/gianpaoloaranha/go-social-network/internal/adapters/out/db/postgres"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/pubsub"
	"github.com/gianpaoloaranha/go-social-network/internal/app/usecase"
	socialnetwork "github.com/gianpaoloaranha/go-social-network/proto/gen"
)

// RegisterServices wires repositories and use cases into the gRPC service
// implementations and registers them on the given server.
func RegisterServices(server *googlegrpc.Server, db *gorm.DB, publisher pubsub.Publisher, subscriber pubsub.Subscriber) {
	userRepository := postgres.NewUserRepository(db)
	postRepository := postgres.NewPostRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository)
	postUsecase := usecase.NewPostUsecase(postRepository, userRepository, publisher, subscriber)

	socialnetwork.RegisterUserServiceServer(server, grpcadapter.NewGrpcUserService(userUsecase))
	socialnetwork.RegisterPostServiceServer(server, grpcadapter.NewGrpcPostService(postUsecase))
}
