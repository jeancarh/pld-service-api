package services

import (
	"crabi-test/internal/domain"
	"fmt"
	"testing"
	"time"
)

// MockUserRepository para testing
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

// MockPLDService para testing
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
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
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
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(true)
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err == nil {
		t.Error("Expected error for blacklisted user")
	}

	if err.Error() != "usuario en lista negra: Usuario en lista negra por actividades sospechosas" {
		t.Errorf("Expected blacklist error, got %v", err)
	}
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

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

	err := userService.CreateUser(newUser)
	if err == nil {
		t.Error("Expected error for duplicate email")
	}

	if err.Error() != "el email ya está registrado" {
		t.Errorf("Expected duplicate email error, got %v", err)
	}
}

func TestUserService_GetUser(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Crear usuario
	user := &domain.User{
		Name:      "Juan Pérez",
		Email:     "juan.perez@email.com",
		Password:  "password123",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Test GetUser
	retrievedUser, err := userService.GetUser(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrievedUser == nil {
		t.Error("Expected user to be returned")
	}
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test GetUser con ID inexistente
	retrievedUser, err := userService.GetUser(999)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrievedUser != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Crear usuario
	user := &domain.User{
		Name:      "Juan Pérez",
		Email:     "juan.perez@email.com",
		Password:  "password123",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Test GetUserByEmail
	retrievedUser, err := userService.GetUserByEmail("juan.perez@email.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrievedUser == nil {
		t.Error("Expected user to be returned")
	}
}

func TestUserService_GetUserByEmail_NotFound(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test GetUserByEmail con email inexistente
	retrievedUser, err := userService.GetUserByEmail("nonexistent@email.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrievedUser != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Crear usuario
	user := &domain.User{
		ID:        1,
		Name:      "Juan Pérez",
		Email:     "juan.perez@email.com",
		Password:  "password123",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Test UpdateUser
	user.Name = "Juan Carlos Pérez"
	err := userService.UpdateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test UpdateUser con usuario inexistente
	user := &domain.User{
		ID:        999,
		Name:      "Usuario Inexistente",
		Email:     "inexistente@email.com",
		Password:  "password123",
		IDNumber:  "99999999",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := userService.UpdateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Crear usuario
	user := &domain.User{
		Name:      "Juan Pérez",
		Email:     "juan.perez@email.com",
		Password:  "password123",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Test DeleteUser
	err := userService.DeleteUser(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test DeleteUser con ID inexistente
	err := userService.DeleteUser(999)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUserService_CreateUser_WithExistingUser(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Crear usuario existente
	existingUser := &domain.User{
		Name:      "Usuario Existente",
		Email:     "existente@email.com",
		Password:  "password123",
		IDNumber:  "11111111",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(existingUser)

	// Crear nuevo usuario
	newUser := &domain.User{
		Name:     "Nuevo Usuario",
		Email:    "nuevo@email.com",
		Password: "password456",
		IDNumber: "22222222",
	}

	err := userService.CreateUser(newUser)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if newUser.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

func TestUserService_GetUser_WithMultipleUsers(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Crear múltiples usuarios
	users := []*domain.User{
		{
			Name:      "Usuario 1",
			Email:     "usuario1@email.com",
			Password:  "password1",
			IDNumber:  "11111111",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Usuario 2",
			Email:     "usuario2@email.com",
			Password:  "password2",
			IDNumber:  "22222222",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Usuario 3",
			Email:     "usuario3@email.com",
			Password:  "password3",
			IDNumber:  "33333333",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Crear usuarios
	for _, user := range users {
		userRepo.Create(user)
	}

	// Test GetUser para cada usuario
	for i, user := range users {
		retrievedUser, err := userService.GetUser(uint(i + 1))
		if err != nil {
			t.Errorf("Expected no error for user %d, got %v", i+1, err)
		}
		if retrievedUser == nil {
			t.Errorf("Expected user %d to be returned", i+1)
		}
		if retrievedUser.Name != user.Name {
			t.Errorf("Expected name %s, got %s", user.Name, retrievedUser.Name)
		}
	}
}

func TestUserService_GetUserByEmail_WithMultipleUsers(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Crear múltiples usuarios
	users := []*domain.User{
		{
			Name:      "Usuario 1",
			Email:     "usuario1@email.com",
			Password:  "password1",
			IDNumber:  "11111111",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Usuario 2",
			Email:     "usuario2@email.com",
			Password:  "password2",
			IDNumber:  "22222222",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Crear usuarios
	for _, user := range users {
		userRepo.Create(user)
	}

	// Test GetUserByEmail para cada usuario
	for _, user := range users {
		retrievedUser, err := userService.GetUserByEmail(user.Email)
		if err != nil {
			t.Errorf("Expected no error for email %s, got %v", user.Email, err)
		}
		if retrievedUser == nil {
			t.Errorf("Expected user with email %s to be returned", user.Email)
		}
		if retrievedUser.Name != user.Name {
			t.Errorf("Expected name %s, got %s", user.Name, retrievedUser.Name)
		}
	}
}

func TestUserService_UpdateUser_WithMultipleUpdates(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Crear usuario
	user := &domain.User{
		ID:        1,
		Name:      "Usuario Original",
		Email:     "original@email.com",
		Password:  "password123",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Múltiples actualizaciones
	updates := []struct {
		name     string
		email    string
		password string
	}{
		{"Usuario Actualizado 1", "actualizado1@email.com", "newpassword1"},
		{"Usuario Actualizado 2", "actualizado2@email.com", "newpassword2"},
		{"Usuario Final", "final@email.com", "finalpassword"},
	}

	for i, update := range updates {
		user.Name = update.name
		user.Email = update.email
		user.Password = update.password

		err := userService.UpdateUser(user)
		if err != nil {
			t.Errorf("Expected no error for update %d, got %v", i+1, err)
		}

		// Verificar actualización
		retrievedUser, err := userService.GetUser(user.ID)
		if err != nil {
			t.Errorf("Expected no error getting user after update %d, got %v", i+1, err)
		}
		if retrievedUser.Name != update.name {
			t.Errorf("Expected name %s after update %d, got %s", update.name, i+1, retrievedUser.Name)
		}
	}
}

func TestUserService_DeleteUser_WithMultipleUsers(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Crear múltiples usuarios
	users := []*domain.User{
		{
			Name:      "Usuario 1",
			Email:     "usuario1@email.com",
			Password:  "password1",
			IDNumber:  "11111111",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Usuario 2",
			Email:     "usuario2@email.com",
			Password:  "password2",
			IDNumber:  "22222222",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "Usuario 3",
			Email:     "usuario3@email.com",
			Password:  "password3",
			IDNumber:  "33333333",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Crear usuarios
	for _, user := range users {
		userRepo.Create(user)
	}

	// Eliminar usuarios en orden inverso
	for i := len(users) - 1; i >= 0; i-- {
		err := userService.DeleteUser(uint(i + 1))
		if err != nil {
			t.Errorf("Expected no error deleting user %d, got %v", i+1, err)
		}

		// Verificar que el usuario fue eliminado
		retrievedUser, err := userService.GetUser(uint(i + 1))
		if err != nil {
			t.Errorf("Expected no error checking deleted user %d, got %v", i+1, err)
		}
		if retrievedUser != nil {
			t.Errorf("Expected user %d to be deleted", i+1)
		}
	}
}

func TestUserService_CreateUser_WithSpecialCharacters(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test con caracteres especiales en el nombre
	user := &domain.User{
		Name:     "José María O'Connor-Smith",
		Email:    "jose.maria@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

func TestUserService_CreateUser_WithUnicodeCharacters(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test con caracteres Unicode
	user := &domain.User{
		Name:     "José María Ñoño",
		Email:    "jose.maria@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

func TestUserService_CreateUser_WithVeryLongName(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test con nombre muy largo
	longName := "Juan Carlos María José Francisco de Paula Juan Nepomuceno María de los Remedios Cipriano de la Santísima Trinidad Ruiz y Picasso"
	user := &domain.User{
		Name:     longName,
		Email:    "picasso@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

func TestUserService_CreateUser_WithPLDServiceError(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := &ErrorMockPLDService{}
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err == nil {
		t.Error("Expected error from PLD service")
	}
}

// ErrorMockPLDService para testing de errores del PLD
type ErrorMockPLDService struct{}

func (m *ErrorMockPLDService) ValidateUser(idNumber, name, email string) (*domain.PLDResponse, error) {
	return nil, fmt.Errorf("PLD service error")
}

func TestUserService_CreateUser_WithRepositoryError(t *testing.T) {
	userRepo := &ErrorMockUserRepository{}
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err == nil {
		t.Error("Expected error from repository")
	}
}

// ErrorMockUserRepository para testing de errores del repositorio
type ErrorMockUserRepository struct{}

func (m *ErrorMockUserRepository) Create(user *domain.User) error {
	return fmt.Errorf("database error")
}

func (m *ErrorMockUserRepository) GetByID(id uint) (*domain.User, error) {
	return nil, fmt.Errorf("database error")
}

func (m *ErrorMockUserRepository) GetByEmail(email string) (*domain.User, error) {
	return nil, fmt.Errorf("database error")
}

func (m *ErrorMockUserRepository) Update(user *domain.User) error {
	return fmt.Errorf("database error")
}

func (m *ErrorMockUserRepository) Delete(id uint) error {
	return fmt.Errorf("database error")
}

func TestUserService_GetUser_WithRepositoryError(t *testing.T) {
	userRepo := &ErrorMockUserRepository{}
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test GetUser con error del repositorio
	user, err := userService.GetUser(1)
	if err == nil {
		t.Error("Expected error from repository")
	}
	if user != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestUserService_GetUserByEmail_WithRepositoryError(t *testing.T) {
	userRepo := &ErrorMockUserRepository{}
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test GetUserByEmail con error del repositorio
	user, err := userService.GetUserByEmail("test@email.com")
	if err == nil {
		t.Error("Expected error from repository")
	}
	if user != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestUserService_UpdateUser_WithRepositoryError(t *testing.T) {
	userRepo := &ErrorMockUserRepository{}
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@email.com",
		Password:  "password123",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test UpdateUser con error del repositorio
	err := userService.UpdateUser(user)
	if err == nil {
		t.Error("Expected error from repository")
	}
}

func TestUserService_DeleteUser_WithRepositoryError(t *testing.T) {
	userRepo := &ErrorMockUserRepository{}
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test DeleteUser con error del repositorio
	err := userService.DeleteUser(1)
	if err == nil {
		t.Error("Expected error from repository")
	}
}

func TestUserService_CreateUser_WithPLDBlacklistedResponse(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := &BlacklistedMockPLDService{}
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err == nil {
		t.Error("Expected error for blacklisted user")
	}

	if err.Error() != "usuario en lista negra: Usuario en lista negra por actividades sospechosas" {
		t.Errorf("Expected blacklist error, got %v", err)
	}
}

// BlacklistedMockPLDService para testing de usuarios en lista negra
type BlacklistedMockPLDService struct{}

func (m *BlacklistedMockPLDService) ValidateUser(idNumber, name, email string) (*domain.PLDResponse, error) {
	return &domain.PLDResponse{
		IsBlacklisted: true,
		Status:        "blacklisted",
		Reason:        "Usuario en lista negra por actividades sospechosas",
	}, nil
}

func TestUserService_CreateUser_WithPLDServiceTimeout(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := &TimeoutMockPLDService{}
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err == nil {
		t.Error("Expected error from PLD service timeout")
	}
}

// TimeoutMockPLDService para testing de timeouts
type TimeoutMockPLDService struct{}

func (m *TimeoutMockPLDService) ValidateUser(idNumber, name, email string) (*domain.PLDResponse, error) {
	return nil, fmt.Errorf("PLD service timeout")
}

func TestUserService_CreateUser_WithDuplicateEmailInRepository(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Crear usuario existente
	existingUser := &domain.User{
		Name:      "Usuario Existente",
		Email:     "duplicate@email.com",
		Password:  "password123",
		IDNumber:  "11111111",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(existingUser)

	// Intentar crear usuario con mismo email
	newUser := &domain.User{
		Name:     "Nuevo Usuario",
		Email:    "duplicate@email.com", // Mismo email
		Password: "password456",
		IDNumber: "22222222",
	}

	err := userService.CreateUser(newUser)
	if err == nil {
		t.Error("Expected error for duplicate email")
	}

	if err.Error() != "el email ya está registrado" {
		t.Errorf("Expected duplicate email error, got %v", err)
	}
}

func TestUserService_CreateUser_WithSpecialCharactersInEmail(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test con email que contiene caracteres especiales
	user := &domain.User{
		Name:     "Test User",
		Email:    "test+special@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

func TestUserService_CreateUser_WithVeryLongEmail(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test con email muy largo
	longEmail := "very.long.email.address.that.exceeds.normal.length.but.should.still.be.valid@very.long.domain.name.com"
	user := &domain.User{
		Name:     "Test User",
		Email:    longEmail,
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

func TestUserService_CreateUser_WithNumericName(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test con nombre que contiene números
	user := &domain.User{
		Name:     "Juan123 Pérez456",
		Email:    "juan123@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

func TestUserService_CreateUser_WithSpecialCharactersInPassword(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test con contraseña que contiene caracteres especiales
	user := &domain.User{
		Name:     "Test User",
		Email:    "test@email.com",
		Password: "p@ssw0rd!@#$%^&*()",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

func TestUserService_CreateUser_WithPLDServiceNetworkError(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := &NetworkErrorMockPLDService{}
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err == nil {
		t.Error("Expected error from PLD service network error")
	}
}

// NetworkErrorMockPLDService para testing de errores de red
type NetworkErrorMockPLDService struct{}

func (m *NetworkErrorMockPLDService) ValidateUser(idNumber, name, email string) (*domain.PLDResponse, error) {
	return nil, fmt.Errorf("network timeout")
}

func TestUserService_GetUser_WithDatabaseConnectionError(t *testing.T) {
	userRepo := &ConnectionErrorMockUserRepository{}
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test GetUser con error de conexión a base de datos
	user, err := userService.GetUser(1)
	if err == nil {
		t.Error("Expected error from database connection")
	}
	if user != nil {
		t.Error("Expected no user to be returned")
	}
}

// ConnectionErrorMockUserRepository para testing de errores de conexión
type ConnectionErrorMockUserRepository struct{}

func (m *ConnectionErrorMockUserRepository) Create(user *domain.User) error {
	return fmt.Errorf("database connection failed")
}

func (m *ConnectionErrorMockUserRepository) GetByID(id uint) (*domain.User, error) {
	return nil, fmt.Errorf("database connection failed")
}

func (m *ConnectionErrorMockUserRepository) GetByEmail(email string) (*domain.User, error) {
	return nil, fmt.Errorf("database connection failed")
}

func (m *ConnectionErrorMockUserRepository) Update(user *domain.User) error {
	return fmt.Errorf("database connection failed")
}

func (m *ConnectionErrorMockUserRepository) Delete(id uint) error {
	return fmt.Errorf("database connection failed")
}

func TestUserService_GetUserByEmail_WithDatabaseConnectionError(t *testing.T) {
	userRepo := &ConnectionErrorMockUserRepository{}
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test GetUserByEmail con error de conexión a base de datos
	user, err := userService.GetUserByEmail("test@email.com")
	if err == nil {
		t.Error("Expected error from database connection")
	}
	if user != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestUserService_UpdateUser_WithDatabaseConnectionError(t *testing.T) {
	userRepo := &ConnectionErrorMockUserRepository{}
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@email.com",
		Password:  "password123",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test UpdateUser con error de conexión a base de datos
	err := userService.UpdateUser(user)
	if err == nil {
		t.Error("Expected error from database connection")
	}
}

func TestUserService_DeleteUser_WithDatabaseConnectionError(t *testing.T) {
	userRepo := &ConnectionErrorMockUserRepository{}
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	// Test DeleteUser con error de conexión a base de datos
	err := userService.DeleteUser(1)
	if err == nil {
		t.Error("Expected error from database connection")
	}
}

func TestUserService_CreateUser_WithPLDServiceUnavailable(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := &UnavailableMockPLDService{}
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err == nil {
		t.Error("Expected error from PLD service unavailable")
	}
}

// UnavailableMockPLDService para testing de servicio no disponible
type UnavailableMockPLDService struct{}

func (m *UnavailableMockPLDService) ValidateUser(idNumber, name, email string) (*domain.PLDResponse, error) {
	return nil, fmt.Errorf("PLD service unavailable")
}

func TestUserService_CreateUser_WithPLDServiceRateLimit(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := &RateLimitMockPLDService{}
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Juan Pérez",
		Email:    "juan.perez@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err == nil {
		t.Error("Expected error from PLD service rate limit")
	}
}

// RateLimitMockPLDService para testing de límite de tasa
type RateLimitMockPLDService struct{}

func (m *RateLimitMockPLDService) ValidateUser(idNumber, name, email string) (*domain.PLDResponse, error) {
	return nil, fmt.Errorf("rate limit exceeded")
}

// Tests adicionales para aumentar cobertura
func TestUserService_CreateUser_WithEmptyName(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "",
		Email:    "test@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error for empty name, got %v", err)
	}
}

func TestUserService_CreateUser_WithEmptyEmail(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Test User",
		Email:    "",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error for empty email, got %v", err)
	}
}

func TestUserService_CreateUser_WithEmptyPassword(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Test User",
		Email:    "test@email.com",
		Password: "",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error for empty password, got %v", err)
	}
}

func TestUserService_CreateUser_WithEmptyIDNumber(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Test User",
		Email:    "test@email.com",
		Password: "password123",
		IDNumber: "",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error for empty ID number, got %v", err)
	}
}

func TestUserService_CreateUser_WithAllEmptyFields(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "",
		Email:    "",
		Password: "",
		IDNumber: "",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error for all empty fields, got %v", err)
	}
}

func TestUserService_GetUser_WithZeroID(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	user, err := userService.GetUser(0)
	if err != nil {
		t.Errorf("Expected no error for zero ID, got %v", err)
	}

	if user != nil {
		t.Error("Expected nil user for zero ID")
	}
}

func TestUserService_GetUserByEmail_WithEmptyEmail(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	user, err := userService.GetUserByEmail("")
	if err != nil {
		t.Errorf("Expected no error for empty email, got %v", err)
	}

	if user != nil {
		t.Error("Expected nil user for empty email")
	}
}

func TestUserService_DeleteUser_WithZeroID(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := NewUserService(userRepo, pldService)

	err := userService.DeleteUser(0)
	if err != nil {
		t.Errorf("Expected no error for zero ID, got %v", err)
	}
}

func TestUserService_CreateUser_WithPLDServiceReturningEmptyResponse(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := &EmptyResponseMockPLDService{}
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Test User",
		Email:    "test@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error for empty PLD response, got %v", err)
	}
}

type EmptyResponseMockPLDService struct{}

func (m *EmptyResponseMockPLDService) ValidateUser(idNumber, name, email string) (*domain.PLDResponse, error) {
	return &domain.PLDResponse{}, nil
}

func TestUserService_CreateUser_WithPLDServiceReturningPartialResponse(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := &PartialResponseMockPLDService{}
	userService := NewUserService(userRepo, pldService)

	user := &domain.User{
		Name:     "Test User",
		Email:    "test@email.com",
		Password: "password123",
		IDNumber: "12345678",
	}

	err := userService.CreateUser(user)
	if err != nil {
		t.Errorf("Expected no error for partial PLD response, got %v", err)
	}
}

type PartialResponseMockPLDService struct{}

func (m *PartialResponseMockPLDService) ValidateUser(idNumber, name, email string) (*domain.PLDResponse, error) {
	return &domain.PLDResponse{
		IsBlacklisted: false,
		Status:        "clean",
	}, nil
}
