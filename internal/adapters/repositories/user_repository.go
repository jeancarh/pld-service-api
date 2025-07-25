package repositories

import (
	"crabi-test/internal/domain"
	"database/sql"
)

// UserRepository implementa el repositorio de usuarios con SQLite
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create crea un nuevo usuario en la base de datos
func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password, id_number, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.IDNumber, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	// Obtener el ID generado
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = uint(id)
	return nil
}

// GetByID obtiene un usuario por su ID
func (r *UserRepository) GetByID(id uint) (*domain.User, error) {
	query := `
		SELECT id, name, email, password, id_number, created_at, updated_at
		FROM users WHERE id = ?
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IDNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password, id_number, created_at, updated_at
		FROM users WHERE email = ?
	`

	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IDNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// Update actualiza un usuario existente
func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET name = ?, email = ?, password = ?, id_number = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.IDNumber, user.UpdatedAt, user.ID)
	return err
}

// Delete elimina un usuario por su ID
func (r *UserRepository) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
