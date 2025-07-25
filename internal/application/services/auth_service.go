package services

import (
	"crabi-test/internal/application/ports"
	"crabi-test/internal/domain"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService implementa la lógica de autenticación
type AuthService struct {
	userRepo ports.UserRepository
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(userRepo ports.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Login autentica un usuario y retorna un token JWT
func (s *AuthService) Login(email, password string) (*domain.User, string, error) {
	// Buscar usuario por email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil || user == nil {
		return nil, "", errors.New("credenciales inválidas")
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("credenciales inválidas")
	}

	// Generar token JWT
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", errors.New("error generando token")
	}

	return user, token, nil
}

// GenerateToken genera un token JWT para un usuario
func (s *AuthService) GenerateToken(user *domain.User) (string, error) {
	// Obtener secret key del environment
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "crabi-jwt-secret-key-for-development-only"
	}

	// Crear claims del token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 horas
		"iat":     time.Now().Unix(),
	}

	// Crear token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar token
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken valida un token JWT y retorna el usuario
func (s *AuthService) ValidateToken(tokenString string) (*domain.User, error) {
	// Obtener secret key del environment
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "crabi-jwt-secret-key-for-development-only"
	}

	// Parsear token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, errors.New("token inválido")
	}

	// Verificar que el token sea válido
	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	// Extraer claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("token inválido")
	}

	// Obtener user_id del token
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("token inválido")
	}

	// Buscar usuario en base de datos
	user, err := s.userRepo.GetByID(uint(userID))
	if err != nil {
		return nil, errors.New("usuario no encontrado")
	}

	return user, nil
}
