package tests

import (
	"crabi-test/internal/application/services"
	"crabi-test/internal/domain"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := services.NewAuthService(userRepo)

	// Crear usuario con contraseña encriptada
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        1,
		Name:      "Juan Pérez",
		Email:     "juan.perez@email.com",
		Password:  string(hashedPassword),
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Act
	resultUser, token, err := authService.Login("juan.perez@email.com", "password123")

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resultUser == nil {
		t.Error("Expected user to be returned")
	}

	if token == "" {
		t.Error("Expected token to be generated")
	}

	if resultUser.Email != "juan.perez@email.com" {
		t.Errorf("Expected email %s, got %s", "juan.perez@email.com", resultUser.Email)
	}
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := services.NewAuthService(userRepo)

	// Act
	user, token, err := authService.Login("nonexistent@email.com", "wrongpassword")

	// Assert
	if err == nil {
		t.Error("Expected error for invalid credentials")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}

	if token != "" {
		t.Error("Expected no token to be generated")
	}

	if err.Error() != "credenciales inválidas" {
		t.Errorf("Expected 'credenciales inválidas' error, got %v", err)
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := services.NewAuthService(userRepo)

	// Crear usuario con contraseña encriptada
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        1,
		Name:      "Juan Pérez",
		Email:     "juan.perez@email.com",
		Password:  string(hashedPassword),
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Act
	resultUser, token, err := authService.Login("juan.perez@email.com", "wrongpassword")

	// Assert
	if err == nil {
		t.Error("Expected error for wrong password")
	}

	if resultUser != nil {
		t.Error("Expected no user to be returned")
	}

	if token != "" {
		t.Error("Expected no token to be generated")
	}

	if err.Error() != "credenciales inválidas" {
		t.Errorf("Expected 'credenciales inválidas' error, got %v", err)
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := services.NewAuthService(userRepo)

	user := &domain.User{
		ID:        1,
		Name:      "Juan Pérez",
		Email:     "juan.perez@email.com",
		Password:  "hashedpassword",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	token, err := authService.GenerateToken(user)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected token to be generated")
	}
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := services.NewAuthService(userRepo)

	user := &domain.User{
		ID:        1,
		Name:      "Juan Pérez",
		Email:     "juan.perez@email.com",
		Password:  "hashedpassword",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Generar token
	token, _ := authService.GenerateToken(user)

	// Act
	validatedUser, err := authService.ValidateToken(token)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if validatedUser == nil {
		t.Error("Expected user to be returned")
	}

	if validatedUser.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, validatedUser.Email)
	}
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	// Arrange
	userRepo := NewMockUserRepository()
	authService := services.NewAuthService(userRepo)

	// Act
	user, err := authService.ValidateToken("invalid.token.here")

	// Assert
	if err == nil {
		t.Error("Expected error for invalid token")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}

	if err.Error() != "token inválido" {
		t.Errorf("Expected 'token inválido' error, got %v", err)
	}
}
