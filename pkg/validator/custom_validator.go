package validator

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// CustomValidator middleware para validación personalizada
func CustomValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		validate := validator.New()

		// Registrar validaciones personalizadas
		validate.RegisterValidation("id_number", validateIDNumber)

		c.Set("validator", validate)
		c.Next()
	}
}

// validateIDNumber valida formato de número de identificación
func validateIDNumber(fl validator.FieldLevel) bool {
	// Implementar lógica de validación específica
	// Por ahora solo verifica que tenga al menos 8 caracteres
	return len(fl.Field().String()) >= 8
}
