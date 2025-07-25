package services

import (
	"crabi-test/internal/application/ports"
	"crabi-test/internal/domain"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserService implementa la lógica de negocio para usuarios
type UserService struct {
	userRepo   ports.UserRepository
	pldService ports.PLDService
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService(userRepo ports.UserRepository, pldService ports.PLDService) *UserService {
	return &UserService{
		userRepo:   userRepo,
		pldService: pldService,
	}
}

// CreateUser crea un nuevo usuario validando contra el servicio PLD
func (s *UserService) CreateUser(user *domain.User) error {
	// Validar que el email no exista
	existingUser, err := s.userRepo.GetByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("el email ya está registrado")
	}

	// Validar contra el servicio PLD
	pldResponse, err := s.pldService.ValidateUser(user.IDNumber, user.Name, user.Email)
	if err != nil {
		return errors.New("error validando usuario con servicio PLD")
	}

	if pldResponse.IsBlacklisted {
		return errors.New("usuario en lista negra: " + pldResponse.Reason)
	}

	// Encriptar contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error encriptando contraseña")
	}
	user.Password = string(hashedPassword)

	// Establecer timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Guardar en base de datos
	return s.userRepo.Create(user)
}

// GetUser obtiene un usuario por ID
func (s *UserService) GetUser(id uint) (*domain.User, error) {
	return s.userRepo.GetByID(id)
}

// GetUserByEmail obtiene un usuario por email
func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.userRepo.GetByEmail(email)
}

// UpdateUser actualiza un usuario
func (s *UserService) UpdateUser(user *domain.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(user)
}

// DeleteUser elimina un usuario
func (s *UserService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}
