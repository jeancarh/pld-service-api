package routes

import (
	"database/sql"

	"crabi-test/internal/adapters/repositories"
	"crabi-test/internal/application/services"
	"crabi-test/internal/infrastructure/external"
	"crabi-test/internal/infrastructure/http/handlers"
	"crabi-test/internal/infrastructure/http/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(r *gin.Engine, db *sql.DB) {
	// Crear instancias de repositorios
	userRepo := repositories.NewUserRepository(db)

	// Crear instancias de servicios externos
	pldClient := external.NewPLDClient()

	// Crear instancias de servicios de aplicación
	userService := services.NewUserService(userRepo, pldClient)
	authService := services.NewAuthService(userRepo)

	// Crear instancias de handlers
	userHandler := handlers.NewUserHandler(userService, authService)
	authHandler := handlers.NewAuthHandler(authService)

	// Crear middleware de autenticación
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Grupo de rutas de la API
	api := r.Group("/api/v1")

	// Rutas públicas
	{
		api.POST("/users", userHandler.CreateUser)
		api.POST("/auth/login", authHandler.Login)
	}

	// Rutas protegidas (requieren autenticación)
	protected := api.Group("")
	protected.Use(authMiddleware.Authenticate())
	{
		protected.GET("/users/me", userHandler.GetUser)
		protected.GET("/users/:id", userHandler.GetUserByID)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
	}
}
