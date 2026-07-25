package graphql

import (
	"gorm.io/gorm"

	"github.com/gianpaoloaranha/go-social-network/internal/adapters/in/graphql/resolver"
	"github.com/gianpaoloaranha/go-social-network/internal/adapters/out/db/postgres"
	"github.com/gianpaoloaranha/go-social-network/internal/app/usecase"
)

// BuildResolvers initializes the resolvers with the necessary use cases and repositories.
func BuildResolvers(db *gorm.DB) *resolver.Resolver {
	userRepository := postgres.NewUserRepository(db)
	postRepository := postgres.NewPostRepository(db)
	commentRepository := postgres.NewCommentRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepository)
	postUsecase := usecase.NewPostUsecase(postRepository, userRepository)
	commentUsecase := usecase.NewCommentUsecase(commentRepository, postRepository, userRepository)

	return &resolver.Resolver{
		UserUsecase:    userUsecase,
		PostUsecase:    postUsecase,
		CommentUsecase: commentUsecase,
	}
}
