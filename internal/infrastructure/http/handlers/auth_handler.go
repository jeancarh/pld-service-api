package handlers

import (
	"crabi-test/internal/application/services"
	"crabi-test/internal/infrastructure/http/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// AuthHandler maneja las solicitudes HTTP relacionadas con autenticación
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler crea una nueva instancia del handler de autenticación
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login godoc
// @Summary Autenticar usuario
// @Description Autentica un usuario con email y contraseña
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Credenciales de login"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Datos de entrada inválidos",
			Details: err.Error(),
		})
		return
	}

	// Validación adicional
	if validate := c.MustGet("validator").(*validator.Validate); validate != nil {
		if err := validate.Struct(req); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Validación fallida",
				Details: err.Error(),
			})
			return
		}
	}

	// Autenticar usuario
	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Error de autenticación",
			Details: err.Error(),
		})
		return
	}

	// Convertir a DTO de respuesta
	userResponse := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IDNumber:  user.IDNumber,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response := dto.LoginResponse{
		Token: token,
		User:  userResponse,
	}

	c.JSON(http.StatusOK, response)
}
