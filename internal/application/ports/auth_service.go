package ports

import "crabi-test/internal/domain"

// AuthService define las operaciones de autenticación
type AuthService interface {
	Login(email, password string) (*domain.User, string, error)
	ValidateToken(token string) (*domain.User, error)
	GenerateToken(user *domain.User) (string, error)
}
