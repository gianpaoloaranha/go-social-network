package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/comment"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/post"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/pubsub"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/user"
)

type commentUsecase struct {
	commentRepository comment.Repository
	postRepository    post.Repository
	userRepository    user.Repository
	publisher         pubsub.Publisher
	subscriber        pubsub.Subscriber
}

func NewCommentUsecase(
	commentRepository comment.Repository,
	postRepository post.Repository,
	userRepository user.Repository,
	publisher pubsub.Publisher,
	subscriber pubsub.Subscriber,
) comment.UseCase {
	return &commentUsecase{
		commentRepository: commentRepository,
		postRepository:    postRepository,
		userRepository:    userRepository,
		publisher:         publisher,
		subscriber:        subscriber,
	}
}

func (uc *commentUsecase) CreateComment(ctx context.Context, input comment.CreateCommentInput) (*domain.Comment, error) {
	if input.Message == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Message is required")
	}

	if input.AuthorID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "AuthorID is required")
	}

	if input.PostID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "PostID is required")
	}

	author, err := uc.userRepository.GetUserByID(input.AuthorID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve author", err)
	}

	if author == nil {
		return nil, domain.NewError(domain.ErrNotFound, "Author not found")
	}

	post, err := uc.postRepository.GetPostByID(input.PostID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve post", err)
	}

	if post == nil {
		return nil, domain.NewError(domain.ErrNotFound, "Post not found")
	}

	createdComment, err := uc.commentRepository.CreateComment(&domain.Comment{
		Message:  input.Message,
		PostID:   input.PostID,
		AuthorID: input.AuthorID,
	})

	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not create comment", err)
	}

	payload, err := json.Marshal(createdComment)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not marshal added comment event", err)
	}

	if err := uc.publisher.Publish(ctx, fmt.Sprintf("comment:added:%s", input.PostID), payload); err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not publish added comment event", err)
	}

	return createdComment, nil
}

func (uc *commentUsecase) GetCommentsByPostID(postID string) ([]domain.Comment, error) {
	if postID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Post ID is required")
	}

	post, err := uc.postRepository.GetPostByID(postID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve post", err)
	}

	if post == nil {
		return nil, domain.NewError(domain.ErrNotFound, "Post not found")
	}

	comments, err := uc.commentRepository.GetCommentsByPostID(postID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve comments", err)
	}

	return comments, nil
}

func (uc *commentUsecase) GetCommentByID(id string) (*domain.Comment, error) {
	if id == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Comment ID is required")
	}

	comment, err := uc.commentRepository.GetCommentByID(id)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve comment", err)
	}

	if comment == nil {
		return nil, domain.NewError(domain.ErrNotFound, "Comment not found")
	}

	return comment, nil
}

func (uc *commentUsecase) UpdateComment(input comment.UpdateCommentInput) (*domain.Comment, error) {
	if input.ID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Comment ID is required")
	}

	if input.Message == nil {
		return nil, domain.NewError(domain.ErrInvalidInput, "Message is required")
	}

	if *input.Message == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Message is required")
	}

	commentToUpdate, err := uc.commentRepository.GetCommentByID(input.ID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve comment", err)
	}

	if commentToUpdate == nil {
		return nil, domain.NewError(domain.ErrNotFound, "Comment not found")
	}

	if *input.Message != commentToUpdate.Message {
		commentToUpdate.Message = *input.Message
	}

	updatedComment, err := uc.commentRepository.UpdateComment(commentToUpdate)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not update comment", err)
	}

	return updatedComment, nil
}

func (uc *commentUsecase) DeleteComment(id string) error {
	if id == "" {
		return domain.NewError(domain.ErrInvalidInput, "Comment ID is required")
	}

	comment, err := uc.commentRepository.GetCommentByID(id)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not retrieve comment", err)
	}

	if comment == nil {
		return domain.NewError(domain.ErrNotFound, "Comment not found")
	}

	err = uc.commentRepository.DeleteComment(id)
	if err != nil {
		return domain.WrapError(domain.ErrInternal, "Could not delete comment", err)
	}

	return nil
}

func (uc *commentUsecase) SubscribeAddedComment(ctx context.Context, postID string) (<-chan *domain.Comment, error) {
	if postID == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Post ID is required")
	}

	if uc.subscriber == nil {
		return nil, fmt.Errorf("comment subscription subscriber is not configured")
	}

	raw, cleanup, err := uc.subscriber.Subscribe(ctx, fmt.Sprintf("comment:added:%s", postID))
	if err != nil {
		return nil, err
	}

	out := make(chan *domain.Comment)
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
				var c domain.Comment
				if err := json.Unmarshal(msg, &c); err != nil {
					continue
				}
				select {
				case <-ctx.Done():
					return
				case out <- &c:
				}
			}
		}
	}()

	return out, nil
}
