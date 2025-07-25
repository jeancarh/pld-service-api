package tests

import (
	"crabi-test/internal/application/services"
	"crabi-test/internal/domain"
	"crabi-test/internal/infrastructure/external"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// TestUserServiceCoverage cubre todas las funciones del UserService
func TestUserServiceCoverage(t *testing.T) {
	userRepo := NewMockUserRepository()
	pldService := NewMockPLDService(false)
	userService := services.NewUserService(userRepo, pldService)

	// Test CreateUser - Success
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

	// Test GetUser
	retrievedUser, err := userService.GetUser(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrievedUser == nil {
		t.Error("Expected user to be returned")
	}

	// Test GetUserByEmail
	userByEmail, err := userService.GetUserByEmail("juan.perez@email.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if userByEmail == nil {
		t.Error("Expected user to be returned")
	}

	// Test UpdateUser
	user.Name = "Juan Carlos Pérez"
	err = userService.UpdateUser(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test DeleteUser
	err = userService.DeleteUser(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestAuthServiceCoverage cubre todas las funciones del AuthService
func TestAuthServiceCoverage(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := services.NewAuthService(userRepo)

	// Crear usuario para testing
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

	// Test Login - Success
	resultUser, token, err := authService.Login("juan.perez@email.com", "password123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if resultUser == nil {
		t.Error("Expected user to be returned")
	}
	if token == "" {
		t.Error("Expected token to be generated")
	}

	// Test GenerateToken
	newToken, err := authService.GenerateToken(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if newToken == "" {
		t.Error("Expected token to be generated")
	}

	// Test ValidateToken - Success
	validatedUser, err := authService.ValidateToken(token)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if validatedUser == nil {
		t.Error("Expected user to be returned")
	}

	// Test ValidateToken - Invalid
	_, err = authService.ValidateToken("invalid.token.here")
	if err == nil {
		t.Error("Expected error for invalid token")
	}
}

// TestPLDClientCoverage cubre las funciones del PLDClient
func TestPLDClientCoverage(t *testing.T) {
	pldClient := external.NewPLDClient()

	// Test ValidateUser con diferentes tipos de nombres
	testCases := []struct {
		name     string
		idNumber string
		email    string
	}{
		{"Juan Pérez", "12345678", "juan.perez@email.com"},
		{"Juan Carlos Pérez González", "87654321", "juan.carlos@email.com"},
		{"Ana", "11111111", "ana@email.com"},
	}

	for _, tc := range testCases {
		response, err := pldClient.ValidateUser(tc.idNumber, tc.name, tc.email)
		if err != nil {
			t.Logf("PLD service not available for %s: %v", tc.name, err)
			continue
		}
		if response == nil {
			t.Errorf("Expected response for %s", tc.name)
		}
	}
}

// TestMockRepositoryCoverage cubre todas las funciones del mock repository
func TestMockRepositoryCoverage(t *testing.T) {
	repo := NewMockUserRepository()

	// Test Create
	user := &domain.User{
		Name:      "Test User",
		Email:     "test@email.com",
		Password:  "password123",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test GetByID
	retrievedUser, err := repo.GetByID(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrievedUser == nil {
		t.Error("Expected user to be returned")
	}

	// Test GetByEmail
	userByEmail, err := repo.GetByEmail("test@email.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if userByEmail == nil {
		t.Error("Expected user to be returned")
	}

	// Test Update
	user.Name = "Updated User"
	err = repo.Update(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test Delete
	err = repo.Delete(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test GetByID after delete
	deletedUser, err := repo.GetByID(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if deletedUser != nil {
		t.Error("Expected no user to be returned after delete")
	}
}

// TestMockPLDServiceCoverage cubre las funciones del mock PLD service
func TestMockPLDServiceCoverage(t *testing.T) {
	// Test clean user
	cleanService := NewMockPLDService(false)
	response, err := cleanService.ValidateUser("12345678", "Juan Pérez", "juan@email.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response == nil {
		t.Error("Expected response to be returned")
	}
	if response.IsBlacklisted {
		t.Error("Expected user to be clean")
	}

	// Test blacklisted user
	blacklistedService := NewMockPLDService(true)
	response, err = blacklistedService.ValidateUser("12345678", "Juan Pérez", "juan@email.com")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response == nil {
		t.Error("Expected response to be returned")
	}
	if !response.IsBlacklisted {
		t.Error("Expected user to be blacklisted")
	}
}
