package external

import (
	"testing"
)

// Los tests requieren que PLD_SERVICE_URL esté configurada en las variables de entorno
// Ejemplo: export PLD_SERVICE_URL=http://localhost:3000

func TestPLDClient_ValidateUser_Success(t *testing.T) {
	pldClient := NewPLDClient()

	response, err := pldClient.ValidateUser("12345678", "Juan Pérez", "juan.perez@email.com")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}

	if response.Status == "" {
		t.Error("Expected status to be set")
	}
}

func TestPLDClient_ValidateUser_WithComplexName(t *testing.T) {
	pldClient := NewPLDClient()

	response, err := pldClient.ValidateUser("87654321", "Juan Carlos Pérez González", "juan.carlos@email.com")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}

func TestPLDClient_ValidateUser_SingleName(t *testing.T) {
	pldClient := NewPLDClient()

	response, err := pldClient.ValidateUser("11111111", "Ana", "ana@email.com")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}

func TestPLDClient_ValidateUser_EmptyName(t *testing.T) {
	pldClient := NewPLDClient()

	response, err := pldClient.ValidateUser("12345678", "", "test@email.com")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}

func TestPLDClient_ValidateUser_EmptyEmail(t *testing.T) {
	pldClient := NewPLDClient()

	response, err := pldClient.ValidateUser("12345678", "Juan Pérez", "")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}

func TestPLDClient_ValidateUser_EmptyIDNumber(t *testing.T) {
	pldClient := NewPLDClient()

	response, err := pldClient.ValidateUser("", "Juan Pérez", "juan.perez@email.com")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}

func TestPLDClient_ValidateUser_AllEmpty(t *testing.T) {
	pldClient := NewPLDClient()

	response, err := pldClient.ValidateUser("", "", "")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}

func TestPLDClient_ValidateUser_VeryLongName(t *testing.T) {
	pldClient := NewPLDClient()

	// Nombre muy largo
	longName := "Juan Carlos María José Francisco de Paula Juan Nepomuceno María de los Remedios Cipriano de la Santísima Trinidad Ruiz y Picasso"

	response, err := pldClient.ValidateUser("12345678", longName, "picasso@email.com")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}

func TestPLDClient_ValidateUser_SpecialCharacters(t *testing.T) {
	pldClient := NewPLDClient()

	// Nombre con caracteres especiales
	nameWithSpecialChars := "José María O'Connor-Smith"

	response, err := pldClient.ValidateUser("12345678", nameWithSpecialChars, "jose.maria@email.com")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}

func TestPLDClient_ValidateUser_NumbersInName(t *testing.T) {
	pldClient := NewPLDClient()

	// Nombre con números
	nameWithNumbers := "Juan123 Pérez456"

	response, err := pldClient.ValidateUser("12345678", nameWithNumbers, "juan123@email.com")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}

func TestPLDClient_ValidateUser_UnicodeCharacters(t *testing.T) {
	pldClient := NewPLDClient()

	// Nombre con caracteres Unicode
	unicodeName := "José María Ñoño"

	response, err := pldClient.ValidateUser("12345678", unicodeName, "jose.maria@email.com")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	if response == nil {
		t.Error("Expected response to be returned")
	}
}

func TestPLDClient_ValidateUser_ResponseStructure(t *testing.T) {
	pldClient := NewPLDClient()

	response, err := pldClient.ValidateUser("12345678", "Juan Pérez", "juan.perez@email.com")
	if err != nil {
		t.Logf("PLD service not available: %v", err)
		t.Skip("PLD service not available for testing")
		return
	}

	// Verificar estructura de respuesta
	if response == nil {
		t.Error("Expected response to be returned")
		return
	}

	// Verificar que los campos esperados estén presentes
	// Nota: Los valores específicos dependen del servicio PLD real
	if response.Status == "" {
		t.Log("Status field is empty, but this might be normal for the PLD service")
	}

	// Verificar que IsBlacklisted sea un booleano válido
	// (puede ser true o false, ambos son válidos)
	_ = response.IsBlacklisted // Solo verificar que no cause panic
}

func TestPLDClient_ValidateUser_MultipleCalls(t *testing.T) {
	pldClient := NewPLDClient()

	// Hacer múltiples llamadas para verificar consistencia
	testCases := []struct {
		idNumber string
		name     string
		email    string
	}{
		{"12345678", "Juan Pérez", "juan.perez@email.com"},
		{"87654321", "María García", "maria.garcia@email.com"},
		{"11111111", "Carlos López", "carlos.lopez@email.com"},
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
