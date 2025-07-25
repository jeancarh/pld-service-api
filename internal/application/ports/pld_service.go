package ports

import "crabi-test/internal/domain"

// PLDService define las operaciones del servicio de PLD
type PLDService interface {
	ValidateUser(idNumber, name, email string) (*domain.PLDResponse, error)
}
