package external

import (
	"bytes"
	"crabi-test/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// PLDClient implementa el cliente para el servicio externo de PLD
type PLDClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewPLDClient crea una nueva instancia del cliente PLD
func NewPLDClient() *PLDClient {
	baseURL := os.Getenv("PLD_SERVICE_URL")
	if baseURL == "" {
		panic("PLD_SERVICE_URL no está configurado en las variables de entorno")
	}

	return &PLDClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ValidateUser valida un usuario contra el servicio PLD
func (c *PLDClient) ValidateUser(idNumber, name, email string) (*domain.PLDResponse, error) {
	// Separar nombre en first_name y last_name
	names := strings.Fields(name)
	firstName := name
	lastName := ""

	if len(names) > 1 {
		firstName = names[0]
		lastName = strings.Join(names[1:], " ")
	}

	// Crear payload para la solicitud según la documentación real
	payload := map[string]string{
		"first_name": firstName,
		"last_name":  lastName,
		"email":      email,
	}

	// Serializar payload
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload: %w", err)
	}

	// Crear solicitud HTTP usando el endpoint real
	req, err := http.NewRequest("POST", c.baseURL+"/check-blacklist", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creando solicitud: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Realizar solicitud
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error realizando solicitud: %w", err)
	}
	defer resp.Body.Close()

	// Verificar código de respuesta (el servicio real retorna 201)
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("servicio PLD retornó código %d", resp.StatusCode)
	}

	// Decodificar respuesta según el formato real
	var pldResponseReal struct {
		IsInBlacklist bool `json:"is_in_blacklist"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&pldResponseReal); err != nil {
		return nil, fmt.Errorf("error decodificando respuesta: %w", err)
	}

	// Convertir a nuestro formato interno
	pldResponse := &domain.PLDResponse{
		IsBlacklisted: pldResponseReal.IsInBlacklist,
		Status:        "clean",
		Reason:        "",
	}

	if pldResponseReal.IsInBlacklist {
		pldResponse.Status = "blacklisted"
		pldResponse.Reason = "Usuario en lista negra"
	}

	return pldResponse, nil
}
