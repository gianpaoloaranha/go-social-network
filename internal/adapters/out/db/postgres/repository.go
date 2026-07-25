package postgres

import (
	commentport "github.com/gianpaoloaranha/go-social-network/internal/app/ports/comment"
	postport "github.com/gianpaoloaranha/go-social-network/internal/app/ports/post"
	userport "github.com/gianpaoloaranha/go-social-network/internal/app/ports/user"
)

var (
	_ userport.Repository    = (*UserRepository)(nil)
	_ postport.Repository    = (*PostRepository)(nil)
	_ commentport.Repository = (*CommentRepository)(nil)
)
