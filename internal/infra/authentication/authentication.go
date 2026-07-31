package authentication

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gianpaoloaranha/go-social-network/internal/infra/config"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

type Claims struct {
	UserID    string `json:"userId"`
	TokenType string `json:"tokenType"`
	jwt.RegisteredClaims
}

type contextKey string

const userIDContextKey contextKey = "authenticatedUserID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		tokenValue, ok := strings.CutPrefix(authorizationHeader, "Bearer ")
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		claims, err := ValidateAccessToken(strings.TrimSpace(tokenValue))
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok && userID != ""
}

func CreateAccessToken(userID string) (string, error) {
	claims := Claims{
		UserID:    userID,
		TokenType: AccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWTSecretKey))
}

func CreateRefreshToken(userID string) (string, error) {
	claims := Claims{
		UserID:    userID,
		TokenType: RefreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.JWTSecretKey))
}

func ValidateAccessToken(accessToken string) (*Claims, error) {
	claims, err := parseToken(accessToken)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != AccessTokenType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func ValidateRefreshToken(refreshToken string) (*Claims, error) {
	claims, err := parseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != RefreshTokenType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func parseToken(tokenValue string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenValue, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(config.JWTSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
