package auth

import "github.com/gianpaoloaranha/go-social-network/internal/app/domain"

type Usecase interface {
	Login(email, password string) (*domain.Session, error)
	RefreshToken(refreshToken string) (*domain.Session, error)
	Logout() bool
}
