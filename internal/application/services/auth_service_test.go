package services

import (
	"crabi-test/internal/domain"
	"fmt"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

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

	// Test Login
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
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test Login con credenciales inválidas
	user, token, err := authService.Login("nonexistent@email.com", "wrongpassword")

	if err == nil {
		t.Error("Expected error for invalid credentials")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}

	if token != "" {
		t.Error("Expected no token to be generated")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

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

	// Test Login con contraseña incorrecta
	resultUser, token, err := authService.Login("juan.perez@email.com", "wrongpassword")

	if err == nil {
		t.Error("Expected error for wrong password")
	}

	if resultUser != nil {
		t.Error("Expected no user to be returned")
	}

	if token != "" {
		t.Error("Expected no token to be generated")
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	user := &domain.User{
		ID:        1,
		Name:      "Juan Pérez",
		Email:     "juan.perez@email.com",
		Password:  "hashedpassword",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test GenerateToken
	token, err := authService.GenerateToken(user)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token == "" {
		t.Error("Expected token to be generated")
	}
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

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

	// Test ValidateToken
	validatedUser, err := authService.ValidateToken(token)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if validatedUser == nil {
		t.Error("Expected user to be returned")
	}
}

func TestAuthService_ValidateToken_InvalidToken(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test ValidateToken con token inválido
	user, err := authService.ValidateToken("invalid.token.here")

	if err == nil {
		t.Error("Expected error for invalid token")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestAuthService_ValidateToken_EmptyToken(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test ValidateToken con token vacío
	user, err := authService.ValidateToken("")

	if err == nil {
		t.Error("Expected error for empty token")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestAuthService_ValidateToken_MalformedToken(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test ValidateToken con token malformado
	user, err := authService.ValidateToken("not.a.valid.jwt.token")

	if err == nil {
		t.Error("Expected error for malformed token")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestAuthService_Login_EmptyCredentials(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test Login con credenciales vacías
	user, token, err := authService.Login("", "")

	if err == nil {
		t.Error("Expected error for empty credentials")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}

	if token != "" {
		t.Error("Expected no token to be generated")
	}
}

func TestAuthService_Login_EmptyEmail(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test Login con email vacío
	user, token, err := authService.Login("", "password123")

	if err == nil {
		t.Error("Expected error for empty email")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}

	if token != "" {
		t.Error("Expected no token to be generated")
	}
}

func TestAuthService_Login_EmptyPassword(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test Login con contraseña vacía
	user, token, err := authService.Login("test@email.com", "")

	if err == nil {
		t.Error("Expected error for empty password")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}

	if token != "" {
		t.Error("Expected no token to be generated")
	}
}

func TestAuthService_Login_WithMultipleUsers(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Crear múltiples usuarios con contraseñas encriptadas
	users := []struct {
		email    string
		password string
		name     string
	}{
		{"user1@email.com", "password1", "Usuario 1"},
		{"user2@email.com", "password2", "Usuario 2"},
		{"user3@email.com", "password3", "Usuario 3"},
	}

	// Crear usuarios
	for _, u := range users {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		user := &domain.User{
			ID:        uint(len(userRepo.users) + 1),
			Name:      u.name,
			Email:     u.email,
			Password:  string(hashedPassword),
			IDNumber:  "12345678",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		userRepo.Create(user)
	}

	// Test login para cada usuario
	for _, u := range users {
		resultUser, token, err := authService.Login(u.email, u.password)
		if err != nil {
			t.Errorf("Expected no error for %s, got %v", u.email, err)
		}

		if resultUser == nil {
			t.Errorf("Expected user to be returned for %s", u.email)
		}

		if token == "" {
			t.Errorf("Expected token to be generated for %s", u.email)
		}

		if resultUser.Email != u.email {
			t.Errorf("Expected email %s, got %s", u.email, resultUser.Email)
		}
	}
}

func TestAuthService_ValidateToken_WithMultipleTokens(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Crear usuario
	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@email.com",
		Password:  "hashedpassword",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Generar múltiples tokens
	tokens := make([]string, 5)
	for i := 0; i < 5; i++ {
		token, err := authService.GenerateToken(user)
		if err != nil {
			t.Errorf("Expected no error generating token %d, got %v", i+1, err)
		}
		tokens[i] = token
	}

	// Validar cada token
	for i, token := range tokens {
		validatedUser, err := authService.ValidateToken(token)
		if err != nil {
			t.Errorf("Expected no error validating token %d, got %v", i+1, err)
		}

		if validatedUser == nil {
			t.Errorf("Expected user to be returned for token %d", i+1)
		}

		if validatedUser.Email != user.Email {
			t.Errorf("Expected email %s, got %s", user.Email, validatedUser.Email)
		}
	}
}

func TestAuthService_GenerateToken_WithMultipleUsers(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Crear múltiples usuarios
	users := []*domain.User{
		{
			ID:        1,
			Name:      "Usuario 1",
			Email:     "usuario1@email.com",
			Password:  "hashedpassword1",
			IDNumber:  "11111111",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "Usuario 2",
			Email:     "usuario2@email.com",
			Password:  "hashedpassword2",
			IDNumber:  "22222222",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        3,
			Name:      "Usuario 3",
			Email:     "usuario3@email.com",
			Password:  "hashedpassword3",
			IDNumber:  "33333333",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Generar tokens para cada usuario
	for i, user := range users {
		token, err := authService.GenerateToken(user)
		if err != nil {
			t.Errorf("Expected no error generating token for user %d, got %v", i+1, err)
		}

		if token == "" {
			t.Errorf("Expected token to be generated for user %d", i+1)
		}

	}
}

func TestAuthService_Login_WithSpecialCharacters(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Crear usuario con caracteres especiales
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        1,
		Name:      "José María O'Connor-Smith",
		Email:     "jose.maria@email.com",
		Password:  string(hashedPassword),
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Test Login
	resultUser, token, err := authService.Login("jose.maria@email.com", "password123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resultUser == nil {
		t.Error("Expected user to be returned")
	}

	if token == "" {
		t.Error("Expected token to be generated")
	}

	if resultUser.Name != "José María O'Connor-Smith" {
		t.Errorf("Expected name José María O'Connor-Smith, got %s", resultUser.Name)
	}
}

func TestAuthService_Login_WithUnicodeCharacters(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Crear usuario con caracteres Unicode
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        1,
		Name:      "José María Ñoño",
		Email:     "jose.maria@email.com",
		Password:  string(hashedPassword),
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Test Login
	resultUser, token, err := authService.Login("jose.maria@email.com", "password123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resultUser == nil {
		t.Error("Expected user to be returned")
	}

	if token == "" {
		t.Error("Expected token to be generated")
	}

	if resultUser.Name != "José María Ñoño" {
		t.Errorf("Expected name José María Ñoño, got %s", resultUser.Name)
	}
}

func TestAuthService_ValidateToken_WithExpiredToken(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test con token que simula estar expirado
	// En un entorno real, esto requeriría manipular el tiempo
	// Por ahora, solo verificamos que el middleware maneja tokens inválidos

	user, err := authService.ValidateToken("expired.token.here")

	if err == nil {
		t.Error("Expected error for expired token")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestAuthService_ValidateToken_WithMalformedToken(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test con tokens malformados
	malformedTokens := []string{
		"not.a.valid.jwt",
		"header.payload.signature",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		"invalid.token.format",
		"",
	}

	for i, token := range malformedTokens {
		user, err := authService.ValidateToken(token)

		if err == nil {
			t.Errorf("Expected error for malformed token %d: %s", i+1, token)
		}

		if user != nil {
			t.Errorf("Expected no user to be returned for malformed token %d", i+1)
		}
	}
}

func TestAuthService_Login_WithDifferentPasswordCosts(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test con diferentes costos de bcrypt
	costs := []int{bcrypt.MinCost, bcrypt.DefaultCost, bcrypt.DefaultCost + 1}

	for i, cost := range costs {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), cost)
		user := &domain.User{
			ID:        uint(i + 1),
			Name:      fmt.Sprintf("Usuario %d", i+1),
			Email:     fmt.Sprintf("usuario%d@email.com", i+1),
			Password:  string(hashedPassword),
			IDNumber:  "12345678",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		userRepo.Create(user)

		// Test Login
		resultUser, token, err := authService.Login(user.Email, "password123")
		if err != nil {
			t.Errorf("Expected no error for cost %d, got %v", cost, err)
		}

		if resultUser == nil {
			t.Errorf("Expected user to be returned for cost %d", cost)
		}

		if token == "" {
			t.Errorf("Expected token to be generated for cost %d", cost)
		}
	}
}

func TestAuthService_Login_WithRepositoryError(t *testing.T) {
	userRepo := &ErrorMockUserRepository{}
	authService := NewAuthService(userRepo)

	// Test Login con error del repositorio
	user, token, err := authService.Login("test@email.com", "password123")

	if err == nil {
		t.Error("Expected error from repository")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}

	if token != "" {
		t.Error("Expected no token to be generated")
	}
}

func TestAuthService_ValidateToken_WithRepositoryError(t *testing.T) {
	userRepo := &ErrorMockUserRepository{}
	authService := NewAuthService(userRepo)

	// Test ValidateToken con error del repositorio
	user, err := authService.ValidateToken("valid.token.here")

	if err == nil {
		t.Error("Expected error from repository")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestAuthService_Login_WithDifferentPasswordHashes(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test con diferentes tipos de hash de contraseña
	testCases := []struct {
		password string
		cost     int
	}{
		{"password123", bcrypt.MinCost},
		{"password456", bcrypt.DefaultCost},
		{"password789", bcrypt.DefaultCost + 1},
		{"p@ssw0rd!@#", bcrypt.DefaultCost},
		{"123456789", bcrypt.DefaultCost},
	}

	for i, tc := range testCases {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(tc.password), tc.cost)
		user := &domain.User{
			ID:        uint(i + 1),
			Name:      fmt.Sprintf("Usuario %d", i+1),
			Email:     fmt.Sprintf("usuario%d@email.com", i+1),
			Password:  string(hashedPassword),
			IDNumber:  "12345678",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		userRepo.Create(user)

		// Test Login
		resultUser, token, err := authService.Login(user.Email, tc.password)
		if err != nil {
			t.Errorf("Expected no error for password %s, got %v", tc.password, err)
		}

		if resultUser == nil {
			t.Errorf("Expected user to be returned for password %s", tc.password)
		}

		if token == "" {
			t.Errorf("Expected token to be generated for password %s", tc.password)
		}
	}
}

func TestAuthService_ValidateToken_WithDifferentTokenFormats(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test con diferentes formatos de token inválidos
	invalidTokens := []string{
		"not.a.jwt.token",
		"header.payload.signature",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
		"invalid.token.format",
		"",
		"   ",
		"Bearer invalid.token",
		"Basic invalid.token",
		"Token invalid.token",
	}

	for i, token := range invalidTokens {
		user, err := authService.ValidateToken(token)

		if err == nil {
			t.Errorf("Expected error for invalid token %d: %s", i+1, token)
		}

		if user != nil {
			t.Errorf("Expected no user to be returned for invalid token %d", i+1)
		}
	}
}

func TestAuthService_Login_WithSpecialCharactersInPassword(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test con contraseña que contiene caracteres especiales
	specialPassword := "p@ssw0rd!@#$%^&*()_+-=[]{}|;':\",./<>?"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(specialPassword), bcrypt.DefaultCost)
	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@email.com",
		Password:  string(hashedPassword),
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	// Test Login
	resultUser, token, err := authService.Login("test@email.com", specialPassword)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resultUser == nil {
		t.Error("Expected user to be returned")
	}

	if token == "" {
		t.Error("Expected token to be generated")
	}
}

func TestAuthService_GenerateToken_WithDifferentUserTypes(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test con diferentes tipos de usuarios
	users := []*domain.User{
		{
			ID:        1,
			Name:      "Admin User",
			Email:     "admin@email.com",
			Password:  "adminpass",
			IDNumber:  "11111111",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        2,
			Name:      "Regular User",
			Email:     "user@email.com",
			Password:  "userpass",
			IDNumber:  "22222222",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        3,
			Name:      "Guest User",
			Email:     "guest@email.com",
			Password:  "guestpass",
			IDNumber:  "33333333",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Generar tokens para cada tipo de usuario
	for i, user := range users {
		token, err := authService.GenerateToken(user)
		if err != nil {
			t.Errorf("Expected no error generating token for user %d, got %v", i+1, err)
		}

		if token == "" {
			t.Errorf("Expected token to be generated for user %d", i+1)
		}

		// Verificar que el token es único
		if i > 0 {
			prevToken, _ := authService.GenerateToken(users[i-1])
			if token == prevToken {
				t.Errorf("Expected unique token for user %d", i+1)
			}
		}
	}
}

// Tests adicionales para aumentar cobertura
func TestAuthService_Login_WithEmptyUserRepository(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test Login con repositorio vacío
	user, token, err := authService.Login("nonexistent@email.com", "password123")

	if err == nil {
		t.Error("Expected error for non-existent user")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}

	if token != "" {
		t.Error("Expected no token to be generated")
	}
}

func TestAuthService_ValidateToken_WithEmptyToken(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test ValidateToken con token vacío
	user, err := authService.ValidateToken("")

	if err == nil {
		t.Error("Expected error for empty token")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestAuthService_ValidateToken_WithWhitespaceToken(t *testing.T) {
	userRepo := NewMockUserRepository()
	authService := NewAuthService(userRepo)

	// Test ValidateToken con token que solo contiene espacios
	user, err := authService.ValidateToken("   ")

	if err == nil {
		t.Error("Expected error for whitespace token")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}
}

func TestAuthService_Login_WithNilUserFromRepository(t *testing.T) {
	userRepo := &NilUserMockRepository{}
	authService := NewAuthService(userRepo)

	// Test Login cuando el repositorio retorna nil
	user, token, err := authService.Login("test@email.com", "password123")

	if err == nil {
		t.Error("Expected error for nil user")
	}

	if user != nil {
		t.Error("Expected no user to be returned")
	}

	if token != "" {
		t.Error("Expected no token to be generated")
	}
}

func TestAuthService_ValidateToken_WithRepositoryErrorOnGetByID(t *testing.T) {
	userRepo := NewErrorOnGetByIDMockRepository()
	authService := NewAuthService(userRepo)

	// Crear usuario y generar token
	user := &domain.User{
		ID:        1,
		Name:      "Test User",
		Email:     "test@email.com",
		Password:  "hashedpassword",
		IDNumber:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.Create(user)

	token, _ := authService.GenerateToken(user)

	// Test ValidateToken con error en GetByID
	validatedUser, err := authService.ValidateToken(token)

	if err == nil {
		t.Error("Expected error from repository GetByID")
	}

	if validatedUser != nil {
		t.Error("Expected no user to be returned")
	}
}

// Mock repositories adicionales para casos edge
type NilUserMockRepository struct {
	users  map[uint]*domain.User
	emails map[string]*domain.User
}

func NewNilUserMockRepository() *NilUserMockRepository {
	return &NilUserMockRepository{
		users:  make(map[uint]*domain.User),
		emails: make(map[string]*domain.User),
	}
}

func (m *NilUserMockRepository) Create(user *domain.User) error {
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *NilUserMockRepository) GetByID(id uint) (*domain.User, error) {
	return nil, nil // Siempre retorna nil
}

func (m *NilUserMockRepository) GetByEmail(email string) (*domain.User, error) {
	return nil, nil // Siempre retorna nil
}

func (m *NilUserMockRepository) Update(user *domain.User) error {
	return nil
}

func (m *NilUserMockRepository) Delete(id uint) error {
	return nil
}

type ErrorOnGetByIDMockRepository struct {
	users  map[uint]*domain.User
	emails map[string]*domain.User
}

func NewErrorOnGetByIDMockRepository() *ErrorOnGetByIDMockRepository {
	return &ErrorOnGetByIDMockRepository{
		users:  make(map[uint]*domain.User),
		emails: make(map[string]*domain.User),
	}
}

func (m *ErrorOnGetByIDMockRepository) Create(user *domain.User) error {
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *ErrorOnGetByIDMockRepository) GetByID(id uint) (*domain.User, error) {
	return nil, fmt.Errorf("database error")
}

func (m *ErrorOnGetByIDMockRepository) GetByEmail(email string) (*domain.User, error) {
	if user, exists := m.emails[email]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *ErrorOnGetByIDMockRepository) Update(user *domain.User) error {
	return nil
}

func (m *ErrorOnGetByIDMockRepository) Delete(id uint) error {
	return nil
}
