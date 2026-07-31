package domain

import "golang.org/x/crypto/bcrypt"

type Session struct {
	AccessToken  string
	RefreshToken string
	User         User
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}