package sqlite

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

// InitDB inicializa la conexión a la base de datos SQLite
func InitDB() (*sql.DB, error) {
	// Obtener ruta de la base de datos del environment
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// En Docker, usar base de datos en memoria
		if os.Getenv("DOCKER_ENV") == "true" {
			dbPath = ":memory:"
		} else {
			dbPath = "./crabi.db"
		}
	}

	// Abrir conexión a SQLite
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Verificar conexión
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Crear tablas si no existen
	if err := createTables(db); err != nil {
		return nil, err
	}

	log.Println("Base de datos SQLite inicializada correctamente")
	return db, nil
}

// createTables crea las tablas necesarias en la base de datos
func createTables(db *sql.DB) error {
	// Tabla de usuarios
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		id_number TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	`

	_, err := db.Exec(createUsersTable)
	if err != nil {
		return err
	}

	log.Println("Tablas creadas correctamente")
	return nil
}
