package resolver

import "github.com/gianpaoloaranha/go-social-network/internal/adapters/in/graphql/generated/model"

func graphQLError(field string, err error) []*model.Error {
	return []*model.Error{
		{
			Field:   field,
			Message: err.Error(),
		},
	}
}

func userPayloadError(field string, err error) *model.UserPayload {
	return &model.UserPayload{
		Errors: graphQLError(field, err),
	}
}

func postPayloadError(field string, err error) *model.PostPayload {
	return &model.PostPayload{
		Errors: graphQLError(field, err),
	}
}

func commentPayloadError(field string, err error) *model.CommentPayload {
	return &model.CommentPayload{
		Errors: graphQLError(field, err),
	}
}

func deletePayloadError(field string, err error) *model.DeletePayload {
	return &model.DeletePayload{
		Status: false,
		Errors: graphQLError(field, err),
	}
}
