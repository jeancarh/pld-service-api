package tests

import (
	"crabi-test/internal/application/services"
	"crabi-test/internal/domain"
	"testing"
)

// MockUserRepository implementa un repositorio mock para testing
type MockUserRepository struct {
	users  map[uint]*domain.User
	emails map[string]*domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:  make(map[uint]*domain.User),
		emails: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) Create(user *domain.User) error {
	user.ID = uint(len(m.users) + 1)
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *MockUserRepository) GetByID(id uint) (*domain.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	if user, exists := m.emails[email]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *MockUserRepository) Update(user *domain.User) error {
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *MockUserRepository) Delete(id uint) error {
	if user, exists := m.users[id]; exists {
		delete(m.users, id)
		delete(m.emails, user.Email)
	}
	return nil
}

// MockPLDService implementa un servicio PLD mock para testing
type MockPLDService struct {
	shouldBlacklist bool
}

func NewMockPLDService(shouldBlacklist bool) *MockPLDService {
	return &MockPLDService{
		shouldBlacklist: shouldBlacklist,
	}
}

func (m *MockPLDService) ValidateUser(idNumber, name, email string) (*domain.PLDResponse, error) {
	if m.shouldBlacklist {
		return &domain.PLDResponse{
			IsBlacklisted: true,
			Status:        "blacklisted",
			Reason:        "Usuario en lista negra por actividades sospechosas",
		}, nil
	}

	return &domain.PLDResponse{
		IsBlacklisted: false,
		Status:        "clean",
		Reason:        "",
	}, nil
}

func TestUserService_CreateUser_Success(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := services.NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	// Act
	err := userService.CreateUser(user)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}

	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}

	if user.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}
}

func TestUserService_CreateUser_Blacklisted(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(true)
	userService := services.NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	// Act
	err := userService.CreateUser(user)

	// Assert
	if err == nil {
		t.Error("Expected error for blacklisted user")
	}

	if err.Error() != "usuario en lista negra: Usuario en lista negra por actividades sospechosas" {
		t.Errorf("Expected blacklist error, got %v", err)
	}
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := services.NewUserService(userRepo, pldService)

	// Crear usuario existente
	existingUser := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}
	userRepo.Create(existingUser)

	// Intentar crear usuario con mismo email
	newUser := &domain.User{
		Name:     "Otro Usuario",
		Email:    "juan.perez@email.com", // Mismo email
		Password: "password456",
		IDNumber: "87654321",
	}

	// Act
	err := userService.CreateUser(newUser)

	// Assert
	if err == nil {
		t.Error("Expected error for duplicate email")
	}

	if err.Error() != "el email ya está registrado" {
		t.Errorf("Expected duplicate email error, got %v", err)
	}
}
