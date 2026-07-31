package usecase

import (
	"github.com/gianpaoloaranha/go-social-network/internal/app/domain"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/auth"
	"github.com/gianpaoloaranha/go-social-network/internal/app/ports/user"
	"github.com/gianpaoloaranha/go-social-network/internal/infra/authentication"
)

type authUsecase struct {
	userRepository user.Repository
}

func NewAuthUsecase(userRepository user.Repository) auth.Usecase {
	return &authUsecase{
		userRepository: userRepository,
	}
}

func (uc *authUsecase) Login(email, password string) (*domain.Session, error) {
	if email == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "E-mail is required")
	}

	if password == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Password is required")
	}

	user, err := uc.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve user", err)
	}

	if user == nil {
		return nil, domain.NewError(domain.ErrInvalidInput, "Invalid e-mail or password")
	}

	if err := domain.VerifyPassword(user.Password, password); err != nil {
		return nil, domain.NewError(domain.ErrInvalidInput, "Invalid e-mail or password")
	}

	return uc.createSession(*user)
}

func (uc *authUsecase) RefreshToken(refreshToken string) (*domain.Session, error) {
	if refreshToken == "" {
		return nil, domain.NewError(domain.ErrInvalidInput, "Refresh token is required")
	}

	claims, err := authentication.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, domain.NewError(domain.ErrInvalidInput, "Invalid refresh token")
	}

	user, err := uc.userRepository.GetUserByID(claims.UserID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not retrieve user", err)
	}

	if user == nil {
		return nil, domain.NewError(domain.ErrInvalidInput, "Invalid refresh token")
	}

	return uc.createSession(*user)
}

func (uc *authUsecase) Logout() bool {
	return true
}

func (uc *authUsecase) createSession(user domain.User) (*domain.Session, error) {
	accessToken, err := authentication.CreateAccessToken(user.ID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not create access token", err)
	}

	refreshToken, err := authentication.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, domain.WrapError(domain.ErrInternal, "Could not create refresh token", err)
	}

	return &domain.Session{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
