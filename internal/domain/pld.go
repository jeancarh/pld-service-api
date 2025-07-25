package domain

// PLDService define las operaciones del servicio de PLD
type PLDService interface {
	ValidateUser(idNumber, name, email string) (*PLDResponse, error)
}

// PLDResponse representa la respuesta del servicio PLD
type PLDResponse struct {
	IsBlacklisted bool   `json:"is_blacklisted"`
	Reason        string `json:"reason,omitempty"`
	Status        string `json:"status"`
}

// PLDRequest representa la solicitud al servicio PLD
type PLDRequest struct {
	IDNumber string `json:"id_number"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}
