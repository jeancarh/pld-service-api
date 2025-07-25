package ports

import "crabi-test/internal/domain"

// UserRepository define las operaciones de persistencia para usuarios
type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id uint) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
}
