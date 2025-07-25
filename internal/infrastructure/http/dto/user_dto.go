package dto

import "time"

// CreateUserRequest representa la solicitud para crear un usuario
// @Description Solicitud para crear un nuevo usuario
type CreateUserRequest struct {
	// @Description Nombre completo del usuario
	// @Example "Juan Pérez"
	// @Required
	Name string `json:"name" binding:"required,min=2,max=100" example:"Juan Pérez"`

	// @Description Email del usuario (debe ser único)
	// @Example "juan.perez@email.com"
	// @Required
	Email string `json:"email" binding:"required,email" example:"juan.perez@email.com"`

	// @Description Contraseña del usuario (mínimo 8 caracteres)
	// @Example "password123"
	// @Required
	Password string `json:"password" binding:"required,min=8" example:"password123"`

	// @Description Número de identificación personal
	// @Example "12345678"
	// @Required
	IDNumber string `json:"id_number" binding:"required,min=8,max=20" example:"12345678"`
}

// LoginRequest representa la solicitud de login
// @Description Solicitud para autenticarse en el sistema
type LoginRequest struct {
	// @Description Email del usuario
	// @Example "juan.perez@email.com"
	// @Required
	Email string `json:"email" binding:"required,email" example:"juan.perez@email.com"`

	// @Description Contraseña del usuario
	// @Example "password123"
	// @Required
	Password string `json:"password" binding:"required" example:"password123"`
}

// UserResponse representa la respuesta de usuario
// @Description Información del usuario
type UserResponse struct {
	// @Description ID único del usuario
	// @Example "1"
	ID uint `json:"id" example:"1"`

	// @Description Nombre completo del usuario
	// @Example "Juan Pérez"
	Name string `json:"name" example:"Juan Pérez"`

	// @Description Email del usuario
	// @Example "juan.perez@email.com"
	Email string `json:"email" example:"juan.perez@email.com"`

	// @Description Número de identificación personal
	// @Example "12345678"
	IDNumber string `json:"id_number" example:"12345678"`

	// @Description Fecha de creación del usuario
	// @Example "2024-01-15T10:30:00Z"
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`

	// @Description Fecha de última actualización del usuario
	// @Example "2024-01-15T10:30:00Z"
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// LoginResponse representa la respuesta de login
// @Description Respuesta de autenticación exitosa
type LoginResponse struct {
	// @Description Token JWT para autenticación
	// @Example "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`

	// @Description Información del usuario autenticado
	User UserResponse `json:"user"`
}

// ErrorResponse representa una respuesta de error
// @Description Respuesta de error
type ErrorResponse struct {
	// @Description Mensaje de error
	// @Example "Error de validación"
	Error string `json:"error" example:"Error de validación"`

	// @Description Detalles adicionales del error
	// @Example "El campo email es requerido"
	Details string `json:"details,omitempty" example:"El campo email es requerido"`
}

// SuccessResponse representa una respuesta exitosa
// @Description Respuesta de operación exitosa
type SuccessResponse struct {
	// @Description Mensaje de éxito
	// @Example "Usuario eliminado correctamente"
	Message string `json:"message" example:"Usuario eliminado correctamente"`
}
