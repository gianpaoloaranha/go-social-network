package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
)

func toGRPCError(err error) error {
	var domainErr domain.Error
	if !errors.As(err, &domainErr) {
		return status.Error(codes.Internal, err.Error())
	}

	code := codes.Internal
	switch domainErr.Type {
	case domain.ErrInvalidInput:
		code = codes.InvalidArgument
	case domain.ErrNotFound:
		code = codes.NotFound
	case domain.ErrConflict:
		code = codes.AlreadyExists
	case domain.ErrUnauthorized:
		code = codes.Unauthenticated
	case domain.ErrInternal:
		code = codes.Internal
	}

	return status.Error(code, domainErr.Message)
}
