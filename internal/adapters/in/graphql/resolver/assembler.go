package resolver

import (
	"github.com/gianpaoloaranha/go-social-network/internal/adapters/in/graphql/generated/model"
	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
)

func (r *Resolver) assembleUser(user domain.User) (*model.User, error) {
	result := userToGraphQL(user)

	following, err := r.UserUsecase.GetFollowing(user.ID)
	if err != nil {
		return nil, err
	}
	result.Following = usersToGraphQL(following)

	followers, err := r.UserUsecase.GetFollowers(user.ID)
	if err != nil {
		return nil, err
	}
	result.Followers = usersToGraphQL(followers)

	posts, err := r.PostUsecase.GetPostsByAuthorID(user.ID)
	if err != nil {
		return nil, err
	}
	result.Posts = make([]*model.Post, 0, len(posts))
	for _, post := range posts {
		assembledPost, err := r.assemblePost(post)
		if err != nil {
			return nil, err
		}

		result.Posts = append(result.Posts, assembledPost)
	}

	return result, nil
}

func (r *Resolver) assemblePost(post domain.Post) (*model.Post, error) {
	result := postToGraphQL(post)

	author, err := r.UserUsecase.GetUserByID(post.AuthorID)
	if err != nil {
		return nil, err
	}
	result.Author = userToGraphQL(*author)

	comments, err := r.CommentUsecase.GetCommentsByPostID(post.ID)
	if err != nil {
		return nil, err
	}
	result.Comments = commentsToGraphQL(comments)
	for index, comment := range comments {
		assembledComment, err := r.assembleComment(comment)
		if err != nil {
			return nil, err
		}

		result.Comments[index] = assembledComment
	}

	return result, nil
}

func (r *Resolver) assembleComment(comment domain.Comment) (*model.Comment, error) {
	result := commentToGraphQL(comment)

	author, err := r.UserUsecase.GetUserByID(comment.AuthorID)
	if err != nil {
		return nil, err
	}
	result.Author = userToGraphQL(*author)

	return result, nil
}
