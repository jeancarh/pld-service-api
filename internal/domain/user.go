package domain

import (
	"time"
)

// User representa la entidad de usuario en el dominio
type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // No se serializa en JSON
	IDNumber  string    `json:"id_number"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository define las operaciones de persistencia para usuarios
type UserRepository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
}

// UserService define las operaciones de negocio para usuarios
type UserService interface {
	CreateUser(user *User) error
	GetUser(id uint) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id uint) error
}
