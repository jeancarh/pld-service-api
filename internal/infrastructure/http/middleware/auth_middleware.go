package middleware

import (
	"crabi-test/internal/application/services"
	"crabi-test/internal/infrastructure/http/dto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware middleware para autenticación JWT
type AuthMiddleware struct {
	authService *services.AuthService
}

// NewAuthMiddleware crea una nueva instancia del middleware de autenticación
func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Authenticate middleware para validar token JWT
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Token de autorización requerido",
			})
			c.Abort()
			return
		}

		// Verificar formato del token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Formato de token inválido. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validar token
		user, err := m.authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:   "Token inválido",
				Details: err.Error(),
			})
			c.Abort()
			return
		}

		// Establecer usuario en el contexto
		c.Set("user", user)
		c.Next()
	}
}
