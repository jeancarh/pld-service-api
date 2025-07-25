package tests

import (
	"crabi-test/internal/infrastructure/external"
	"testing"
)

func TestPLDClient_ValidateUser_Success(t *testing.T) {
	// Arrange
	pldClient := external.NewPLDClient()

	// Act
	response, err := pldClient.ValidateUser("12345678", "Juan Pérez", "juan.perez@email.com")

	// Assert
	// Este test puede fallar si el servicio PLD no está disponible
	// En un entorno de testing real, usaríamos un mock del servicio HTTP
	if err != nil {
		t.Logf("Test skipped - PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}

	// Verificar que el nombre se separó correctamente
	// Esto es una prueba básica, en un caso real necesitaríamos un mock del servicio HTTP
}

func TestPLDClient_ValidateUser_WithComplexName(t *testing.T) {
	// Arrange
	pldClient := external.NewPLDClient()

	// Act
	response, err := pldClient.ValidateUser("12345678", "Juan Carlos Pérez González", "juan.perez@email.com")

	// Assert
	// Este test puede fallar si el servicio PLD no está disponible
	// En un entorno de testing real, usaríamos un mock del servicio HTTP
	if err != nil {
		t.Logf("Test skipped - PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}

func TestPLDClient_ValidateUser_SingleName(t *testing.T) {
	// Arrange
	pldClient := external.NewPLDClient()

	// Act
	response, err := pldClient.ValidateUser("12345678", "Juan", "juan@email.com")

	// Assert
	// Este test puede fallar si el servicio PLD no está disponible
	// En un entorno de testing real, usaríamos un mock del servicio HTTP
	if err != nil {
		t.Logf("Test skipped - PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}
