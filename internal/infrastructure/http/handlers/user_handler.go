package handlers

import (
	"crabi-test/internal/application/services"
	"crabi-test/internal/domain"
	"crabi-test/internal/infrastructure/http/dto"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// UserHandler maneja las solicitudes HTTP relacionadas con usuarios
type UserHandler struct {
	userService *services.UserService
	authService *services.AuthService
}

// NewUserHandler crea una nueva instancia del handler de usuarios
func NewUserHandler(userService *services.UserService, authService *services.AuthService) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
	}
}

// CreateUser godoc
// @Summary Crear un nuevo usuario
// @Description Crea un nuevo usuario validando contra el servicio PLD
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "Datos del usuario"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Usuario en lista negra"
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

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

	// Crear usuario en el dominio
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		IDNumber: req.IDNumber,
	}

	// Crear usuario usando el servicio
	if err := h.userService.CreateUser(user); err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "el email ya está registrado" {
			statusCode = http.StatusConflict
		} else if err.Error() == "usuario en lista negra" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Error creando usuario",
			Details: err.Error(),
		})
		return
	}

	// Convertir a DTO de respuesta
	response := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IDNumber:  user.IDNumber,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser godoc
// @Summary Obtener información del usuario autenticado
// @Description Obtiene la información del usuario autenticado
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	// Obtener usuario del contexto (establecido por el middleware de autenticación)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Usuario no autenticado",
		})
		return
	}

	userDomain := user.(*domain.User)

	// Convertir a DTO de respuesta
	response := dto.UserResponse{
		ID:        userDomain.ID,
		Name:      userDomain.Name,
		Email:     userDomain.Email,
		IDNumber:  userDomain.IDNumber,
		CreatedAt: userDomain.CreatedAt,
		UpdatedAt: userDomain.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetUserByID godoc
// @Summary Obtener usuario por ID
// @Description Obtiene la información de un usuario por su ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "ID inválido",
			Details: err.Error(),
		})
		return
	}

	user, err := h.userService.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Error obteniendo usuario",
			Details: err.Error(),
		})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Usuario no encontrado",
		})
		return
	}

	// Convertir a DTO de respuesta
	response := dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IDNumber:  user.IDNumber,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser godoc
// @Summary Eliminar usuario
// @Description Elimina un usuario por su ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID del usuario"
// @Security BearerAuth
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "ID inválido",
			Details: err.Error(),
		})
		return
	}

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "usuario no encontrado" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Error eliminando usuario",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Usuario eliminado correctamente",
	})
}
