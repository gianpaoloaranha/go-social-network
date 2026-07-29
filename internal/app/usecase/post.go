package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/post"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/pubsub"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/user"
)

type postUsecase struct {
	postRepository post.Repository
	userRepository user.Repository
	publisher      pubsub.Publisher
	subscriber     pubsub.Subscriber
}

func NewPostUsecase(
	postRepository post.Repository,
	userRepository user.Repository,
	publisher pubsub.Publisher,
	subscriber pubsub.Subscriber,
) post.UseCase {
	return &postUsecase{
		postRepository: postRepository,
		userRepository: userRepository,
		publisher:      publisher,
		subscriber:     subscriber,
	}
}

func (uc *postUsecase) CreatePost(ctx context.Context, post post.CreatePostInput) (*domain.Post, error) {
	if post.Description == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Description is required")
	}

	if post.AuthorID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "AuthorID is required")
	}

	author, err := uc.userRepository.GetUserByID(post.AuthorID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve author", err)
	}

	if author == nil {
		return nil, domain.NewError(domain.ErrNotFound, "Author not found")
	}

	createdPost, err := uc.postRepository.CreatePost(&domain.Post{
		Description: post.Description,
		AuthorID:    post.AuthorID,
	})

	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not create post", err)
	}

	payload, err := json.Marshal(createdPost)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not marshal created post event", err)
	}

	if err := uc.publisher.Publish(ctx, "post:created", payload); err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not publish created post event", err)
	}

	return createdPost, nil
}

func (uc *postUsecase) GetPosts() ([]domain.Post, error) {
	posts, err := uc.postRepository.GetPosts()
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve posts", err)
	}

	return posts, nil
}

func (uc *postUsecase) GetPostsByAuthorID(authorID string) ([]domain.Post, error) {
	if authorID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "AuthorID is required")
	}

	author, err := uc.userRepository.GetUserByID(authorID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve author", err)
	}

	if author == nil {
		return nil, domain.NewError(domain.ErrNotFound, "Author not found")
	}

	posts, err := uc.postRepository.GetPostsByAuthorID(authorID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve posts", err)
	}

	return posts, nil
}

func (uc *postUsecase) GetPostByID(id string) (*domain.Post, error) {
	if id == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Post ID is required")
	}

	post, err := uc.postRepository.GetPostByID(id)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve post", err)
	}

	if post == nil {
		return nil, domain.NewError(domain.ErrNotFound, "Post not found")
	}

	return post, nil
}

func (uc *postUsecase) UpdatePost(post post.UpdatePostInput) (*domain.Post, error) {
	if post.ID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Post ID is required")
	}

	if post.Description == nil && post.AuthorID == nil {
		return nil, domain.NewError(domain.ErrInvalidInput, "At least one field is required")
	}

	if post.Description != nil && *post.Description == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Description is required")
	}

	if post.AuthorID != nil && *post.AuthorID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "AuthorID is required")
	}

	postToUpdate, err := uc.postRepository.GetPostByID(post.ID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve post", err)
	}

	if postToUpdate == nil {
		return nil, domain.NewError(domain.ErrNotFound, "Post not found")
	}

	if post.Description != nil && *post.Description != postToUpdate.Description {
		postToUpdate.Description = *post.Description
	}

	if post.AuthorID != nil {
		author, err := uc.userRepository.GetUserByID(*post.AuthorID)
		if err != nil {
			return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve author", err)
		}

		if author == nil {
			return nil, domain.NewError(domain.ErrNotFound, "Author not found")
		}

		postToUpdate.AuthorID = *post.AuthorID
	}

	updatedPost, err := uc.postRepository.UpdatePost(postToUpdate)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not update post", err)
	}

	return updatedPost, nil
}

func (uc *postUsecase) DeletePost(id string) error {
	if id == "" {
		return domain.NewError(domain.ErrInvalidInput, "Post ID is required")
	}

	post, err := uc.postRepository.GetPostByID(id)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not retrieve post", err)
	}

	if post == nil {
		return domain.NewError(domain.ErrNotFound, "Post not found")
	}

	err = uc.postRepository.DeletePost(id)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not delete post", err)
	}

	return nil
}

func (uc *postUsecase) SubscribeCreatedPost(ctx context.Context) (<-chan *domain.Post, error) {
	if uc.subscriber == nil {
		return nil, fmt.Errorf("post subscription subscriber is not configured")
	}

	raw, cleanup, err := uc.subscriber.Subscribe(ctx, "post:created")
	if err != nil {
		return nil, err
	}

	out := make(chan *domain.Post)
	go func() {
		defer close(out)
		defer cleanup()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-raw:
				if !ok {
					return
				}
				var p domain.Post
				if err := json.Unmarshal(msg, &p); err != nil {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case out <- &p:
				}
			}
		}
	}()

	return out, nil
}
