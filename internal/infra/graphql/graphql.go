package graphql

import (
	"gorm.io/gorm"

	"github.com/gianpaoloaranha/go-social-network/internal/adapters/in/graphql/resolver"
	"github.com/gianpaoloaranha/go-social-network/internal/adapters/out/db/postgres"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/pubsub"
	"github.com/gianpaoloaranha/go-social-network/internal/app/usecase"
)

// BuildResolvers initializes the resolvers with the necessary use cases and repositories.
func BuildResolvers(db *gorm.DB, publisher pubsub.Publisher, subscriber pubsub.Subscriber) *resolver.Resolver {
	userRepository := postgres.NewUserRepository(db)
	postRepository := postgres.NewPostRepository(db)
	commentRepository := postgres.NewCommentRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository)
	postUsecase := usecase.NewPostUsecase(postRepository, userRepository, publisher, subscriber)
	commentUsecase := usecase.NewCommentUsecase(commentRepository, postRepository, userRepository, publisher, subscriber)

	return &resolver.Resolver{
		UserUsecase:    userUsecase,
		PostUsecase:    postUsecase,
		CommentUsecase: commentUsecase,
	}
}
